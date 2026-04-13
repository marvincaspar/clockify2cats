package report

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type RepositoryInterface interface {
	FetchClockifyData(start string) ([]ClockifyTimeEntry, error)
}

type Repository struct {
	WorkspaceID string
	UserID      string
	ApiKey      string
	BaseURL     string
	HTTPClient  *http.Client
}

func (r Repository) FetchClockifyData(start string) ([]ClockifyTimeEntry, error) {
	startDate, _ := time.Parse(timeFormat, start)
	endDate := startDate.AddDate(0, 0, 7).Add(-time.Second)

	baseURL := r.BaseURL
	if baseURL == "" {
		baseURL = "https://api.clockify.me"
	}

	url := fmt.Sprintf("%s/api/v1/workspaces/%s/user/%s/time-entries?hydrated=1&start=%s&end=%s&page-size=1000",
		baseURL, r.WorkspaceID, r.UserID, startDate.Format(timeFormat), endDate.Format(timeFormat))

	client := r.HTTPClient
	if client == nil {
		client = &http.Client{}
	}

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("X-Api-Key", r.ApiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Clockify API error: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var timeEntries []ClockifyTimeEntry
	if err := json.Unmarshal(body, &timeEntries); err != nil {
		return nil, fmt.Errorf("error parsing Clockify response: %w", err)
	}

	return timeEntries, nil
}
