package main

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/prometheus/alertmanager/template"
	log "github.com/sirupsen/logrus"
)

// Processor performs whatever actions are needed in response to
// a received webhook request body
type Processor interface {
	Process(request io.Reader) error
}

// LoggingProcessor is a webhook request processor that logs received webhook alerts
type LoggingProcessor struct{}

// NewLoggingProcessor returns a new LoggingProcessor instance.
func NewLoggingProcessor() *LoggingProcessor {
	return &LoggingProcessor{}
}

// Process the supplied request, in this case by logging
// at info level if an alert was received and at error level
// if the received request was invalid.
func (p *LoggingProcessor) Process(request io.Reader) error {
	// Godoc: https://godoc.org/github.com/prometheus/alertmanager/template#Data
	data := template.Data{}
	err := json.NewDecoder(request).Decode(&data)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to decode webhook request JSON")
		return errors.New("Failed to decode webhook request")
	}

	if len(data.Alerts) < 1 {
		log.WithFields(log.Fields{
			"data": data,
		}).Error("Webhook request JSON did not contain any alerts")
		return errors.New("Webhook request must include at least one Alert")
	}

	log.WithFields(log.Fields{
		"request": data,
	}).Info("Received webhook request")
	return nil
}
