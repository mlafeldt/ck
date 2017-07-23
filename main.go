package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/ryanuber/columnize"

	"github.com/mlafeldt/ck/convertkit"
)

func main() {
	var csvFormat = flag.Bool("csv", false, "Output in CSV format")
	flag.Parse()

	config := convertkit.DefaultConfig()
	config.HTTPClient = &http.Client{Timeout: 10 * time.Second}
	client, _ := convertkit.NewClient(config)
	subscribers, err := client.Subscribers()
	if err != nil {
		abort("%s", err)
	}

	if *csvFormat {
		err = outputCSV(os.Stdout, subscribers)
	} else {
		err = outputTable(os.Stdout, subscribers)
	}
	if err != nil {
		abort("%s", err)
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

func abort(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", a...)
	os.Exit(1)
}
