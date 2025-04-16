package report

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReporter_Generate_withTwoTimeEntriesForTheSameProjectAndDescription(t *testing.T) {
	reporter := Reporter{
		DescriptionDelimiter: "#",
		Repository: repositoryMock{
			data: []ClockifyTimeEntry{
				{
					Description: "Task description",
					TimeInterval: struct {
						Start    string `json:"start"`
						End      string `json:"end"`
						Duration string `json:"duration"`
					}{
						Start:    "2022-01-01T07:00:00Z",
						End:      "2022-01-01T08:00:00Z",
						Duration: "PT1H",
					},
					Project: struct {
						Name string `json:"name"`
					}{
						Name: "Project name (123)",
					},
				},
				{
					Description: "Task description",
					TimeInterval: struct {
						Start    string `json:"start"`
						End      string `json:"end"`
						Duration string `json:"duration"`
					}{
						Start:    "2022-01-01T09:00:00Z",
						End:      "2022-01-01T10:00:00Z",
						Duration: "PT1H",
					},
					Project: struct {
						Name string `json:"name"`
					}{
						Name: "Project name (123)",
					},
				},
			},
		},
	}

	report := reporter.Generate(2022, 1, "Category", true)

	assert.NotEmpty(t, report, "Report should not be empty")

	entities := strings.Split(report, "\n")
	entities = entities[:len(entities)-1]
	assert.Equal(t, 1, len(entities), "Time entries should contain 1 lines")

	parts := strings.Split(entities[0], "\t")
	assert.Equal(t, 23, len(parts), "Time entries should contain 23 columns")
	assert.Equal(t, "123", parts[0])
	assert.Equal(t, "", parts[2])
	assert.Equal(t, "Task description", parts[3])
	assert.Equal(t, "", parts[4])
	assert.Equal(t, "Category", parts[5])
	assert.Equal(t, "2,00", parts[6])
}

func TestReporter_Generate_withTwoTimeEntriesForTheSameProjectAndDifferentDescription(t *testing.T) {
	reporter := Reporter{
		DescriptionDelimiter: "#",
		Repository: repositoryMock{
			data: []ClockifyTimeEntry{
				{
					Description: "Task description",
					TimeInterval: struct {
						Start    string `json:"start"`
						End      string `json:"end"`
						Duration string `json:"duration"`
					}{
						Start:    "2022-01-01T07:00:00Z",
						End:      "2022-01-01T08:00:00Z",
						Duration: "PT1H",
					},
					Project: struct {
						Name string `json:"name"`
					}{
						Name: "Project name (123)",
					},
				},
				{
					Description: "Task description 2",
					TimeInterval: struct {
						Start    string `json:"start"`
						End      string `json:"end"`
						Duration string `json:"duration"`
					}{
						Start:    "2022-01-01T09:00:00Z",
						End:      "2022-01-01T10:00:00Z",
						Duration: "PT1H",
					},
					Project: struct {
						Name string `json:"name"`
					}{
						Name: "Project name (123)",
					},
				},
			},
		},
	}

	report := reporter.Generate(2022, 1, "Category", true)

	assert.NotEmpty(t, report, "Report should not be empty")

	entities := strings.Split(report, "\n")
	entities = entities[:len(entities)-1]
	assert.Equal(t, 2, len(entities), "Time entries should contain 2 lines")

	parts := strings.Split(entities[0], "\t")
	assert.Equal(t, 23, len(parts), "Time entries should contain 23 columns")
	assert.Equal(t, "123", parts[0])
	assert.Equal(t, "", parts[2])
	assert.Equal(t, "Task description", parts[3])
	assert.Equal(t, "", parts[4])
	assert.Equal(t, "Category", parts[5])
	assert.Equal(t, "1,00", parts[6])

	parts = strings.Split(entities[1], "\t")
	assert.Equal(t, 23, len(parts), "Time entries should contain 23 columns")
	assert.Equal(t, "123", parts[0])
	assert.Equal(t, "", parts[2])
	assert.Equal(t, "Task description 2", parts[3])
	assert.Equal(t, "", parts[4])
	assert.Equal(t, "Category", parts[5])
	assert.Equal(t, "1,00", parts[6])
}

