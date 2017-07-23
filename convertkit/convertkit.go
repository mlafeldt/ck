package convertkit

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
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
	ID           int               `json:"id"`
	FirstName    string            `json:"first_name"`
	EmailAddress string            `json:"email_address"`
	State        string            `json:"state"`
	CreatedAt    time.Time         `json:"created_at"`
	Fields       map[string]string `json:"fields"`
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
