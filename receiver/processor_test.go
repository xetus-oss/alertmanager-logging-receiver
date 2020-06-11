package main

import (
	"encoding/json"
	"testing"

	"github.com/prometheus/alertmanager/template"
	"github.com/sirupsen/logrus"
	logrusTest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestProcessorLogsOnValidRequest(t *testing.T) {
	// register a hook with the global logrus logger
	// TODO: figure out how to NOT use a global logger!
	hook := logrusTest.NewGlobal()

	request := "{ \"alerts\": [{ \"fake\": \"fake\" }] }"
	err := NewLoggingProcessor().Process(readCloser(request))
	if err != nil {
		assert.Fail(t, "Unexpected error processing valid request: %s", err)
	}

	expectedData := template.Data{}
	json.NewDecoder(readCloser(request)).Decode(&expectedData)
	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.InfoLevel, hook.LastEntry().Level)
	assert.Equal(t, "Received webhook request", hook.LastEntry().Message)
	assert.Equal(t, expectedData, hook.LastEntry().Data["request"])
}

func TestProcessLogsErrorOnInvalidRequestJSON(t *testing.T) {
	hook := logrusTest.NewGlobal()

	request := "I am not valid json"
	err := NewLoggingProcessor().Process(readCloser(request))
	if err == nil {
		assert.Fail(t, "Expected error while processing invalid request: %s", request)
	}

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.NotNil(t, hook.LastEntry().Data["error"])
}

func TestProcessorLogsErrorOnEmptyAlertsJSON(t *testing.T) {
	hook := logrusTest.NewGlobal()

	request := "{}"
	err := NewLoggingProcessor().Process(readCloser(request))
	if err == nil {
		assert.Fail(t, "Expected error while processing request without alerts: %s", err)
	}

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.NotNil(t, hook.LastEntry().Data["data"])
}
