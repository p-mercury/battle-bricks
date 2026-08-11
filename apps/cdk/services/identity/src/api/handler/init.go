package main

import (
	"connectkit"
	"context"
	"contracts/dist/policy/v1/policyconnect"
	"log"
	"log/slog"
	"net/http"
	"os"
	"sigv4http"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-xray-sdk-go/instrumentation/awsv2"
	"github.com/aws/aws-xray-sdk-go/xray"
)

var (
	Logger      *slog.Logger
	DynamoRead  *dynamodb.Client
	DynamoWrite *dynamodb.Client
	EventBus    *eventbridge.Client
	Cognito     *cognitoidentityprovider.Client
	Http        *http.Client

	PolicyService policyconnect.InternalClient

	StackName          string
	Namespace          string
	ApiUrl             string
	Hostname           *string
	TableName          *string
	EventBusName       *string
	EventBusEndpointId *string

	UserPoolId           *string
	UserPoolUrl          *string
	UserPoolProviderUrl  *string
	UserPoolClientId     *string
	UserPoolClientSecret *string

	AuthHandler            *OAuthHandler
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
	Hostname = new(os.Getenv("HOSTNAME"))
	EventBusName = new(os.Getenv("EVENT_BUS_NAME"))
	EventBusEndpointId = new(os.Getenv("EVENT_BUS_ENDPOINT_ID"))

	UserPoolId = new(os.Getenv("USER_POOL_ID"))
	UserPoolUrl = new(os.Getenv("USER_POOL_URL"))
	UserPoolProviderUrl = new(os.Getenv("USER_POOL_PROVIDER_URL"))
	UserPoolClientId = new(os.Getenv("USER_POOL_CLIENT_ID"))
	UserPoolClientSecret = new(os.Getenv("USER_POOL_CLIENT_SECRET"))

	{
		cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(os.Getenv("TABLE_WRITE_REGION")))
		if err != nil {
			log.Fatalf("Unable to configure aws services, %v", err)
			return
		}

		awsv2.AWSV2Instrumentor(&cfg.APIOptions)

		DynamoWrite = dynamodb.NewFromConfig(cfg)
		Cognito = cognitoidentityprovider.NewFromConfig(cfg)
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		Logger.Error("Unable to configure aws services", slog.Any("error", err))
		log.Fatal()
	}

	awsv2.AWSV2Instrumentor(&cfg.APIOptions)

	DynamoRead = dynamodb.NewFromConfig(cfg)
	EventBus = eventbridge.NewFromConfig(cfg)
	Http = xray.Client(nil)

	{
		authHandler, err := NewAuthHandler(
			context.Background(),
			*UserPoolUrl,
			*UserPoolProviderUrl,
			*UserPoolClientId,
			*UserPoolClientSecret)
		if err != nil {
			Logger.Error("Error creating auth handler", slog.Any("error", err))
			log.Fatal()
		}
		AuthHandler = authHandler
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
