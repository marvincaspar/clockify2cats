package report

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetFirstDayOfWeek(t *testing.T) {
	tests := []struct {
		name        string
		year        int
		week        int
		wantYear    int
		wantMonth   time.Month
		wantDay     int
		wantWeekday time.Weekday
	}{
		{
			name:        "week 1 of 2022 starts on Monday Jan 3",
			year:        2022,
			week:        1,
			wantYear:    2022,
			wantMonth:   time.January,
			wantDay:     3,
			wantWeekday: time.Monday,
		},
		{
			name:        "week 1 of 2024 starts on Monday Jan 1",
			year:        2024,
			week:        1,
			wantYear:    2024,
			wantMonth:   time.January,
			wantDay:     1,
			wantWeekday: time.Monday,
		},
		{
			name:        "week 5 of 2022 starts on Monday Jan 31",
			year:        2022,
			week:        5,
			wantYear:    2022,
			wantMonth:   time.January,
			wantDay:     31,
			wantWeekday: time.Monday,
		},
		{
			name:        "week 52 of 2022 starts on Monday Dec 26",
			year:        2022,
			week:        52,
			wantYear:    2022,
			wantMonth:   time.December,
			wantDay:     26,
			wantWeekday: time.Monday,
		},
		{
			name:        "week 10 of 2023 starts on Monday Mar 6",
			year:        2023,
			week:        10,
			wantYear:    2023,
			wantMonth:   time.March,
			wantDay:     6,
			wantWeekday: time.Monday,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getFirstDayOfWeek(tt.year, tt.week)
			assert.Equal(t, tt.wantWeekday, got.Weekday(), "should be a Monday")
			assert.Equal(t, tt.wantYear, got.Year())
			assert.Equal(t, tt.wantMonth, got.Month())
			assert.Equal(t, tt.wantDay, got.Day())
		})
	}
}
