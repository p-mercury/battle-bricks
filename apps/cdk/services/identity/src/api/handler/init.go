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
	"github.com/auth0/go-auth0/management"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-xray-sdk-go/instrumentation/awsv2"
	"github.com/aws/aws-xray-sdk-go/xray"
)

var (
	Logger      *slog.Logger
	DynamoRead  *dynamodb.Client
	DynamoWrite *dynamodb.Client
	Auth0       *management.Management
	EventBus    *eventbridge.Client
	Http        *http.Client

	PolicyService policyconnect.InternalClient

	StackName          string
	Namespace          string
	ApiUrl             string
	Hostname           *string
	TableName          *string
	EventBusName       *string
	EventBusEndpointId *string

	Auth0Audience       *string
	Auth0ClientId       *string
	Auth0ClientDomain   *string
	Auth0ClientSecret   *string
	Auth0ConnnectinName *string

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

	Auth0Audience = new(os.Getenv("AUTH0_AUDIENCE"))
	Auth0ClientId = new(os.Getenv("AUTH0_CLIENT_ID"))
	Auth0ClientDomain = new(os.Getenv("AUTH0_CLIENT_DOMAIN"))
	Auth0ClientSecret = new(os.Getenv("AUTH0_CLIENT_SECRET"))
	Auth0ConnnectinName = new(os.Getenv("AUTH0_CONNECTION_NAME"))

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
	EventBus = eventbridge.NewFromConfig(cfg)
	Http = xray.Client(nil)

	{
		auth0, err := management.New(
			os.Getenv("AUTH0_MANAGEMENT_CLIENT_DOMAIN"),
			management.WithClientCredentials(context.Background(), os.Getenv("AUTH0_MANAGEMENT_CLIENT_ID"), os.Getenv("AUTH0_MANAGEMENT_CLIENT_SECRET")),
		)
		if err != nil {
			Logger.Error("Error initializing the auth0 management API client", slog.Any("error", err))
			log.Fatal()
		}
		Auth0 = auth0
	}

	{
		authHandler, err := NewAuthHandler(
			context.Background(),
			*Auth0ClientDomain,
			*Auth0ClientId,
			*Auth0ClientSecret)
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
