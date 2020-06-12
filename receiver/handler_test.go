package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type mockProcessor struct {
	err       string
	processed bool
}

func (p *mockProcessor) Process(request io.Reader) error {
	p.processed = true
	if p.err != "" {
		log.Info(fmt.Sprintf("MockProcessor is returning an error: %v", p.err))
		return errors.New(p.err)
	}
	log.Info("MockProcessor is not returning an error")
	return nil
}

func TestHandlerServiceRegistersExpectedHandlersAtExpectedPaths(t *testing.T) {
	tests := map[string]struct {
		path          string
		method        string
		processor     *mockProcessor
		wantStatus    int
		wantText      string
		wantProcessed bool
	}{
		"webhook processor error should return error": {
			path: "/webhook", method: "POST", processor: &mockProcessor{err: "mock error"},
			wantStatus: 400, wantProcessed: true, wantText: "{\"Status\":400,\"Message\":\"mock error\"}",
		},
		"webhook processor success should return ok": {
			path: "/webhook", method: "POST", processor: &mockProcessor{},
			wantStatus: 200, wantProcessed: true, wantText: "{\"Status\":200,\"Message\":\"success\"}",
		},
		"healthz endpoint should return ok": {
			path: "/healthz", method: "GET", processor: &mockProcessor{},
			wantStatus: 200, wantProcessed: false, wantText: "Ok!",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mux := GenerateMux(tc.processor)
			request := httptest.NewRequest(tc.method, tc.path, nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, request)
			response := w.Result()

			assert.Equal(t, tc.wantStatus, response.StatusCode)
			assert.Equal(t, tc.wantProcessed, tc.processor.processed)

			defer response.Body.Close()
			gotBytes, _ := ioutil.ReadAll(response.Body)
			assert.Equal(t, tc.wantText, string(gotBytes))
		})
	}
}
