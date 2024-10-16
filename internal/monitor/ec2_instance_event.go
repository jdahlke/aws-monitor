package monitor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

type ec2InstanceEvent struct {
	event  events.EventBridgeEvent
	detail ec2InstanceDetail
}

type ec2InstanceDetail struct {
	InstanceId string `json:"instance-id,omitempty"`
	State      string `json:"state,omitempty"`
}

// NewEc2InstanceEvent returns new Ec2InstanceEvent
func NewEc2InstanceEvent(evt events.EventBridgeEvent) (*ec2InstanceEvent, error) {
	var detail ec2InstanceDetail
	if err := json.Unmarshal(evt.Detail, &detail); err != nil {
		return nil, err
	}

	return &ec2InstanceEvent{
		event:  evt,
		detail: detail,
	}, nil
}

// CreateSlackMessage creates new SlackMessage for Ec2InstanceEvent state change event
func (e ec2InstanceEvent) CreateSlackMessage(ctx context.Context) (*SlackMessage, error) {
	var title, link, severity string

	title = fmt.Sprintf("*EC2 Instance %v: %v*", e.detail.InstanceId, e.detail.State)
	link = fmt.Sprintf("<https://%v.console.aws.amazon.com/ec2/v2/home#InstanceDetails:instanceId=%v|view details>", e.event.Region, e.detail.InstanceId)

	switch e.detail.State {
	case "stopping":
	case "stopped":
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

func (e ec2InstanceEvent) ReportEvent(ctx context.Context) bool {
	return true
}
