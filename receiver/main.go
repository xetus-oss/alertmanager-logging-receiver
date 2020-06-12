package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	// TODO: figure out how to avoid using a global logger while
	// still configuring JSON output and static fields
	log.SetFormatter(&log.JSONFormatter{})

	listenAddress := ":8080"
	if os.Getenv("PORT") != "" {
		listenAddress = ":" + os.Getenv("PORT")
	}

	log.Info(fmt.Sprintf("listening on: %v", listenAddress))
	server := &http.Server{
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		Handler:      GenerateMux(NewLoggingProcessor()),
		Addr:         listenAddress,
	}
	log.Fatal(server.ListenAndServe())
}
