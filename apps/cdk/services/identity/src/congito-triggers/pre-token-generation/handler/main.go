package main

import (
	"context"
	"errors"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-xray-sdk-go/xray"
)

func handleRequest(
	ctx context.Context,
	event events.CognitoEventUserPoolsPreTokenGenV2_0,
) (events.CognitoEventUserPoolsPreTokenGenV2_0, error) {
	ctx, segment := xray.BeginSubsegment(ctx, "Handler")
	segment.AddAnnotation("Stack", StackName)
	defer segment.Close(nil)

	Logger.Info("Event", slog.Any("event", event))

	userID := event.Request.UserAttributes["custom:userId"]
	if userID == "" {
		err := errors.New("custom:userId is missing")

		Logger.Error(
			"Cannot generate token without userId",
			slog.String("username", event.UserName),
			slog.Any("error", err),
		)

		return event, err
	}

	claims := map[string]any{
		"userId": userID,
	}

	event.Response.ClaimsAndScopeOverrideDetails =
		events.ClaimsAndScopeOverrideDetailsV2_0{
			IDTokenGeneration: events.IDTokenGenerationV2_0{
				ClaimsToAddOrOverride: claims,
			},
			AccessTokenGeneration: events.AccessTokenGenerationV2_0{
				ClaimsToAddOrOverride: claims,
			},
		}

	return event, nil
}

func main() {
	lambda.Start(handleRequest)
}
