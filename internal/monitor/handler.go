package monitor

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
)

const (
	StateChangeCloudWatchAlarm    string = "CloudWatch Alarm State Change"
	StateChangeCodePipelineAction string = "CodePipeline Action Execution State Change"
	StateChangeEC2Instance        string = "EC2 Instance State-change Notification"
)

func Handler(ctx context.Context, evt events.EventBridgeEvent, slackClient SlackClient) error {
	var err error

	slog.Info(fmt.Sprintf("event type: %v\n", evt.DetailType))
	slog.Info(fmt.Sprintf("event detail: %v\n", string(evt.Detail)))

	var event AwsEvent
	switch evt.DetailType {
	case StateChangeCloudWatchAlarm:
		event, err = NewCloudwatchAlarmEvent(evt)
	case StateChangeCodePipelineAction:
		event, err = NewCodePipelineEvent(evt)
	case StateChangeEC2Instance:
		event, err = NewEc2InstanceEvent(evt)
	default:
		return fmt.Errorf("error event type not recognized: %+v", evt.DetailType)
	}

	if err != nil {
		return fmt.Errorf("error parsing event: %+v", err)
	}

	if !event.ReportEvent(ctx) {
		slog.Info("event ignored")
		return nil
	}

	var msg *SlackMessage
	msg, err = event.CreateSlackMessage(ctx)
	if err != nil {
		return fmt.Errorf("error create message: %+v", err)
	}

	err = slackClient.PostMessage(ctx, msg)
	if err != nil {
		return fmt.Errorf("error post message: %+v", err)
	}

	slog.Info("success posted to Slack")

	return nil
}
