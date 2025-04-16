package report

import "time"

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
	Billable bool `json:"billable"`
}

type CatsEntity struct {
	CatsID       string
	Text         string
	Text2        string
	TextExternal string
	Durations    map[string]time.Duration
}
