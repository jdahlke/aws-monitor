package monitor

import (
	"context"
	"net/http"
)

type AwsEvent interface {
	// CreateSlackMessage creates new SlackMessage
	CreateSlackMessage(ctx context.Context) (*SlackMessage, error)

	// ReportEvent returns true or false
	ReportEvent(ctx context.Context) bool
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type SlackClient interface {
	// PostMessage posts SlackMessage to pre-defined Slack channel
	PostMessage(ctx context.Context, message *SlackMessage) error
}

type SlackMessage struct {
	Subject  string
	Message  string
	Details  map[string]string
	Severity string
}

type Severity string

const (
	SeverityDanger Severity = "danger"
	SeverityGood   Severity = "good"
)
