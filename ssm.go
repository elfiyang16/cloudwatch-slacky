package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"os"
	"strings"
)

const (
	EnvKeyWebhookURL = "SLACK_WEBHOOK_URL"
)

type SSMAdapter interface {
	get(string) (string, error)
}

func NewSSMAdapter() (SSMAdapter, error) {
	p := os.Getenv(EnvKeyWebhookURL)
	if strings.HasPrefix(p, "https://hooks.slack.com") {
		return NewEnvStore()
	}
	if strings.HasPrefix(p, "ssm:") {
		key := strings.Replace(p, "ssm:", "", 1)
		return NewParameterStore(key)
	}
	return nil, fmt.Errorf("no ssm adapter found, %s. please set %s", p, EnvKeyWebhookURL)
}

// Use SSM parameter store
type ParamStore struct {
	svc        *ssm.SSM
	webhookKey string
}

func NewParameterStore(key string) (ParamStore, error) {
	sess, err := session.NewSessionWithOptions(session.Options{})
	if err != nil {
		return ParamStore{}, err
	}
	svc := ssm.New(sess)
	return ParamStore{
		svc:        svc,
		webhookKey: key,
	}, nil
}

func (p ParamStore) get(key string) (string, error) {
	if key == EnvKeyWebhookURL {
		key = p.webhookKey
	}
	res, err := p.svc.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", nil
	}
	return *res.Parameter.Value, nil
}

// Use env var
type EnvStore struct {
}

func NewEnvStore() (EnvStore, error) {
	return EnvStore{}, nil
}

func (p EnvStore) get(key string) (string, error) {
	return os.Getenv(key), nil
}
