package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type response struct {
	Status  int
	Message string
}

// GenerateMux generates an *http.ServeMux multiplexer pre-configured
// with receiver routes
func GenerateMux(p Processor) *http.ServeMux {
	handler := http.NewServeMux()
	handler.HandleFunc("/healthz", buildHealthzHandler())
	handler.HandleFunc("/webhook", buildWebhookHandler(p))
	return handler
}

func buildWebhookHandler(processor Processor) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		status := http.StatusOK
		resp := &response{
			Status:  status,
			Message: "success",
		}
		if err := processor.Process(r.Body); err != nil {
			status = http.StatusBadRequest
			resp = &response{
				Status:  status,
				Message: err.Error(),
			}
		}
		jsonBytes, jsonErr := json.Marshal(resp)
		if jsonErr != nil {
			log.WithFields(log.Fields{
				"error":    jsonErr,
				"response": resp,
			}).Error("Failed to marshal response to JSON")
			jsonBytes = []byte("{ \"Status\": 200 }")
		}
		w.WriteHeader(status)
		fmt.Fprint(w, string(jsonBytes[:]))
	}
}

func buildHealthzHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Ok!")
	}
}
