package monitor_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jdahlke/aws-monitor/internal/monitor"
	"github.com/jdahlke/aws-monitor/internal/monitor/mocks"
	"github.com/jdahlke/aws-monitor/internal/monitor/slack"
	"github.com/stretchr/testify/require"
)

func Test_Handler(t *testing.T) {
	tests := []struct {
		evt             events.EventBridgeEvent
		err             error
		httpPostInvoked bool
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
			httpPostInvoked: true,
		},

		// event build fails because of unknown AWS event
		{
			evt: events.EventBridgeEvent{
				ID:         "CWA-event-id",
				DetailType: "Unknown AWS Event",
				Source:     "aws.cloudwatch",
				AccountID:  "123456789012",
				Time:       time.Date(2017, 0o4, 22, 3, 31, 47, 0, time.UTC),
				Region:     "us-east-1",
				Resources: []string{
					"arn:aws:cloudwatch:us-east-1:123456789012:alarm:ServerCpuTooHigh",
				},
				Detail: []byte(`
				{
					"instance-id": "i-abcd1111",
					"state": "stopped"
				 }`),
			},
			err: fmt.Errorf("error event type not recognized: Unknown AWS Event"),
		},

		// event ignored because of Code Pipeline state succeeded
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
		},
	}

	slackUrl := "http://www.example.com"
	slackChannel := "ABCDE"
	httpClient := mocks.MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader("success")),
			}, nil
		},
	}

	slackClient, _ := slack.NewClient(slackUrl, slackChannel, &httpClient)

	for _, tc := range tests {
		ctx := context.TODO()

		httpClient.DoFuncInvoked = false
		err := monitor.Handler(ctx, tc.evt, slackClient)

		require.Equal(t, tc.err, err)
		require.Equal(t, tc.httpPostInvoked, httpClient.DoFuncInvoked)
	}
}
