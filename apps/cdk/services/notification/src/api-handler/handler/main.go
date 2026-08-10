package main

import (
	"contracts/dist/notification/v1/notificationconnect"
	"log"
	"net/http"

	"connectrpc.com/connect"
)

type Handler struct{}

func main() {
	mux := http.NewServeMux()

	mux.Handle(notificationconnect.NewInternalHandler(
		&Handler{},
		connect.WithInterceptors(MonitoringInterceptor, ValidationInterceptor, AuthInterceptor, IdempotencyInterceptor)))

	mux.Handle(notificationconnect.NewExternalHandler(
		&Handler{},
		connect.WithInterceptors(MonitoringInterceptor, ValidationInterceptor, AuthInterceptor, IdempotencyInterceptor)))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(srv.ListenAndServe())
}
