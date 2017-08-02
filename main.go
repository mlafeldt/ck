package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"

	"github.com/mlafeldt/ck/convertkit"
)

func main() {
	config := convertkit.DefaultConfig()
	config.HTTPClient = &http.Client{Timeout: 10 * time.Second}

	rootCmd := &cobra.Command{
		Use:          "ck",
		Short:        "The ConvertKit Tool",
		SilenceUsage: true,
	}
	rootCmd.PersistentFlags().StringVar(&config.Key, "api-key", "", "Set API key for ConvertKit account")
	rootCmd.PersistentFlags().StringVar(&config.Secret, "api-secret", "", "Set API secret for ConvertKit account")
	rootCmd.PersistentFlags().StringVar(&config.Endpoint, "api-endpoint", "", "Set ConvertKit API endpoint")

	var (
		query convertkit.SubscriberQuery
		csv   bool
	)
	subscribersCmd := &cobra.Command{
		Use:   "subscribers",
		Short: "List subscribers",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, _ := convertkit.NewClient(config)
			subscribers, err := client.Subscribers(&query)
			if err != nil {
				return err
			}
			if csv {
				return outputCSV(os.Stdout, subscribers)
			}
			return outputTable(os.Stdout, subscribers)
		},
	}
	subscribersCmd.Flags().StringVar(&query.Since, "since", "", "Filter subscribers added on or after this date")
	subscribersCmd.Flags().StringVar(&query.Until, "until", "", "Filter subscribers added on or before this date")
	subscribersCmd.Flags().BoolVar(&query.Reverse, "reverse", false, "List subscribers in reverse order")
	subscribersCmd.Flags().BoolVar(&query.Cancelled, "cancelled", false, "List cancelled subscribers")
	subscribersCmd.Flags().StringVar(&query.EmailAddress, "email", "", "Filter subscribers by email address")
	subscribersCmd.Flags().BoolVar(&csv, "csv", false, "Output in CSV format")
	rootCmd.AddCommand(subscribersCmd)

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show program version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("ck %s %s/%s %s\n", Version,
				runtime.GOOS, runtime.GOARCH, runtime.Version())
		},
	}
	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func outputCSV(w io.Writer, subscribers []convertkit.Subscriber) error {
	records := [][]string{
		{"Email", "Signed up"},
	}
	for _, s := range subscribers {
		records = append(records, []string{
			s.EmailAddress,
			s.CreatedAt.Format(time.RFC3339),
		})
	}
	cw := csv.NewWriter(w)
	cw.WriteAll(records)
	return cw.Error()
}

func outputTable(w io.Writer, subscribers []convertkit.Subscriber) error {
	lines := []string{"#|Email|Signed up"}
	for i, s := range subscribers {
		lines = append(lines, fmt.Sprintf("%d|%s|%s",
			i+1,
			s.EmailAddress,
			s.CreatedAt.Format(time.RFC3339),
		))
	}
	_, err := fmt.Fprintln(w, columnize.SimpleFormat(lines))
	return err
}