func TestReporter_Generate_withTwoTimeEntriesForTheDifferentProjectAndDifferentDescriptionButWithoutTextOutput(t *testing.T) {
	reporter := Reporter{
		DescriptionDelimiter: "#",
		Repository: repositoryMock{
			data: []ClockifyTimeEntry{
				{
					Description: "Task description",
					TimeInterval: struct {
						Start    string `json:"start"`
						End      string `json:"end"`
						Duration string `json:"duration"`
					}{
						Start:    "2022-01-01T07:00:00Z",
						End:      "2022-01-01T08:00:00Z",
						Duration: "PT1H",
					},
					Project: struct {
						Name string `json:"name"`
					}{
						Name: "Project name (123)",
					},
				},
				{
					Description: "Task description",
					TimeInterval: struct {
						Start    string `json:"start"`
						End      string `json:"end"`
						Duration string `json:"duration"`
					}{
						Start:    "2022-01-01T09:00:00Z",
						End:      "2022-01-01T10:00:00Z",
						Duration: "PT1H",
					},
					Project: struct {
						Name string `json:"name"`
					}{
						Name: "Project name 2 (456)",
					},
				},
			},
		},
	}

	report := reporter.Generate(2022, 1, "CategoryChanged", false)

	assert.NotEmpty(t, report, "Report should not be empty")

	entities := strings.Split(report, "\n")
	entities = entities[:len(entities)-1]
	assert.Equal(t, 2, len(entities), "Time entries should contain 2 lines")

	parts := strings.Split(entities[0], "\t")
	assert.Equal(t, 23, len(parts), "Time entries should contain 23 columns")
	assert.Equal(t, "123", parts[0])
	assert.Equal(t, "", parts[2])
	assert.Equal(t, "CategoryChanged", parts[5])
	assert.Equal(t, "1,00", parts[6])

	parts = strings.Split(entities[1], "\t")
	assert.Equal(t, 23, len(parts), "Time entries should contain 23 columns")
	assert.Equal(t, "456", parts[0])
	assert.Equal(t, "", parts[2])
	assert.Equal(t, "CategoryChanged", parts[5])
	assert.Equal(t, "1,00", parts[6])
}

func TestReporter_Generate_withSharedTimeEntriesForTwoProjects(t *testing.T) {
	reporter := Reporter{
		DescriptionDelimiter: "#",
		Repository: repositoryMock{
			data: []ClockifyTimeEntry{
				{
					Description: "Task description",
					TimeInterval: struct {
						Start    string `json:"start"`
						End      string `json:"end"`
						Duration string `json:"duration"`
					}{
						Start:    "2022-01-01T07:00:00Z",
						End:      "2022-01-01T08:00:00Z",
						Duration: "PT1H",
					},
					Project: struct {
						Name string `json:"name"`
					}{
						Name: "Project name (123, 456)",
					},
				},
			},
		},
	}

	report := reporter.Generate(2022, 1, "Category", true)

	assert.NotEmpty(t, report, "Report should not be empty")

	entities := strings.Split(report, "\n")
	entities = entities[:len(entities)-1]
	assert.Equal(t, 2, len(entities), "Time entries should contain 2 lines")

	parts := strings.Split(entities[0], "\t")
	assert.Equal(t, 23, len(parts), "Time entries should contain 23 columns")
	assert.Equal(t, "123", parts[0])
	assert.Equal(t, "", parts[2])
	assert.Equal(t, "Task description", parts[3])
	assert.Equal(t, "", parts[4])
	assert.Equal(t, "Category", parts[5])
	assert.Equal(t, "0,50", parts[6])

	parts = strings.Split(entities[1], "\t")
	assert.Equal(t, 23, len(parts), "Time entries should contain 23 columns")
	assert.Equal(t, "456", parts[0])
	assert.Equal(t, "", parts[2])
	assert.Equal(t, "Task description", parts[3])
	assert.Equal(t, "", parts[4])
	assert.Equal(t, "Category", parts[5])
	assert.Equal(t, "0,50", parts[6])
}

