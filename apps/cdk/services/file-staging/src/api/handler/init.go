package main

import (
	"connectkit"
	"context"
	"contracts/dist/policy/v1/policyconnect"
	"log"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"sigv4http"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-xray-sdk-go/instrumentation/awsv2"
	"github.com/aws/aws-xray-sdk-go/xray"
)

var (
	Logger      *slog.Logger
	DynamoRead  *dynamodb.Client
	DynamoWrite *dynamodb.Client
	EventBus    *eventbridge.Client
	S3          *s3.Client
	S3Presign   *s3.PresignClient
	Http        *http.Client

	PolicyService policyconnect.InternalClient

	StackName          string
	Namespace          string
	ApiUrl             string
	TableName          *string
	EventBusName       *string
	EventBusEndpointId *string
	BucketName         *string

	MonitoringInterceptor connect.UnaryInterceptorFunc
	AuthInterceptor       connect.UnaryInterceptorFunc
	ValidationInterceptor *validate.Interceptor
)

func init() {
	Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	mime.AddExtensionType(".pdf", "application/pdf")
	mime.AddExtensionType(".zip", "application/zip")
	mime.AddExtensionType(".xml", "application/xml")
	mime.AddExtensionType(".txt", "text/plain")
	mime.AddExtensionType(".png", "image/png")
	mime.AddExtensionType(".gif", "image/gif")
	mime.AddExtensionType(".jpeg", "image/jpeg")
	mime.AddExtensionType(".jpg", "image/jpeg")
	mime.AddExtensionType(".csv", "text/csv")
	mime.AddExtensionType(".doc", "application/msword")
	mime.AddExtensionType(".docx", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")

	StackName = os.Getenv("STACK_NAME")
	Namespace = os.Getenv("NAMESPACE")
	ApiUrl = os.Getenv("API_URL")
	TableName = new(os.Getenv("TABLE_NAME"))
	EventBusName = new(os.Getenv("EVENT_BUS_NAME"))
	EventBusEndpointId = new(os.Getenv("EVENT_BUS_ENDPOINT_ID"))
	BucketName = new(os.Getenv("BUCKET_NAME"))

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
		log.Fatalf("Unable to configure aws services, %v", err)
		return
	}

	awsv2.AWSV2Instrumentor(&cfg.APIOptions)

	DynamoRead = dynamodb.NewFromConfig(cfg)
	EventBus = eventbridge.NewFromConfig(cfg)
	S3 = s3.NewFromConfig(cfg)
	S3Presign = s3.NewPresignClient(S3)
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
}
