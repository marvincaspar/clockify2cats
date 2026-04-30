package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/atotto/clipboard"
	"github.com/marvincaspar/clockify2cats/internal/report"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	flagWeek            int
	flagLastWeek        bool
	flagCurrentWeek     bool
	flagMonthChange     string
	flagCopyToClipboard bool
	flagCategory        string
	flagWithText        bool
)

func newGenerateCmd(t time.Time, reporter report.ReporterInterface) *cobra.Command {
	return &cobra.Command{
		Use:   "generate",
		Short: "Generate report for a specific week",
		Long:  `Generate a report from your clockify data for a specific week and print it to stdout. You can also copy it to the clipboard.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if flagWeek > 53 {
				return fmt.Errorf("invalid value %d for --week: must be between 1 and 53", flagWeek)
			}
			if flagMonthChange != "" && flagMonthChange != "start" && flagMonthChange != "end" {
				return fmt.Errorf("invalid value %q for --month-boundary: must be \"start\" or \"end\"", flagMonthChange)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			var week int
			year, currentWeek := t.ISOWeek()

			if flagCurrentWeek {
				week = currentWeek
			} else if flagLastWeek {
				week = currentWeek - 1

				// if the current week is the first week of the year we need to go back to the previous year
				if currentWeek == 1 {
					year, currentWeek = t.Add(-time.Hour * 24 * 7).ISOWeek()
					week = currentWeek
				}
			} else if flagWeek > 0 {
				weekInput := flagWeek

				if weekInput > currentWeek {
					year = year - 1
				}

				week = weekInput

			} else {
				fmt.Fprintf(os.Stderr, "Error: no week specified")
				os.Exit(1)
			}

			report, totalHours, err := reporter.Generate(
				year,
				week,
				flagCategory,
				flagWithText,
				flagMonthChange,
			)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(1)
			}

			fmt.Println(report)

			if flagCopyToClipboard {
				if err := clipboard.WriteAll(report); err != nil {
					fmt.Fprintf(os.Stderr, "Error: could not copy to clipboard: %s\nMake sure xclip, xsel (X11) or wl-clipboard (Wayland) is installed.\n", err)
				}
			}

			fmt.Printf("Total: %.2fh\n", totalHours)
		},
	}
}

func init() {
	initConfig()
	workspaceID := viper.GetString("workspace-id")
	userID := viper.GetString("user-id")
	apiKey := viper.GetString("api-key")
	descriptionDelimiter := viper.GetString("description-delimiter")

	clockifyRepository := report.Repository{
		WorkspaceID: workspaceID,
		UserID:      userID,
		ApiKey:      apiKey,
	}

	reporter := report.Reporter{
		Repository:           clockifyRepository,
		DescriptionDelimiter: descriptionDelimiter,
	}

	t := time.Now()
	generateCmd := newGenerateCmd(t, &reporter)

	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")
	generateCmd.Flags().IntVarP(&flagWeek, "week", "w", 0, "Week number")
	generateCmd.Flags().BoolVarP(&flagLastWeek, "last", "l", false, "Last week")
	generateCmd.Flags().BoolVarP(&flagCurrentWeek, "current", "c", false, "Current week")
	generateCmd.MarkFlagsOneRequired("week", "last", "current")
	generateCmd.MarkFlagsMutuallyExclusive("week", "last", "current")
	generateCmd.Flags().StringVarP(&flagMonthChange, "month-boundary", "m", "", `Filter entries for weeks spanning a month boundary: "start" keeps the new month, "end" keeps the current month`)

	generateCmd.Flags().BoolVarP(&flagCopyToClipboard, "copy", "C", false, "Copy report to clipboard")

	generateCmd.Flags().StringVar(&flagCategory, "category", "ID", "Category identifier")
	generateCmd.Flags().BoolVarP(&flagWithText, "text", "t", false, "Print with text")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
