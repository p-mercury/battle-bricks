package main

import (
	"log"
	"net/http"

	"contracts/dist/catalogue/v1/catalogueconnect"

	"connectrpc.com/connect"
)

type Handler struct{}

func main() {
	mux := http.NewServeMux()

	mux.Handle(catalogueconnect.NewInternalHandler(
		&Handler{},
		connect.WithInterceptors(MonitoringInterceptor, ValidationInterceptor, AuthInterceptor, IdempotencyInterceptor)))

	mux.Handle(catalogueconnect.NewExternalHandler(
		&Handler{},
		connect.WithInterceptors(MonitoringInterceptor, ValidationInterceptor, AuthInterceptor, IdempotencyInterceptor)))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(srv.ListenAndServe())
}
