package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-xray-sdk-go/instrumentation/awsv2"
)

var (
	Logger   *slog.Logger
	Dynamo   *dynamodb.Client
	S3       *s3.Client
	EventBus *eventbridge.Client

	StackName          string
	Namespace          string
	TableName          *string
	EventBusName       *string
	EventBusEndpointId *string
	BucketName         *string
)

func init() {
	Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	StackName = os.Getenv("STACK_NAME")
	Namespace = os.Getenv("NAMESPACE")
	TableName = new(os.Getenv("TABLE_NAME"))
	EventBusName = new(os.Getenv("EVENT_BUS_NAME"))
	EventBusEndpointId = new(os.Getenv("EVENT_BUS_ENDPOINT_ID"))
	BucketName = new(os.Getenv("BUCKET_NAME"))

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		log.Fatalf("Unable to configure aws services, %v", err)
	}

	awsv2.AWSV2Instrumentor(&cfg.APIOptions)

	Dynamo = dynamodb.NewFromConfig(cfg)
	S3 = s3.NewFromConfig(cfg)
	EventBus = eventbridge.NewFromConfig(cfg)
}
