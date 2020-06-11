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

// HandlerService provides all of the
type HandlerService struct {
	processor Processor
	handler   *http.ServeMux
}

// NewHandlerService creates a new ServeMux instance configured with the handler endpoints
func NewHandlerService(p Processor) *HandlerService {
	service := &HandlerService{processor: p}
	handler := http.NewServeMux()
	handler.HandleFunc("/healthz", service.handleHealthz)
	handler.HandleFunc("/webhook", service.handleWebhook)
	service.handler = handler
	return service
}

// Handler returns the http.ServeMux instance configured for the handler
func (svc *HandlerService) Handler() *http.ServeMux {
	return svc.handler
}

func (svc *HandlerService) handleWebhook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	status := http.StatusOK
	resp := &response{
		Status:  status,
		Message: "success",
	}
	if err := svc.processor.Process(r.Body); err != nil {
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

func (svc *HandlerService) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Ok!")
}
