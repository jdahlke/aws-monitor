package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type codePipelineEvent struct {
	event  events.EventBridgeEvent
	detail codePipelineDetail
}

type codePipelineDetail struct {
	Pipeline    string `json:"pipeline,omitempty"`
	ExecutionID string `json:"execution-id,omitempty"`
	Stage       string `json:"stage,omitempty"`
	Action      string `json:"action,omitempty"`
	State       string `json:"state,omitempty"`
}

// NewCodePipelineActionEvent returns new CodePipeline
func NewCodePipelineEvent(evt events.EventBridgeEvent) (*codePipelineEvent, error) {
	var detail codePipelineDetail
	if err := json.Unmarshal(evt.Detail, &detail); err != nil {
		return nil, err
	}

	return &codePipelineEvent{
		event:  evt,
		detail: detail,
	}, nil
}

// CreateSlackMessage creates new SlackMessage for CodePipelineAction state change event
func (e codePipelineEvent) CreateSlackMessage(ctx context.Context) (*SlackMessage, error) {
	var title, link, severity string

	lowerState := strings.ToLower(string(e.detail.State))

	title = fmt.Sprintf("*Code Pipeline %v: %v %v*", e.detail.Pipeline, e.detail.Action, lowerState)
	link = fmt.Sprintf("<https://%v.console.aws.amazon.com/codesuite/codepipeline/pipelines/%v/executions/%v/timeline|view details>", e.event.Region, e.detail.Pipeline, e.detail.ExecutionID)

	switch e.detail.State {
	case "FAILED":
		severity = "danger"
	default:
		severity = "good"
	}

	return &SlackMessage{
		Details: map[string]string{
			title: link,
		},
		Severity: severity,
	}, nil
}

func (e codePipelineEvent) ReportEvent(ctx context.Context) bool {
	switch e.detail.State {
	case "FAILED":
		return true
	case "SUCCEEDED":
		return false
	default:
		return false
	}
}
