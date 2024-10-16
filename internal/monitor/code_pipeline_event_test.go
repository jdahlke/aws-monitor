package monitor_test

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jdahlke/aws-monitor/internal/monitor"
	"github.com/stretchr/testify/require"
)

func Test_CodePipelineEvent_CreateSlackMessage(t *testing.T) {
	tests := []struct {
		evt    events.EventBridgeEvent
		expect monitor.SlackMessage
	}{
		{
			evt: events.EventBridgeEvent{
				ID:         "CWE-event-id",
				DetailType: "CodePipeline Action Execution State Change",
				Source:     "aws.codepipeline",
				AccountID:  "123456789012",
				Time:       time.Date(2017, 0o4, 22, 3, 31, 47, 0, time.UTC),
				Region:     "us-east-1",
				Resources: []string{
					"arn:aws:codepipeline:us-east-1:123456789012:pipeline:myPipeline",
				},
				Detail: []byte(`
				{
					"pipeline": "myPipeline",
					"execution-id": "01234567-0123-0123-0123-012345678901",
					"stage": "Prod",
					"action": "myAction",
					"state": "STARTED"
				 }`),
			},
			expect: monitor.SlackMessage{
				Details: map[string]string{
					"*Code Pipeline myPipeline: myAction started*": "<https://us-east-1.console.aws.amazon.com/codesuite/codepipeline/pipelines/myPipeline/executions/01234567-0123-0123-0123-012345678901/timeline|view details>",
				},
				Severity: "good",
			},
		},
	}

	for _, tc := range tests {
		var err error
		ctx := context.TODO()

		event, err := monitor.NewCodePipelineEvent(tc.evt)
		require.NoError(t, err)

		actual, err := event.CreateSlackMessage(ctx)
		require.NoError(t, err)

		require.Equal(t, &tc.expect, actual)
	}
}
