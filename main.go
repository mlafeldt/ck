package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mlafeldt/ck/convertkit"

	"github.com/ryanuber/columnize"
)

func main() {
	client, _ := convertkit.NewClient(convertkit.DefaultConfig())
	subscribers, err := client.Subscribers()
	if err != nil {
		abort("%s", err)
	}

	lines := []string{"ID|Signed up|Email"}
	for _, s := range subscribers {
		lines = append(lines, fmt.Sprintf("%d|%s|%s",
			s.ID,
			s.CreatedAt.Format(time.RFC3339),
			s.EmailAddress,
		))
	}

	fmt.Println(columnize.SimpleFormat(lines))
}

func abort(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", a...)
	os.Exit(1)
}
