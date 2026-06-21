package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	webhookURL string
	httpClient *http.Client
}

func NewClient(webhookURL string) *Client {
	return &Client{
		webhookURL: webhookURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) Send(message string) error {
	body, err := json.Marshal(struct {
		Content string `json:"content"`
	}{Content: message})
	if err != nil {
		return fmt.Errorf("encode Discord message: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt < 2; attempt++ {
		if err := c.send(body); err == nil {
			return nil
		} else {
			lastErr = err
		}
	}

	return fmt.Errorf("send Discord message after retry: %w", lastErr)
}

func (c *Client) send(body []byte) error {
	req, err := http.NewRequest(http.MethodPost, c.webhookURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create Discord request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("unexpected HTTP status %s", resp.Status)
	}
	return nil
}
