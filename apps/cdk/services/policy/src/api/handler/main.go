package main

import (
	"log"
	"net/http"

	"contracts/dist/policy/v1/policyconnect"

	"connectrpc.com/connect"
)

type Handler struct{}

func main() {
	mux := http.NewServeMux()

	mux.Handle(policyconnect.NewInternalHandler(
		&Handler{},
		connect.WithInterceptors(MonitoringInterceptor, ValidationInterceptor, AuthInterceptor)))

	mux.Handle(policyconnect.NewExternalHandler(
		&Handler{},
		connect.WithInterceptors(MonitoringInterceptor, ValidationInterceptor, AuthInterceptor)))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(srv.ListenAndServe())
}
