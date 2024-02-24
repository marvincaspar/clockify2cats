package report

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

type ClockifyTimeEntry struct {
	Description  string `json:"description"`
	TimeInterval struct {
		Start    string `json:"start"`
		End      string `json:"end"`
		Duration string `json:"duration"`
	} `json:"timeInterval"`
	Project struct {
		Name string `json:"name"`
	} `json:"project"`
}

type CatsEntity struct {
	CatsID    string
	Text      string
	Durations map[string]time.Duration
}

func Generate(workspaceID string, userID string, apiKey string, year int, week int, category string, withText bool) string {
	start := getFirstDayOfWeek(year, week).Format(timeFormat)

	timeEntries := fetchClockifyData(workspaceID, userID, apiKey, start)

	convertedTimeEntries := convertTimeEntries(start, timeEntries, withText)

	report := generateCatsReportData(convertedTimeEntries, category, withText)
	return report
}

func fetchClockifyData(workspaceID string, userID string, apiKey string, start string) []ClockifyTimeEntry {
	startDate, _ := time.Parse(timeFormat, start)
	endDate := startDate.AddDate(0, 0, 7).Add(-time.Second)

	url := fmt.Sprintf("https://api.clockify.me/api/v1/workspaces/%s/user/%s/time-entries?hydrated=1&start=%s&end=%s",
		workspaceID, userID, startDate.Format(timeFormat), endDate.Format(timeFormat))

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("X-Api-Key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var timeEntries []ClockifyTimeEntry
	_ = json.Unmarshal([]byte(body), &timeEntries)

	return timeEntries
}

func convertTimeEntries(start string, timeEntries []ClockifyTimeEntry, withText bool) []CatsEntity {
	startToDate, _ := time.Parse(timeFormat, start)
	catsEntries := []CatsEntity{}

	for _, timeEntry := range timeEntries {
		startDate, _ := time.Parse(timeFormat, timeEntry.TimeInterval.Start)
		cleanDuration := strings.ToLower(strings.Replace(timeEntry.TimeInterval.Duration, "PT", "", 1))
		duration, _ := time.ParseDuration(cleanDuration)

		catsID := getCatsID(timeEntry)
		text := ""
		if withText {
			text = timeEntry.Description
		}

		index := findCatsEntryID(catsEntries, catsID, text)

		if index == -1 {
			catsEntries = append(catsEntries, CatsEntity{
				CatsID: catsID,
				Text:   text,
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
		currentTime += duration

		catsEntries[index].Durations[startDate.Format("2006-01-02")] = currentTime
	}

	return catsEntries
}

func findCatsEntryID(catsEntries []CatsEntity, catsID string, text string) int {
	for i, catsEntry := range catsEntries {
		if catsEntry.CatsID == catsID && catsEntry.Text == text {
			return i
		}
	}

	return -1
}

func getCatsID(t ClockifyTimeEntry) string {
	r, _ := regexp.Compile("\\((.*)\\)")

	// Get CATS ID from project name
	if t.Project.Name != "" {
		match := r.FindAllStringSubmatch(t.Project.Name, -1)
		if len(match) > 0 {
			return match[0][1]
		}
	}

	return "-"
}

func generateCatsReportData(catsEntries []CatsEntity, category string, withText bool) string {
	var output string
	p := message.NewPrinter(language.Make("de-DE"))
	catsMeta := "%s\t\t%s\t%s\t"
	text := ""

	for _, catsEntry := range catsEntries {
		var values []interface{}

		// we have to sort the dates because looping over a map is not guaranteed to be in order
		dateKeys := make([]string, 0)
		for k, _ := range catsEntry.Durations {
			dateKeys = append(dateKeys, k)
		}
		sort.Strings(dateKeys)

		for _, date := range dateKeys {
			values = append(values, p.Sprintf("%.2f", catsEntry.Durations[date].Hours()))
		}
		if withText {
			text = catsEntry.Text
		}

		catsStart := fmt.Sprintf(catsMeta, catsEntry.CatsID, text, category)
		catsStart += strings.Repeat("%s\t\t", len(values))

		output += fmt.Sprintf(catsStart, values...) + "\n"
	}

	fmt.Printf(output)

	return output
}