func TestReporter_Generate_withDescriptionBlockTimeEntry(t *testing.T) {
	reporter := Reporter{
		DescriptionDelimiter: "#",
		Repository: repositoryMock{
			data: []ClockifyTimeEntry{
				{
					Description: "Part 1 # Part 2 # Part 3",
					TimeInterval: struct {
						Start    string `json:"start"`
						End      string `json:"end"`
						Duration string `json:"duration"`
					}{
						Start:    "2022-01-01T07:00:00Z",
						End:      "2022-01-01T08:00:00Z",
						Duration: "PT1H",
					},
					Project: struct {
						Name string `json:"name"`
					}{
						Name: "Project name (123)",
					},
				},
			},
		},
	}

	report := reporter.Generate(2022, 1, "Category", true)

	assert.NotEmpty(t, report, "Report should not be empty")

	entities := strings.Split(report, "\n")
	entities = entities[:len(entities)-1]
	assert.Equal(t, 1, len(entities), "Time entries should contain 1 line")

	parts := strings.Split(entities[0], "\t")
	assert.Equal(t, 23, len(parts), "Time entries should contain 23 columns")
	assert.Equal(t, "123", parts[0])
	assert.Equal(t, "Part 1", parts[2])
	assert.Equal(t, "Part 2", parts[3])
	assert.Equal(t, "Part 3", parts[4])
	assert.Equal(t, "Category", parts[5])
	assert.Equal(t, "1,00", parts[6])
}

func TestReporter_ShareBillableEntries(t *testing.T) {
	reporter := Reporter{
		DescriptionDelimiter: "#",
		Repository: repositoryMock{
			data: []ClockifyTimeEntry{
				{
					Description: "Shared Task",
					TimeInterval: struct {
						Start    string `json:"start"`
						End      string `json:"end"`
						Duration string `json:"duration"`
					}{
						Start:    "2022-01-01T07:00:00Z",
						End:      "2022-01-01T08:00:00Z",
						Duration: "PT1H",
					},
					Project: struct {
						Name string `json:"name"`
					}{
						Name: "Project name (*)",
					},
					Billable: false,
				},
				{
					Description: "Billable Task",
					TimeInterval: struct {
						Start    string `json:"start"`
						End      string `json:"end"`
						Duration string `json:"duration"`
					}{
						Start:    "2022-01-01T09:00:00Z",
						End:      "2022-01-01T10:00:00Z",
						Duration: "PT1H",
					},
					Project: struct {
						Name string `json:"name"`
					}{
						Name: "Project name (123)",
					},
					Billable: true,
				},
				{
					Description: "Third Billable Task",
					TimeInterval: struct {
						Start    string `json:"start"`
						End      string `json:"end"`
						Duration string `json:"duration"`
					}{
						Start:    "2022-01-01T11:00:00Z",
						End:      "2022-01-01T12:00:00Z",
						Duration: "PT1H",
					},
					Project: struct {
						Name string `json:"name"`
					}{
						Name: "Project name (789)",
					},
					Billable: true,
				},
				{
					Description: "NOT Billable Task",
					TimeInterval: struct {
						Start    string `json:"start"`
						End      string `json:"end"`
						Duration string `json:"duration"`
					}{
						Start:    "2022-01-01T13:00:00Z",
						End:      "2022-01-01T14:00:00Z",
						Duration: "PT1H",
					},
					Project: struct {
						Name string `json:"name"`
					}{
						Name: "Project name (000)",
					},
					Billable: false,
				},
			},
		},
	}

	report := reporter.Generate(2022, 1, "Category", true)

	assert.NotEmpty(t, report, "Report should not be empty")

	entities := strings.Split(report, "\n")
	entities = entities[:len(entities)-1]
	assert.Equal(t, 3, len(entities), "Time entries should contain 3 lines")

	// Check shared entry distribution
	entry123 := strings.Split(entities[0], "\t")
	assert.Equal(t, "123", entry123[0])
	assert.Equal(t, "1,50", entry123[6]) // 1 hour + 0.5 hours shared

	entry789 := strings.Split(entities[1], "\t")
	assert.Equal(t, "789", entry789[0])
	assert.Equal(t, "1,50", entry789[6]) // 1 hour + 0.5 hours shared

	entry000 := strings.Split(entities[2], "\t")
	assert.Equal(t, "000", entry000[0])
	assert.Equal(t, "1,00", entry000[6]) // 1 hour without shared hours because not billable
}

type repositoryMock struct {
	data []ClockifyTimeEntry
}

func (r repositoryMock) FetchClockifyData(start string) []ClockifyTimeEntry {
	return r.data
}
