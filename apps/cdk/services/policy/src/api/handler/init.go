package main

import (
	"bytes"
	"connectkit"
	"context"
	"contracts/dist/policy/v1/policyconnect"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"sigv4http"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-xray-sdk-go/instrumentation/awsv2"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/open-policy-agent/opa/v1/sdk"
)

var (
	Logger      *slog.Logger
	DynamoRead  *dynamodb.Client
	DynamoWrite *dynamodb.Client
	Opa         *sdk.OPA
	Http        *http.Client

	PolicyService policyconnect.InternalClient

	StackName          string
	Namespace          string
	ApiUrl             string
	TableName          *string
	EventBusName       *string
	EventBusEndpointId *string

	MonitoringInterceptor connect.UnaryInterceptorFunc
	AuthInterceptor       connect.UnaryInterceptorFunc
	ValidationInterceptor *validate.Interceptor
)

func init() {
	Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	StackName = os.Getenv("STACK_NAME")
	Namespace = os.Getenv("NAMESPACE")
	ApiUrl = os.Getenv("API_URL")
	TableName = new(os.Getenv("TABLE_NAME"))
	EventBusName = new(os.Getenv("EVENT_BUS_NAME"))
	EventBusEndpointId = new(os.Getenv("EVENT_BUS_ENDPOINT_ID"))

	{
		cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(os.Getenv("TABLE_WRITE_REGION")))
		if err != nil {
			log.Fatalf("Unable to configure aws services, %v", err)
			return
		}

		awsv2.AWSV2Instrumentor(&cfg.APIOptions)

		DynamoWrite = dynamodb.NewFromConfig(cfg)
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		Logger.Error("Unable to configure aws services", slog.Any("error", err))
		log.Fatal()
	}

	awsv2.AWSV2Instrumentor(&cfg.APIOptions)

	DynamoRead = dynamodb.NewFromConfig(cfg)
	Http = xray.Client(nil)

	{
		signer := sigv4http.NewSigner(cfg, Http.Transport)
		PolicyService = policyconnect.NewInternalClient(&http.Client{Transport: signer}, ApiUrl)
	}

	MonitoringInterceptor = connectkit.NewMonitoringInterceptor(connectkit.NewMonitoringInterceptorInput{
		Logger: Logger,
		Metadata: map[string]any{
			"stackName": StackName,
		},
	})
	AuthInterceptor = connectkit.NewAuthInterceptor(PolicyService)
	ValidationInterceptor = validate.NewInterceptor()

	var config []byte
	config = fmt.Appendf(config, `
    services:
      s3:
        url: %q
        credentials:
          s3_signing:
            service: "s3"
            signature_Version: "4"
            environment_credentials: {}
    bundles:
      bundle:
        service: s3
        resource: "bundle.tar.gz"
  `, os.Getenv("BUCKET_URL"))

	opa, err := sdk.New(context.Background(), sdk.Options{
		Config: bytes.NewReader(config),
	})
	if err != nil {
		Logger.Error("Error creting OPA client", slog.Any("error", err))
		log.Fatal()
	}
	Opa = opa
}
