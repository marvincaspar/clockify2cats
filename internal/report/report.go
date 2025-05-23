package report

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	timeFormat = "2006-01-02T15:04:05.999Z"
	category   string
)

type ReporterInterface interface {
	Generate(year int, week int, category string, withText bool) (string, error)
}

type Reporter struct {
	Repository           RepositoryInterface
	DescriptionDelimiter string
}

func (r Reporter) Generate(year int, week int, category string, withText bool) (string, error) {
	start := getFirstDayOfWeek(year, week).Format(timeFormat)

	timeEntries := r.Repository.FetchClockifyData(start)

	convertedTimeEntries, err := r.convertTimeEntries(start, timeEntries, withText)

	if err != nil {
		return "", err
	}

	report := r.generateCatsReportData(convertedTimeEntries, category, withText)
	return report, nil
}

func (r Reporter) convertTimeEntries(start string, timeEntries []ClockifyTimeEntry, withText bool) ([]CatsEntity, error) {
	startToDate, _ := time.Parse(timeFormat, start)
	catsEntries := []CatsEntity{}
	nonBillableCatsEntries := []CatsEntity{}
	sharedEntries := []ClockifyTimeEntry{}
	billableSum := 0.0

	for _, timeEntry := range timeEntries {
		startDate, _ := time.Parse(timeFormat, timeEntry.TimeInterval.Start)
		cleanDuration := strings.ToLower(strings.Replace(timeEntry.TimeInterval.Duration, "PT", "", 1))
		duration, _ := time.ParseDuration(cleanDuration)

		catsIDs := r.getCatsIDs(timeEntry)

		// Split entries into shared, billable and non-billable entries
		if catsIDs[0] == "*" {
			sharedEntries = append(sharedEntries, timeEntry)
		} else {
			if timeEntry.Billable {
				billableSum += duration.Hours()
				catsEntries = r.generateCATsEntriesFromTimeEntry(withText, timeEntry, catsIDs, duration, catsEntries, startToDate, startDate)
			} else {
				nonBillableCatsEntries = r.generateCATsEntriesFromTimeEntry(withText, timeEntry, catsIDs, duration, nonBillableCatsEntries, startToDate, startDate)
			}
		}
	}

	//Cancel when no billable entries are given to distribute existing shared ("*") entries
	if billableSum == 0 && len(sharedEntries) > 0 {
		return nil, errors.New("No billable time entries found! Please distribute the shared time manually: https://app.clockify.me/timesheet")
	}

	r.distributeSharedEntriesToBillableEntries(sharedEntries, catsEntries, billableSum)

	return append(catsEntries, nonBillableCatsEntries...), nil
}

func (r Reporter) distributeSharedEntriesToBillableEntries(sharedEntries []ClockifyTimeEntry, catsEntries []CatsEntity, billableSum float64) {
	if len(sharedEntries) > 0 {
		for _, sharedTimeEntry := range sharedEntries {
			cleanDuration := strings.ToLower(strings.Replace(sharedTimeEntry.TimeInterval.Duration, "PT", "", 1))
			sharedDuration, _ := time.ParseDuration(cleanDuration)
			for _, catsEntry := range catsEntries {
				for key, duration := range catsEntry.Durations {
					sharedDurationForEntry := (duration.Hours() / billableSum) * sharedDuration.Hours()
					catsEntry.Durations[key] += time.Duration(sharedDurationForEntry * float64(time.Hour))
				}
			}
		}
	}
}

