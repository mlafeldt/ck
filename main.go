package main

import (
	"fmt"
	"os"

	"github.com/mlafeldt/ck/convertkit"

	"github.com/ryanuber/columnize"
)

func main() {
	client, _ := convertkit.NewClient(convertkit.DefaultConfig())
	subscribers, err := client.Subscribers()
	if err != nil {
		abort("%s", err)
	}

	lines := []string{"ID|Email"}
	for _, s := range subscribers {
		lines = append(lines, fmt.Sprintf("%d|%s",
			s.ID,
			s.Email,
		))
	}
	fmt.Println(columnize.SimpleFormat(lines))
	fmt.Printf("%d subscribers\n", len(subscribers))
}

func abort(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", a...)
	os.Exit(1)
}
