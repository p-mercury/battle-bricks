package main

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log/slog"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	eventTypes "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/aws/aws-xray-sdk-go/xray"
)

type identity struct {
	ProviderType string
	UserId       string
}

func generateLedgerKey(event events.CognitoEventUserPoolsPostConfirmation) (*string, error) {
	var raw string
	if i, e := event.Request.UserAttributes["identities"]; event.Request.UserAttributes["cognito:user_status"] == "EXTERNAL_PROVIDER" && e {
		var identities []identity
		err := json.Unmarshal([]byte(i), &identities)
		if err != nil {
			return nil, err
		}

		if len(identities) != 1 {
			return nil, errors.New("Identities should only contain one item")
		}

		raw = identities[0].ProviderType + ":" + identities[0].UserId
	} else {
		raw = "Email:" + event.Request.UserAttributes["email"]
	}

	sum := sha512.Sum512([]byte(raw))
	key := hex.EncodeToString(sum[:])
	return &key, nil
}

func handleRequest(ctx context.Context, event events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error) {
	ctx, segment := xray.BeginSubsegment(ctx, "Handler")
	segment.AddAnnotation("Stack", StackName)
	defer segment.Close(nil)

	timestamp := time.Now().UnixMilli()

	Logger.Info("Event", slog.Any("event", event))

	if event.TriggerSource != "PostConfirmation_ConfirmSignUp" {
		return event, nil
	}

	item := map[string]dynamoTypes.AttributeValue{
		"pk":           &dynamoTypes.AttributeValueMemberS{Value: "USER#" + event.Request.UserAttributes["sub"]},
		"sk":           &dynamoTypes.AttributeValueMemberS{Value: "USER"},
		"emailAddress": &dynamoTypes.AttributeValueMemberS{Value: event.Request.UserAttributes["email"]},
		"name":         &dynamoTypes.AttributeValueMemberS{Value: event.Request.UserAttributes["fullname"]},
		"language":     &dynamoTypes.AttributeValueMemberN{Value: "0"},
		"provider":     &dynamoTypes.AttributeValueMemberS{Value: "Cognito"},
		"createdTime":  &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(timestamp, 10)},
		"modifiedTime": &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(timestamp, 10)},
	}

	if i, e := event.Request.UserAttributes["identities"]; event.Request.UserAttributes["cognito:user_status"] == "EXTERNAL_PROVIDER" && e {
		var identities []identity
		err := json.Unmarshal([]byte(i), &identities)
		if err == nil {
			if len(identities) == 1 {
				item["provider"] = &dynamoTypes.AttributeValueMemberS{Value: identities[0].ProviderType}
			}
		}
	}

	_, err := DynamoWrite.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           TableName,
		Item:                item,
		ConditionExpression: new("attribute_not_exists(pk)"),
	})
	if err != nil {
		Logger.Error("Error creating dynamo user", slog.Any("error", err))
		return event, err
	}

	{
		eventDetail, err := json.Marshal(map[string]string{
			"id": event.Request.UserAttributes["sub"],
		})
		if err == nil {
			_, err := EventBus.PutEvents(ctx, &eventbridge.PutEventsInput{
				Entries: []eventTypes.PutEventsRequestEntry{
					{
						EventBusName: EventBusName,
						Source:       new(Namespace + ".identity"),
						DetailType:   new("USER_CREATED"),
						Detail:       new(string(eventDetail)),
					},
				},
			})
			if err != nil {
				Logger.Error("Error emitting 'USER_CREATED' event", slog.Any("error", err))
			}
		} else {
			Logger.Error("Error generating 'USER_CREATED' event", slog.Any("error", err))
		}
	}

	return event, nil
}

func main() {
	lambda.Start(handleRequest)
}
