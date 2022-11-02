package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Event interface {
	readTemplate() (*template.Template, error)
	generateMessage(*template.Template) (string, error)
}

func GenerateMessage(source, detailType string, detail json.RawMessage) (string, error) {
	e, err := NewEvent(source, detailType, detail)
	if err != nil {
		return "", fmt.Errorf("NewEvent, %w", err)
	}
	tmpl, err := e.readTemplate()
	if err != nil {
		return "", fmt.Errorf("readTemplate, %w", err)
	}
	return e.generateMessage(tmpl)
}

func NewEvent(source, detailType string, detail json.RawMessage) (Event, error) {
	switch source {
	case "aws.cloudwatch":
		if detailType == "CloudWatch Alarm State Change" {
			return NewEventCloudWatchAlarm(source, detailType, detail), nil
		}
	}
	return nil, fmt.Errorf("cannot find matched template: %s, %s", source, detailType)
}

func readTemplate(paths ...string) (*template.Template, error) {
	dir := os.Getenv(EnvKeyTemplateDir)
	path := make([]string, 0, 10)

	if dir == "" {
		root := filepath.Dir(os.Args[0])
		path = []string{root, "templates"}
	} else {
		path = []string{dir}
	}
	for _, p := range paths {
		path = append(path, p)
	}
	p := strings.Join(path, "/")
	return template.ParseFiles(p)
}
