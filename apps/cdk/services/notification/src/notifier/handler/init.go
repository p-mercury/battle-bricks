package main

import (
	"context"
	"contracts/dist/identity/v1/identityconnect"
	"dynamokit/dynamolease"
	"log"
	"log/slog"
	"net/http"
	"os"
	"sigv4http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-xray-sdk-go/instrumentation/awsv2"
	"github.com/aws/aws-xray-sdk-go/xray"
)

var (
	Logger      *slog.Logger
	Ses         *ses.Client
	DynamoRead  *dynamodb.Client
	DynamoWrite *dynamodb.Client
	Leaser      *dynamolease.Client
	Http        *http.Client

	IdentityService identityconnect.InternalClient

	Stage     string
	StackName string
	Namespace string
	Hostname  string
	ApiUrl    string
	TableName *string

	InfoSesTemplateName            *string
	SingleActionSesTemplateName    *string
	PurchasingOrderSesTemplateName *string
)

func init() {
	Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	Stage = os.Getenv("STAGE")
	StackName = os.Getenv("STACK_NAME")
	Namespace = os.Getenv("NAMESPACE")
	Hostname = os.Getenv("HOSTNAME")
	ApiUrl = os.Getenv("API_URL")
	TableName = new(os.Getenv("TABLE_NAME"))

	InfoSesTemplateName = new(os.Getenv("INFO_SES_TEMPLATE_NAME"))
	SingleActionSesTemplateName = new(os.Getenv("SINGLE_ACTION_SES_TEMPLATE_NAME"))
	PurchasingOrderSesTemplateName = new(os.Getenv("PURCHASING_ORDER_SES_TEMPLATE_NAME"))

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
	Ses = ses.NewFromConfig(cfg)
	Leaser = dynamolease.NewClient(DynamoWrite, *TableName, "pk", nil, new("ttl"))
	Http = xray.Client(nil)

	{
		signer := sigv4http.NewSigner(cfg, Http.Transport)
		IdentityService = identityconnect.NewInternalClient(&http.Client{Transport: signer}, ApiUrl)
	}
}
