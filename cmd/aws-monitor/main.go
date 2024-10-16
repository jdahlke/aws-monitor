package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/jdahlke/aws-monitor/internal/monitor"
	"github.com/jdahlke/aws-monitor/internal/monitor/slack"
)

func main() {
	var c monitor.SlackClient
	var err error

	ctx := context.Background()

	_, err = config.LoadDefaultConfig(ctx)
	if err != nil {
		slog.Error(fmt.Sprintf("unable to load AWS SDK config, %v", err))
		os.Exit(1)
	}

	c, err = makeSlackClient()
	if err != nil {
		slog.Error(fmt.Sprintf("unable create Slack client, %v", err))
		os.Exit(1)
	}

	slog.Info("starting lambda")
	lambda.Start(func(evt events.EventBridgeEvent) error {
		err = monitor.Handler(context.TODO(), evt, c)
		if err != nil {
			slog.Error(fmt.Sprintf("%v", err))
			return err
		}
		return nil
	})
	slog.Info("exiting lambda")
}

func makeSlackClient() (monitor.SlackClient, error) {
	return slack.NewClient(
		os.Getenv("SLACK_URL"),
		os.Getenv("SLACK_CHANNEL"),
		&http.Client{},
	)
}
