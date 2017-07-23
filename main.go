package main

import (
	"encoding/csv"
	"flag"
	"fmt"
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
		records := [][]string{
			{"Email", "Signed up"},
		}
		for _, s := range subscribers {
			records = append(records, []string{
				s.EmailAddress,
				s.CreatedAt.Format(time.RFC3339),
			})
		}
		w := csv.NewWriter(os.Stdout)
		w.WriteAll(records)
		if err := w.Error(); err != nil {
			abort("error writing CSV:", err)
		}
	} else {
		lines := []string{"#|Email|Signed Up"}
		for i, s := range subscribers {
			lines = append(lines, fmt.Sprintf("%d|%s|%s",
				i+1,
				s.EmailAddress,
				s.CreatedAt.Format(time.RFC3339),
			))
		}
		fmt.Println(columnize.SimpleFormat(lines))
	}
}

func abort(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", a...)
	os.Exit(1)
}
