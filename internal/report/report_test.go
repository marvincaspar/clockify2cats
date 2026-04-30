package report

import (
	"errors"
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

	report, err := reporter.Generate(2022, 1, "Category", true, "")

	assert.Nil(t, err)
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

	report, err := reporter.Generate(2022, 1, "Category", true, "")
	assert.Nil(t, err)

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

	report, err := reporter.Generate(2022, 1, "CategoryChanged", false, "")
	assert.Nil(t, err)

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

	report, err := reporter.Generate(2022, 1, "Category", true, "")
	assert.Nil(t, err)

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

	report, err := reporter.Generate(2022, 1, "Category", true, "")
	assert.Nil(t, err)

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

	report, err := reporter.Generate(2022, 1, "Category", true, "")
	assert.Nil(t, err)

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

func TestReporter_ErrorOnShareBillableEntriesWhenNoBillableEntryIsGiven(t *testing.T) {
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
					Description: "NOT Billable Task 1",
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
						Name: "Project name (123)",
					},
					Billable: false,
				},
				{
					Description: "NOT Billable Task 2",
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

	report, err := reporter.Generate(2022, 1, "Category", true, "")

	assert.Equal(t, err.Error(), "No billable time entries found! Please distribute the shared time manually: https://app.clockify.me/timesheet")
	assert.Empty(t, report, "Report should be empty")
}

// Week 5 of 2022 runs Mon Jan 31 – Sun Feb 6, spanning two months.
func makeWeek5Entries() []ClockifyTimeEntry {
	makeEntry := func(start, duration, project string) ClockifyTimeEntry {
		return ClockifyTimeEntry{
			Description: "Task",
			TimeInterval: struct {
				Start    string `json:"start"`
				End      string `json:"end"`
				Duration string `json:"duration"`
			}{Start: start, Duration: duration},
			Project: struct {
				Name string `json:"name"`
			}{Name: project},
		}
	}
	return []ClockifyTimeEntry{
		makeEntry("2022-01-31T08:00:00.000Z", "PT1H", "Project (CATS1)"), // January (same as week start)
		makeEntry("2022-02-03T08:00:00.000Z", "PT2H", "Project (CATS2)"), // February (new month)
	}
}

func TestReporter_Generate_monthChangeEnd_keepsCurrentMonthEntries(t *testing.T) {
	reporter := Reporter{
		DescriptionDelimiter: "#",
		Repository:           repositoryMock{data: makeWeek5Entries()},
	}

	report, err := reporter.Generate(2022, 5, "ID", false, "end")
	assert.Nil(t, err)

	entities := strings.Split(strings.TrimRight(report, "\n"), "\n")
	assert.Equal(t, 1, len(entities), "only the January entry should be kept")
	assert.Contains(t, entities[0], "CATS1")
}

func TestReporter_Generate_monthChangeStart_keepsNewMonthEntries(t *testing.T) {
	reporter := Reporter{
		DescriptionDelimiter: "#",
		Repository:           repositoryMock{data: makeWeek5Entries()},
	}

	report, err := reporter.Generate(2022, 5, "ID", false, "start")
	assert.Nil(t, err)

	entities := strings.Split(strings.TrimRight(report, "\n"), "\n")
	assert.Equal(t, 1, len(entities), "only the February entry should be kept")
	assert.Contains(t, entities[0], "CATS2")
}

func TestReporter_Generate_projectWithoutParentheses_usesDashAsCatsID(t *testing.T) {
	reporter := Reporter{
		DescriptionDelimiter: "#",
		Repository: repositoryMock{data: []ClockifyTimeEntry{
			{
				Description: "Task",
				TimeInterval: struct {
					Start    string `json:"start"`
					End      string `json:"end"`
					Duration string `json:"duration"`
				}{Start: "2022-01-03T08:00:00.000Z", Duration: "PT1H"},
				Project: struct {
					Name string `json:"name"`
				}{Name: "Project without CATS ID"},
			},
		}},
	}

	report, err := reporter.Generate(2022, 1, "ID", false, "")
	assert.Nil(t, err)
	assert.NotEmpty(t, report)

	parts := strings.Split(strings.TrimRight(report, "\n"), "\t")
	assert.Equal(t, "-", parts[0])
}

