// Package convertkit provides a client to the ConvertKit API v3.
// See http://help.convertkit.com/article/33-api-documentation-v3
package convertkit

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"golang.org/x/sync/errgroup"
)

var (
	// ErrKeyMissing is returned when the API key is required, but not present.
	ErrKeyMissing = errors.New("ConvertKit API key missing")

	// ErrSecretMissing is returned when the API secret is required, but not present.
	ErrSecretMissing = errors.New("ConvertKit API secret missing")
)

// Config is used to configure the creation of the client.
type Config struct {
	Endpoint string
	Key      string
	Secret   string

	HTTPClient *http.Client
}

// DefaultConfig returns a default configuration for the client. It parses the
// environment variables CONVERTKIT_ENDPOINT, CONVERTKIT_API_KEY, and
// CONVERTKIT_API_SECRET.
func DefaultConfig() *Config {
	c := Config{
		Endpoint:   "https://api.convertkit.com",
		HTTPClient: http.DefaultClient,
	}
	if v := os.Getenv("CONVERTKIT_API_ENDPOINT"); v != "" {
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

// Client is the client to the ConvertKit API. Create a client with NewClient.
type Client struct {
	config *Config
}

// NewClient returns a new client for the given configuration.
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

// Subscriber describes a ConvertKit subscriber.
type Subscriber struct {
	ID           int               `json:"id"`
	FirstName    string            `json:"first_name"`
	EmailAddress string            `json:"email_address"`
	State        string            `json:"state"`
	CreatedAt    time.Time         `json:"created_at"`
	Fields       map[string]string `json:"fields"`
}

type subscriberPage struct {
	TotalSubscribers int          `json:"total_subscribers"`
	Page             int          `json:"page"`
	TotalPages       int          `json:"total_pages"`
	Subscribers      []Subscriber `json:"subscribers"`
}

// Subscribers returns a list of all confirmed subscribers.
func (c *Client) Subscribers() ([]Subscriber, error) {
	p, err := c.subscriberPage(1)
	if err != nil {
		return nil, err
	}

	total := p.TotalPages
	pages := make([]subscriberPage, total)
	pages[0] = *p

	// TODO: limit number of Go routines to be nicer to the API
	var g errgroup.Group
	for i := 2; i <= total; i++ {
		i := i // see https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {
			p, err := c.subscriberPage(i)
			if err == nil {
				pages[i-1] = *p
			}
			return err
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	var subscribers []Subscriber
	for i := 0; i < total; i++ {
		subscribers = append(subscribers, pages[i].Subscribers...)
	}

	return subscribers, nil
}

// TotalSubscribers returns the number of confirmed subscribers.
func (c *Client) TotalSubscribers() (int, error) {
	p, err := c.subscriberPage(1)
	if err != nil {
		return 0, err
	}
	return p.TotalSubscribers, nil
}

func (c *Client) subscriberPage(page int) (*subscriberPage, error) {
	if c.config.Secret == "" {
		return nil, ErrSecretMissing
	}

	url := fmt.Sprintf("%s/v3/subscribers?api_secret=%s&page=%d",
		c.config.Endpoint, c.config.Secret, page)

	var p subscriberPage
	if err := c.sendRequest("GET", url, nil, &p); err != nil {
		return nil, err
	}

	return &p, nil
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
