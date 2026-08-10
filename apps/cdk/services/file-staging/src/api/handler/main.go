package main

import (
	"contracts/dist/filestaging/v1/filestagingconnect"
	"log"
	"net/http"

	"connectrpc.com/connect"
)

type Handler struct{}

func main() {
	mux := http.NewServeMux()

	mux.Handle(filestagingconnect.NewInternalHandler(
		&Handler{},
		connect.WithInterceptors(MonitoringInterceptor, ValidationInterceptor, AuthInterceptor)))

	mux.Handle(filestagingconnect.NewExternalHandler(
		&Handler{},
		connect.WithInterceptors(MonitoringInterceptor, ValidationInterceptor, AuthInterceptor)))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(srv.ListenAndServe())
}
