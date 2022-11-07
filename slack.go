package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

const DefaultSlackTimeout = 5 * time.Second

type SlackClient struct {
	WebHookUrl string
	UserName   string
	Channel    string
	Timeout    time.Duration
}

func NewSlackClient(username string, adapter SSMAdapter) (SlackClient, error) {
	webHookUrl, err := adapter.get(EnvKeyWebhookURL)
	if err != nil {
		return SlackClient{}, fmt.Errorf("can not get webHookUrl, %w", err)
	}

	return SlackClient{
		WebHookUrl: webHookUrl,
		UserName:   username,
		Timeout:    DefaultSlackTimeout,
	}, nil
}

func (sc SlackClient) sendHttpRequest(body string) error {
	req, err := http.NewRequest(
		http.MethodPost,
		sc.WebHookUrl,
		bytes.NewBuffer([]byte(body)),
	)
	if err != nil {
		return fmt.Errorf("error with sendHttpRequest, %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	if sc.Timeout == 0 {
		sc.Timeout = DefaultSlackTimeout
	}
	client := &http.Client{Timeout: sc.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return err
	}

	if buf.String() != "ok" {
		return fmt.Errorf("error returned from Slack, %s", buf.String())

	}
	return nil
}
