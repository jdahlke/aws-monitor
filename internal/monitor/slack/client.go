package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/jdahlke/aws-monitor/internal/monitor"
)

type client struct {
	url        string
	channel    string
	httpClient monitor.HttpClient
}

type payload struct {
	Channel     string       `json:"channel,omitempty"`
	Username    string       `json:"username,omitempty"`
	Text        string       `json:"text,omitempty"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	Attachments []attachment `json:"attachments,omitempty"`
}

type attachment struct {
	Color   string  `json:"color,omitempty"`
	Pretext string  `json:"pretext,omitempty"`
	Fields  []field `json:"fields,omitempty"`
}

type field struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
	Short bool   `json:"short,omitempty"`
}

func NewClient(url, channel string, httpClient monitor.HttpClient) (*client, error) {
	if url == "" {
		return &client{}, fmt.Errorf("[Slack] url cannot be empty")
	}

	if channel == "" {
		return &client{}, fmt.Errorf("[Slack] channel cannot be empty")
	}

	return &client{
		url:        url,
		channel:    channel,
		httpClient: httpClient,
	}, nil
}

func (c *client) PostMessage(ctx context.Context, message *monitor.SlackMessage) error {
	// build payload
	fields := []field{}
	for key, value := range message.Details {
		fields = append(fields, field{Title: key, Value: value, Short: false})
	}
	pl := payload{
		Channel:   c.channel,
		Username:  "AWS Monitor Bot",
		Text:      message.Subject,
		IconEmoji: ":aws:",
		Attachments: []attachment{
			{
				Color:   message.Severity,
				Pretext: message.Message,
				Fields:  fields,
			},
		},
	}

	// serialisation
	body, err := json.Marshal(pl)
	if err != nil {
		return fmt.Errorf("[Slack] error marshal JSON: %v", err)
	}

	slog.Info(fmt.Sprintf("[Slack] payload: %s", body))

	// create a HTTP post request
	req, err := http.NewRequest("POST", c.url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("[Slack] error create HTTP request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("[Slack] error make HTTP request: %v", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("[Slack] error read response body: %v", err)
		}

		return fmt.Errorf("[Slack] error response: %s %v", resBody, res.StatusCode)
	}

	return nil
}
