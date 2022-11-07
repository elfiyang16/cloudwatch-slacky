package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
)

const (
	EnvKeySlackName   = "SLACK_NAME"
	EnvKeyTemplateDir = "SLACK_TEMPLATE_DIR"
)

func Handler(context context.Context, event events.CloudWatchEvent) error {
	// Debug.
	// j, _ := json.MarshalIndent(event, "", "  ")
	// fmt.Printf("Source = %s\n", string(j))

	msgBody, err := GenerateMessage(event.Source, event.DetailType, event.Detail)
	if err != nil {
		return fmt.Errorf("failed to generate message %w", err)
	}
	adapter, err := NewSSMAdapter()
	if err != nil {
		return fmt.Errorf("failed to init adapter, %w", err)
	}
	ss, err := NewSlackClient(getSlackName(), adapter)
	if err != nil {
		return fmt.Errorf("failed to init new slack client, %w", err)
	}

	if err := ss.sendHttpRequest(msgBody); err != nil {
		return fmt.Errorf("failed to send sendHttpRequest, %w", err)
	}

	return nil
}

func getSlackName() string {
	p := os.Getenv(EnvKeySlackName)
	if p == "" {
		return "Alert"
	}
	return p
}

func main() {
	lambda.Start(Handler())
}
