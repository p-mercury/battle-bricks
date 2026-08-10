package main

import (
	"log"
	"net/http"

	"contracts/dist/identity/v1/identityconnect"

	"connectrpc.com/connect"
)

type Handler struct{}

func main() {
	mux := http.NewServeMux()

	mux.Handle(identityconnect.NewInternalHandler(
		&Handler{},
		connect.WithInterceptors(MonitoringInterceptor, ValidationInterceptor, AuthInterceptor, IdempotencyInterceptor)))

	mux.Handle(identityconnect.NewExternalHandler(
		&Handler{},
		connect.WithInterceptors(MonitoringInterceptor, ValidationInterceptor, AuthInterceptor, IdempotencyInterceptor)))

	mux.Handle("/"+identityconnect.ExternalName+"/auth/{path...}", AuthHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(srv.ListenAndServe())
}