func TestReporter_Generate_withOneDelimiterInDescription(t *testing.T) {
	reporter := Reporter{
		DescriptionDelimiter: "#",
		Repository: repositoryMock{data: []ClockifyTimeEntry{
			{
				Description: "Text field # Text2 field",
				TimeInterval: struct {
					Start    string `json:"start"`
					End      string `json:"end"`
					Duration string `json:"duration"`
				}{Start: "2022-01-03T08:00:00.000Z", Duration: "PT1H"},
				Project: struct {
					Name string `json:"name"`
				}{Name: "Project (123)"},
			},
		}},
	}

	report, err := reporter.Generate(2022, 1, "ID", true, "")
	assert.Nil(t, err)

	parts := strings.Split(strings.TrimRight(report, "\n"), "\t")
	assert.Equal(t, "Text field", parts[2])  // Text
	assert.Equal(t, "Text2 field", parts[3]) // Text2
	assert.Equal(t, "", parts[4])            // TextExternal empty
}

// TestReporter_ShareBillableEntries_MultipleSharedEntries is a regression test for a bug where
// multiple `*` entries caused over-distribution. Each subsequent shared entry used already-inflated
// durations as proportion base while billableSum remained the original value, producing a total
// larger than the actual tracked time.
//
// Example: 8h billable, shared1=2h, shared2=2h → expected 12h, buggy code produced 12.5h.
func TestReporter_ShareBillableEntries_MultipleSharedEntries(t *testing.T) {
	reporter := Reporter{
		DescriptionDelimiter: "#",
		Repository: repositoryMock{
			data: []ClockifyTimeEntry{
				// 8h billable entry
				{
					Description: "Billable Task",
					TimeInterval: struct {
						Start    string `json:"start"`
						End      string `json:"end"`
						Duration string `json:"duration"`
					}{
						Start:    "2022-01-03T07:00:00Z",
						Duration: "PT8H",
					},
					Project: struct {
						Name string `json:"name"`
					}{Name: "Project (123)"},
					Billable: true,
				},
				// first shared entry: 2h
				{
					Description: "Shared Task 1",
					TimeInterval: struct {
						Start    string `json:"start"`
						End      string `json:"end"`
						Duration string `json:"duration"`
					}{
						Start:    "2022-01-03T15:00:00Z",
						Duration: "PT2H",
					},
					Project: struct {
						Name string `json:"name"`
					}{Name: "Shared (*)"},
					Billable: false,
				},
				// second shared entry: 2h
				{
					Description: "Shared Task 2",
					TimeInterval: struct {
						Start    string `json:"start"`
						End      string `json:"end"`
						Duration string `json:"duration"`
					}{
						Start:    "2022-01-03T17:00:00Z",
						Duration: "PT2H",
					},
					Project: struct {
						Name string `json:"name"`
					}{Name: "Shared (*)"},
					Billable: false,
				},
			},
		},
	}

	report, err := reporter.Generate(2022, 1, "Category", false, "")
	assert.Nil(t, err)
	assert.NotEmpty(t, report)

	entities := strings.Split(strings.TrimRight(report, "\n"), "\n")
	assert.Equal(t, 1, len(entities), "should have 1 entry")

	parts := strings.Split(entities[0], "\t")
	assert.Equal(t, "123", parts[0])
	// 8h billable + 2h shared1 + 2h shared2 = 12h total
	// Bug produced 12,50 because second shared entry used inflated duration as proportion base
	assert.Equal(t, "12,00", parts[6], "total hours must equal billable + all shared (not over-distributed)")
}

type repositoryMock struct {
	data []ClockifyTimeEntry
	err  error
}

func (r repositoryMock) FetchClockifyData(start string) ([]ClockifyTimeEntry, error) {
	return r.data, r.err
}

func TestReporter_Generate_propagatesFetchError(t *testing.T) {
	reporter := Reporter{
		DescriptionDelimiter: "#",
		Repository:           repositoryMock{err: errors.New("network failure")},
	}

	report, err := reporter.Generate(2022, 1, "ID", false, "")
	assert.EqualError(t, err, "network failure")
	assert.Empty(t, report)
}
