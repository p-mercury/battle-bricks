package main

import (
	"contracts/dist/policy/v1/policyconnect"
	"connectkit"
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"sigv4http"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-xray-sdk-go/instrumentation/awsv2"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/opensearch-project/opensearch-go/v4"
	requestsigner "github.com/opensearch-project/opensearch-go/v4/signer/awsv2"
)

var (
	Logger      *slog.Logger
	DynamoRead  *dynamodb.Client
	DynamoWrite *dynamodb.Client
	EventBus    *eventbridge.Client
	OpenSearch  *opensearch.Client
	Http        *http.Client

	PolicyService policyconnect.InternalClient

	StackName          string
	Namespace          string
	ApiUrl             string
	TableName          *string
	EventBusName       *string
	EventBusEndpointId *string

	MonitoringInterceptor  connect.UnaryInterceptorFunc
	ValidationInterceptor  *validate.Interceptor
	AuthInterceptor        connect.UnaryInterceptorFunc
	IdempotencyInterceptor connect.UnaryInterceptorFunc
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
		log.Fatalf("Unable to configure aws services, %v", err)
		return
	}

	awsv2.AWSV2Instrumentor(&cfg.APIOptions)

	DynamoRead = dynamodb.NewFromConfig(cfg)
	EventBus = eventbridge.NewFromConfig(cfg)
	Http = xray.Client(nil)

	{
		stsClient := sts.NewFromConfig(cfg)
		assumeRoleProvider := stscreds.NewAssumeRoleProvider(stsClient, os.Getenv("OPEN_SEARCH_ROLE_ARN"))

		assumedCfg := cfg
		assumedCfg.Credentials = aws.NewCredentialsCache(assumeRoleProvider)

		signer, err := requestsigner.NewSignerWithService(assumedCfg, "es")
		if err != nil {
			Logger.Error("Error creating request signer", slog.Any("error", err))
			log.Fatal()
		}

		openSearch, err := opensearch.NewClient(opensearch.Config{
			Addresses: []string{os.Getenv("OPEN_SEARCH_ENDPOINT")},
			Signer:    signer,
		})
		if err != nil {
			Logger.Error("Error creating openSearchClient", slog.Any("error", err))
			log.Fatal()
		}

		OpenSearch = openSearch
	}

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
	IdempotencyInterceptor = connectkit.NewIdempotencyInterceptor(DynamoRead, DynamoWrite, TableName)
}
