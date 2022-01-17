package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	timeFormat = "2006-01-02T15:04:05.999Z"
	start      string
	week       int
	category   string
	withText   bool
	fileOutput bool
)

const usage = `Usage of clockify2cats:
  -w, --week week number for report (don't use in combination with start)
  -s, --start Startdate for report YYYY-MM-DD (don't use in combination with week)
  -c, --category Category identifyer (default: ID)
  -t, --text Add Clockify description as text to report
  -f, --file Write report to file
`

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

func main() {
	godotenv.Load()
	validateEnvironmentVariables()
	parseFlags()

	if start == "" && week == 0 {
		log.Fatal("No start date or calendar week specified")
	}

	start = start + "T00:00:00Z" //"2022-01-10T00:00:00Z"
	if week > 0 {
		year := time.Now().Year()
		timezone := time.FixedZone("Europe/Berlin", 0)
		start = firstDayOfISOWeek(year, week, timezone).Format(timeFormat)
	}

	timeEntries := fetchClockifyData(start)

	convertedTimeEntries := convertTimeEntries(start, timeEntries, withText)

	data := generateCatsReportData(convertedTimeEntries, withText)

	if fileOutput {
		writeToFile(start, data)
	}
}

func validateEnvironmentVariables() {
	workspaceId := os.Getenv("CLOCKIFY_WORKSPACE_ID")
	if workspaceId == "" {
		log.Fatal("CLOCKIFY_WORKSPACE_ID is not set")
	}
	userId := os.Getenv("CLOCKIFY_USER_ID")
	if userId == "" {
		log.Fatal("CLOCKIFY_USER_ID is not set")
	}
	apiKey := os.Getenv("CLOCKIFY_API_KEY")
	if apiKey == "" {
		log.Fatal("CLOCKIFY_API_KEY is not set")
	}
}

func parseFlags() {
	flag.IntVar(&week, "week", 0, "Calendar week")
	flag.IntVar(&week, "w", 0, "Calendar week")
	flag.StringVar(&start, "start", "", "Start date")
	flag.StringVar(&start, "s", "", "Start date")
	flag.StringVar(&category, "category", "ID", "Category identifyer")
	flag.StringVar(&category, "c", "ID", "Category identifyer")
	flag.BoolVar(&withText, "text", false, "Print with text")
	flag.BoolVar(&withText, "t", false, "Print with text")
	flag.BoolVar(&fileOutput, "file", false, "Write report to file")
	flag.BoolVar(&fileOutput, "f", false, "Write report to file")
	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()
}

func fetchClockifyData(start string) []ClockifyTimeEntry {
	startDate, _ := time.Parse(timeFormat, start)
	endDate := startDate.AddDate(0, 0, 7)

	workspaceId := os.Getenv("CLOCKIFY_WORKSPACE_ID")
	userId := os.Getenv("CLOCKIFY_USER_ID")
	apiKey := os.Getenv("CLOCKIFY_API_KEY")

	url := fmt.Sprintf("https://api.clockify.me/api/v1/workspaces/%s/user/%s/time-entries?hydrated=1&start=%s&end=%s",
		workspaceId, userId, startDate.Format(timeFormat), endDate.Format(timeFormat))

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("X-Api-Key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
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
				Text:   timeEntry.Description,
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

func generateCatsReportData(catsEntries []CatsEntity, withText bool) string {
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

func writeToFile(start string, output string) {
	startDate, _ := time.Parse(timeFormat, start)
	filename := fmt.Sprintf("cats_report_%s.csv", startDate.Format("2006-01-02"))
	ioutil.WriteFile(filename, []byte(output), 0644)
}

// https://xferion.com/golang-reverse-isoweek-get-the-date-of-the-first-day-of-iso-week/
func firstDayOfISOWeek(year int, week int, timezone *time.Location) time.Time {
	date := time.Date(year, 0, 0, 0, 0, 0, 0, timezone)
	isoYear, isoWeek := date.ISOWeek()

	// iterate back to Monday
	for date.Weekday() != time.Monday {
		date = date.AddDate(0, 0, -1)
		isoYear, isoWeek = date.ISOWeek()
	}

	// iterate forward to the first day of the first week
	for isoYear < year {
		date = date.AddDate(0, 0, 7)
		isoYear, isoWeek = date.ISOWeek()
	}

	// iterate forward to the first day of the given week
	for isoWeek < week {
		date = date.AddDate(0, 0, 7)
		isoYear, isoWeek = date.ISOWeek()
	}

	return date
}
