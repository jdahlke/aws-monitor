package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type cloudwatchAlarmEvent struct {
	event  events.EventBridgeEvent
	detail cloudwatchAlarmDetail
}

type cloudwatchAlarmDetail struct {
	AlarmName     string                       `json:"alarmName,omitempty"`
	Configuration cloudwatchAlarmConfiguration `json:"configuration,omitempty"`
	State         cloudwatchAlarmState         `json:"state,omitempty"`
}

type cloudwatchAlarmConfiguration struct {
	Description string `json:"description,omitempty"`
}

type cloudwatchAlarmState struct {
	Value string `json:"value,omitempty"`
}

// NewEc2InstanceEvent returns new Ec2InstanceEvent
func NewCloudwatchAlarmEvent(evt events.EventBridgeEvent) (*cloudwatchAlarmEvent, error) {
	var detail cloudwatchAlarmDetail
	if err := json.Unmarshal(evt.Detail, &detail); err != nil {
		return nil, err
	}

	return &cloudwatchAlarmEvent{
		event:  evt,
		detail: detail,
	}, nil
}

// CreateSlackMessage creates new SlackMessage for CloudwatchAlarmEvent state change event
func (e cloudwatchAlarmEvent) CreateSlackMessage(ctx context.Context) (*SlackMessage, error) {
	var title, link, severity string

	title = fmt.Sprintf("*Alarm %v changed to %v*", e.detail.AlarmName, e.detail.State.Value)
	link = fmt.Sprintf("<https://%v.console.aws.amazon.com/cloudwatch/home#alarmsV2:alarm/%v|view details>", e.event.Region, e.detail.AlarmName)

	switch e.detail.State.Value {
	case "ALARM":
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

func (e cloudwatchAlarmEvent) ReportEvent(ctx context.Context) bool {
	desc := e.detail.Configuration.Description
	return !strings.Contains(desc, "TargetTrackingScaling policy")
}
