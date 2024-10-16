package monitor_test

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jdahlke/aws-monitor/internal/monitor"
	"github.com/stretchr/testify/require"
)

func Test_CloudwatchAlarmEvent_CreateSlackMessage(t *testing.T) {
	tests := []struct {
		evt    events.EventBridgeEvent
		expect monitor.SlackMessage
	}{
		{
			evt: events.EventBridgeEvent{
				ID:         "CWA-event-id",
				DetailType: "CloudWatch Alarm State Change",
				Source:     "aws.cloudwatch",
				AccountID:  "123456789012",
				Time:       time.Date(2017, 0o4, 22, 3, 31, 47, 0, time.UTC),
				Region:     "us-east-1",
				Resources: []string{
					"arn:aws:cloudwatch:us-east-1:123456789012:alarm:ServerCpuTooHigh",
				},
				Detail: []byte(`
				{
					"alarmName": "ServerCpuTooHigh",
					"configuration": {
						"description": "Goes into alarm when server CPU utilization is too high!"
					},
					"state": {
						"value": "ALARM"
					}
				 }`),
			},
			expect: monitor.SlackMessage{
				Details: map[string]string{
					"*Alarm ServerCpuTooHigh changed to ALARM*": "<https://us-east-1.console.aws.amazon.com/cloudwatch/home#alarmsV2:alarm/ServerCpuTooHigh|view details>",
				},
				Severity: "danger",
			},
		},
	}

	for _, tc := range tests {
		var err error
		ctx := context.TODO()

		event, err := monitor.NewCloudwatchAlarmEvent(tc.evt)
		require.NoError(t, err)

		actual, err := event.CreateSlackMessage(ctx)
		require.NoError(t, err)

		require.Equal(t, &tc.expect, actual)
	}
}
