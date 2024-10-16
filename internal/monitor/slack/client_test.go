package slack_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/jdahlke/aws-monitor/internal/monitor"
	"github.com/jdahlke/aws-monitor/internal/monitor/mocks"
	"github.com/jdahlke/aws-monitor/internal/monitor/slack"
	"github.com/stretchr/testify/require"
)

func Test_SlackClient_PostMessage(t *testing.T) {
	tests := []struct {
		message         monitor.SlackMessage
		requestBody     string
		requestUrl      string
		err             error
		httpPostInvoked bool
	}{
		{
			message: monitor.SlackMessage{
				Details: map[string]string{
					"*Alarm ServerCpuTooHigh changed to ALARM*": "<https://us-east-1.console.aws.amazon.com/cloudwatch/home#alarmsV2:alarm/ServerCpuTooHigh|view details>",
				},
				Severity: "danger",
			},
			err:             nil,
			httpPostInvoked: true,
		},
	}

	slackUrl := "http://www.example.com"
	slackChannel := "ABCDE"
	mockHttpClient := mocks.MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader("success")),
			}, nil
		},
	}

	c, _ := slack.NewClient(slackUrl, slackChannel, &mockHttpClient)

	for _, tc := range tests {
		ctx := context.TODO()

		err := c.PostMessage(ctx, &tc.message)
		require.NoError(t, err)

		require.Equal(t, tc.httpPostInvoked, mockHttpClient.DoFuncInvoked)
	}
}
