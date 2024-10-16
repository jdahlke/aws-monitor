package monitor_test

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jdahlke/aws-monitor/internal/monitor"
	"github.com/stretchr/testify/require"
)

func Test_Ec2InstanceEvent_CreateSlackMessage(t *testing.T) {
	tests := []struct {
		evt    events.EventBridgeEvent
		expect monitor.SlackMessage
	}{
		{
			evt: events.EventBridgeEvent{
				ID:         "EC2-event-id",
				DetailType: "EC2 Instance State-change Notification",
				Source:     "aws.ec2",
				AccountID:  "123456789012",
				Time:       time.Date(2017, 0o4, 22, 3, 31, 47, 0, time.UTC),
				Region:     "us-east-1",
				Resources: []string{
					"arn:aws:ec2:us-east-1:123456789012:instance/i-abcd1111",
				},
				Detail: []byte(`
				{
					"instance-id": "i-abcd1111",
					"state": "stopped"
				 }`),
			},
			expect: monitor.SlackMessage{
				Details: map[string]string{
					"*EC2 Instance i-abcd1111: stopped*": "<https://us-east-1.console.aws.amazon.com/ec2/v2/home#InstanceDetails:instanceId=i-abcd1111|view details>",
				},
				Severity: "danger",
			},
		},
	}

	for _, tc := range tests {
		var err error
		ctx := context.TODO()

		event, err := monitor.NewEc2InstanceEvent(tc.evt)
		require.NoError(t, err)

		actual, err := event.CreateSlackMessage(ctx)
		require.NoError(t, err)

		require.Equal(t, &tc.expect, actual)
	}
}