func (r Reporter) generateCATsEntriesFromTimeEntry(withText bool, timeEntry ClockifyTimeEntry, catsIDs []string, duration time.Duration, catsEntries []CatsEntity, startToDate time.Time, startDate time.Time) []CatsEntity {
	text := []string{"", "", ""}
	if withText {
		text = r.splitDescription(timeEntry.Description)
	}

	for _, catsID := range catsIDs {
		durationShared := duration / time.Duration(len(catsIDs))
		trimmedCatsID := strings.TrimSpace(catsID)
		index := r.findCatsEntryID(catsEntries, trimmedCatsID, text[0], text[1], text[2])

		if index == -1 {
			catsEntries = append(catsEntries, CatsEntity{
				CatsID:       trimmedCatsID,
				Text:         text[0],
				Text2:        text[1],
				TextExternal: text[2],
				Durations: map[string]time.Duration{
					startToDate.AddDate(0, 0, 0).Format("2006-01-02"): time.Duration(0), // Monday
					startToDate.AddDate(0, 0, 1).Format("2006-01-02"): time.Duration(0), // Tuesday
					startToDate.AddDate(0, 0, 2).Format("2006-01-02"): time.Duration(0), // Wednesday
					startToDate.AddDate(0, 0, 3).Format("2006-01-02"): time.Duration(0), // Thursday
					startToDate.AddDate(0, 0, 4).Format("2006-01-02"): time.Duration(0), // Friday
					startToDate.AddDate(0, 0, 5).Format("2006-01-02"): time.Duration(0), // Saturday
					startToDate.AddDate(0, 0, 6).Format("2006-01-02"): time.Duration(0), // Sunday
				},
			},
			)
			index = len(catsEntries) - 1
		}

		currentTime := catsEntries[index].Durations[startDate.Format("2006-01-02")]
		currentTime += durationShared

		catsEntries[index].Durations[startDate.Format("2006-01-02")] = currentTime
	}
	return catsEntries
}

func (r Reporter) findCatsEntryID(catsEntries []CatsEntity, catsID string, text string, text2 string, textExternal string) int {
	for i, catsEntry := range catsEntries {
		if catsEntry.CatsID == catsID && catsEntry.Text == text && catsEntry.Text2 == text2 && catsEntry.TextExternal == textExternal {
			return i
		}
	}

	return -1
}

func (r Reporter) getCatsIDs(t ClockifyTimeEntry) []string {
	rg, _ := regexp.Compile("\\((.*)\\)")

	// Get CATS IDs from project name
	if t.Project.Name != "" {
		match := rg.FindAllStringSubmatch(t.Project.Name, -1)
		if len(match) > 0 {
			return strings.Split(match[0][1], ",")
		}
	}

	return []string{"-"}
}

func (r Reporter) splitDescription(description string) []string {
	parts := strings.Split(description, r.DescriptionDelimiter)
	if len(parts) >= 3 {
		return []string{strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), strings.TrimSpace(parts[2])}
	}

	if len(parts) == 2 {
		return []string{strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), ""}
	}

	if len(parts) == 1 {
		return []string{"", strings.TrimSpace(parts[0]), ""}
	}

	return []string{"", description, ""}
}

func (r Reporter) generateCatsReportData(catsEntries []CatsEntity, category string, withText bool) string {
	var output string
	p := message.NewPrinter(language.Make("de-DE"))
	catsMeta := "%s\t\t%s\t%s\t%s\t%s\t" // Rec. order, Description (empty), Text, Text 2, Text External, Category
	text := ""
	text2 := ""
	textExternal := ""

	for _, catsEntry := range catsEntries {
		var values []interface{}

		// we have to sort the dates because looping over a map is not guaranteed to be in order
		dateKeys := make([]string, 0)
		for k := range catsEntry.Durations {
			dateKeys = append(dateKeys, k)
		}
		sort.Strings(dateKeys)

		for _, date := range dateKeys {
			values = append(values, p.Sprintf("%.2f", catsEntry.Durations[date].Hours()))
		}
		if withText {
			text = catsEntry.Text
			text2 = catsEntry.Text2
			textExternal = catsEntry.TextExternal
		}

		catsStart := fmt.Sprintf(catsMeta, catsEntry.CatsID, text, text2, textExternal, category)
		catsStart += strings.Repeat("%s\t\t", len(values))

		output += fmt.Sprintf(catsStart, values...) + "\n"
	}

	return output
}
