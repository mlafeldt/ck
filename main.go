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
	var (
		apiKey, apiSecret, apiEndpoint string
	)

	rootCmd := &cobra.Command{
		Use:          "ck",
		Short:        "The ConvertKit Tool",
		SilenceUsage: true,
	}
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "Set API key for ConvertKit account")
	rootCmd.PersistentFlags().StringVar(&apiSecret, "api-secret", "", "Set API secret for ConvertKit account")
	rootCmd.PersistentFlags().StringVar(&apiEndpoint, "api-endpoint", "", "Set ConvertKit API endpoint")

	subscribersCmd := &cobra.Command{
		Use:   "subscribers",
		Short: "List all confirmed subscribers",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := convertkit.DefaultConfig()
			if apiKey != "" {
				config.Key = apiKey
			}
			if apiSecret != "" {
				config.Secret = apiSecret
			}
			if apiEndpoint != "" {
				config.Endpoint = apiEndpoint
			}
			config.HTTPClient = &http.Client{Timeout: 10 * time.Second}
			client, _ := convertkit.NewClient(config)
			subscribers, err := client.Subscribers()
			if err != nil {
				return err
			}
			if csv, _ := cmd.Flags().GetBool("csv"); csv {
				return outputCSV(os.Stdout, subscribers)
			}
			return outputTable(os.Stdout, subscribers)
		},
	}
	subscribersCmd.Flags().Bool("csv", false, "Output in CSV format")
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
