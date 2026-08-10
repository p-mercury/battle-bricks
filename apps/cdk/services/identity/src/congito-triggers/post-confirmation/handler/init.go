package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-xray-sdk-go/instrumentation/awsv2"
	"github.com/aws/aws-xray-sdk-go/strategy/ctxmissing"
	"github.com/aws/aws-xray-sdk-go/xray"
)

var (
	Logger   *slog.Logger
	Dynamo   *dynamodb.Client
	EventBus *eventbridge.Client
	Http     *http.Client

	StackName    string
	ServiceName  string
	Namespace    *string
	TableName    *string
	EventBusName *string
)

func init() {
	xray.Configure(xray.Config{
		ContextMissingStrategy: ctxmissing.NewDefaultIgnoreErrorStrategy(),
	})

	Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	StackName = os.Getenv("STACK_NAME")
	ServiceName = os.Getenv("SERVICE_NAME")
	Namespace = aws.String(os.Getenv("NAMESPACE"))
	TableName = aws.String(os.Getenv("TABLE_NAME"))
	EventBusName = aws.String(os.Getenv("EVENT_BUS_NAME"))

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		log.Fatalf("Unable to configure aws services, %v", err)
		return
	}

	awsv2.AWSV2Instrumentor(&cfg.APIOptions)

	Dynamo = dynamodb.NewFromConfig(cfg)
	EventBus = eventbridge.NewFromConfig(cfg)
	Http = xray.Client(nil)
}
