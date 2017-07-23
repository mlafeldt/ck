package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ryanuber/columnize"
)

type Config struct {
	Endpoint string
	Key      string
	Secret   string

	HTTPClient *http.Client
}

func DefaultConfig() *Config {
	c := Config{
		Endpoint:   "https://api.convertkit.com",
		HTTPClient: http.DefaultClient,
	}
	if v := os.Getenv("CONVERTKIT_ENDPOINT"); v != "" {
		c.Endpoint = v
	}
	if v := os.Getenv("CONVERTKIT_API_KEY"); v != "" {
		c.Key = v
	}
	if v := os.Getenv("CONVERTKIT_API_SECRET"); v != "" {
		c.Secret = v
	}
	return &c
}

type Client struct {
	config *Config
}

func NewClient(c *Config) (*Client, error) {
	defConfig := DefaultConfig()
	if c.Endpoint == "" {
		c.Endpoint = defConfig.Endpoint
	}
	if c.Key == "" {
		c.Key = defConfig.Key
	}
	if c.Secret == "" {
		c.Secret = defConfig.Secret
	}
	if c.HTTPClient == nil {
		c.HTTPClient = defConfig.HTTPClient
	}
	return &Client{config: c}, nil
}

type Subscriber struct {
	ID    int    `json:"id"`
	Email string `json:"email_address"`
}

type subscriberResponse struct {
	TotalSubscribers int          `json:"total_subscribers"`
	Page             int          `json:"page"`
	TotalPages       int          `json:"total_pages"`
	Subscribers      []Subscriber `json:"subscribers"`
}

func (c *Client) Subscribers() ([]Subscriber, error) {
	var subscribers []Subscriber
	page := 1

	for {
		url := fmt.Sprintf("%s/v3/subscribers?api_secret=%s&page=%d",
			c.config.Endpoint, c.config.Secret, page)

		var resp subscriberResponse
		if err := c.sendRequest("GET", url, nil, &resp); err != nil {
			return nil, err
		}
		subscribers = append(subscribers, resp.Subscribers...)

		if page >= resp.TotalPages {
			break
		}
		page += 1
	}

	return subscribers, nil
}

func (c *Client) sendRequest(method, url string, body io.Reader, out interface{}) error {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %s", resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(out)
}

func abort(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", a...)
	os.Exit(1)
}

func main() {
	config := DefaultConfig()
	client, _ := NewClient(config)
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
