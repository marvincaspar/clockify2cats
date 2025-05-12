package cmd

import (
	"testing"
	"time"

	"github.com/marvincaspar/clockify2cats/internal/report"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGenerateCmd(t *testing.T) {
	reporter := &report.Reporter{}
	cmd := newGenerateCmd(time.Now(), reporter)
	assert.NotNil(t, cmd)
}

func TestGenerateCmd_WithOptions(t *testing.T) {
	type args struct {
		flagCurrentWeek bool
		flagLastWeek    bool
		flagWeek        int
		wantYear        int
		wantWeek        int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "with current week",
			args: args{
				flagCurrentWeek: true,
				flagLastWeek:    false,
				flagWeek:        0,

				wantYear: 2024,
				wantWeek: 1,
			},
		},
		{
			name: "with last week",
			args: args{
				flagCurrentWeek: false,
				flagLastWeek:    true,
				flagWeek:        0,

				wantYear: 2023,
				wantWeek: 52,
			},
		},
		{
			name: "with week",
			args: args{
				flagCurrentWeek: false,
				flagLastWeek:    false,
				flagWeek:        10,

				wantYear: 2023,
				wantWeek: 10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(reporterMock)

			m.On("Generate", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("").Once()

			flagCurrentWeek = tt.args.flagCurrentWeek
			flagLastWeek = tt.args.flagLastWeek
			flagWeek = tt.args.flagWeek

			testTime := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
			cmd := newGenerateCmd(testTime, m)

			cmd.Run(cmd, []string{})

			m.AssertCalled(t, "Generate", tt.args.wantYear, tt.args.wantWeek, "ID", false)
		})
	}
}

// func TestGenerateCmd_WithDate1JanuarAndFlagCurrent(t *testing.T) {
// 	m := new(reporterMock)

// 	m.On("Generate", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("").Once()

// 	flagCurrentWeek = true
// 	testTime := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
// 	cmd := newGenerateCmd(testTime, m)

// 	cmd.Run(cmd, []string{})

// 	m.AssertCalled(t, "Generate", 2024, 1, "ID", false)
// }

// func TestGenerateCmd_WithDate1JanuarAndFlagLast(t *testing.T) {
// 	m := new(reporterMock)

// 	m.On("Generate", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("").Once()

// 	flagLastWeek = true
// 	testTime := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
// 	cmd := newGenerateCmd(testTime, m)

// 	cmd.Run(cmd, []string{})

// 	m.AssertCalled(t, "Generate", 2023, 52, "ID", false)
// }

// func TestGenerateCmd_WithDateDecemberAndFlagWeeknumber(t *testing.T) {
// 	m := new(reporterMock)

// 	m.On("Generate", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("").Once()

// 	flagWeek = 5
// 	testTime := time.Date(2024, time.December, 1, 0, 0, 0, 0, time.UTC)
// 	cmd := newGenerateCmd(testTime, m)

// 	cmd.Run(cmd, []string{})

// 	m.AssertCalled(t, "Generate", 2024, 5, "ID", false)
// }

type reporterMock struct{ mock.Mock }

func (m *reporterMock) Generate(year int, week int, category string, withText bool) (string, error) {
	args := m.Called(year, week, category, withText)
	return args.String(0), nil
}
