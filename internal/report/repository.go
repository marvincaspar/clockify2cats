package report

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type RepositoryInterface interface {
	FetchClockifyData(start string) []ClockifyTimeEntry
}

type Repository struct {
	WorkspaceID string
	UserID      string
	ApiKey      string
}

func (r Repository) FetchClockifyData(start string) []ClockifyTimeEntry {
	startDate, _ := time.Parse(timeFormat, start)
	endDate := startDate.AddDate(0, 0, 7).Add(-time.Second)

	url := fmt.Sprintf("https://api.clockify.me/api/v1/workspaces/%s/user/%s/time-entries?hydrated=1&start=%s&end=%s&page-size=1000",
		r.WorkspaceID, r.UserID, startDate.Format(timeFormat), endDate.Format(timeFormat))

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("X-Api-Key", r.ApiKey)

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
