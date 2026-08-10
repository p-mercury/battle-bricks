package main

import (
	"context"
	"contracts/dist/identity/v1/identityconnect"
	"log"
	"log/slog"
	"net/http"
	"os"
	"sigv4http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-xray-sdk-go/instrumentation/awsv2"
	"github.com/aws/aws-xray-sdk-go/xray"
)

var (
	Logger *slog.Logger
	Http   *http.Client

	IdentityService identityconnect.InternalClient

	StackName string
	ApiUrl    string
)

func init() {
	Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	StackName = os.Getenv("STACK_NAME")
	ApiUrl = os.Getenv("API_URL")

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		log.Fatalf("Unable to configure aws services, %v", err)
		return
	}

	awsv2.AWSV2Instrumentor(&cfg.APIOptions)

	Http = xray.Client(nil)

	{
		signer := sigv4http.NewSigner(cfg, Http.Transport)
		IdentityService = identityconnect.NewInternalClient(&http.Client{Transport: signer}, ApiUrl)
	}
}
