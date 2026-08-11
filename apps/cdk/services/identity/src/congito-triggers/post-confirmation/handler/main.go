package main

import (
	"connectkit"
	"context"
	"contracts/dist/common/v1"
	identityevents "contracts/dist/identity/events"
	"contracts/dist/identity/v1"
	"encoding/json"
	"log/slog"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	cognitoTypes "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	eventTypes "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/aws/aws-xray-sdk-go/xray"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type cognitoIdentity struct {
	ProviderType string
	UserId       string
}

func handleRequest(ctx context.Context, event events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error) {
	ctx, segment := xray.BeginSubsegment(ctx, "Handler")
	segment.AddAnnotation("Stack", StackName)
	defer segment.Close(nil)

	timestamp := time.Now().Truncate(time.Millisecond)

	Logger.Info("Event", slog.Any("event", event))

	id, err := connectkit.NewBase62Id("", 12)
	if err != nil {
		Logger.Error("Error generating base62 id", slog.Any("error", err))
		return event, err
	}

	if event.TriggerSource != "PostConfirmation_ConfirmSignUp" {
		return event, nil
	}

	item := map[string]dynamoTypes.AttributeValue{
		"pk":           &dynamoTypes.AttributeValueMemberS{Value: "USER#" + id},
		"sk":           &dynamoTypes.AttributeValueMemberS{Value: "USER"},
		"cognitoId":    &dynamoTypes.AttributeValueMemberS{Value: event.Request.UserAttributes["sub"]},
		"emailAddress": &dynamoTypes.AttributeValueMemberS{Value: event.Request.UserAttributes["email"]},
		"name":         &dynamoTypes.AttributeValueMemberS{Value: event.Request.UserAttributes["fullname"]},
		"language":     &dynamoTypes.AttributeValueMemberN{Value: "0"},
		"provider":     &dynamoTypes.AttributeValueMemberS{Value: "Cognito"},
		"createdTime":  &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(timestamp.UnixMilli(), 10)},
		"modifiedTime": &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(timestamp.UnixMilli(), 10)},
	}

	if i, e := event.Request.UserAttributes["identities"]; event.Request.UserAttributes["cognito:user_status"] == "EXTERNAL_PROVIDER" && e {
		var identities []cognitoIdentity
		err := json.Unmarshal([]byte(i), &identities)
		if err == nil {
			if len(identities) == 1 {
				item["provider"] = &dynamoTypes.AttributeValueMemberS{Value: identities[0].ProviderType}
			}
		}
	}

	_, err = DynamoWrite.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           TableName,
		Item:                item,
		ConditionExpression: new("attribute_not_exists(pk)"),
	})
	if err != nil {
		Logger.Error("Error creating dynamo user", slog.Any("error", err))
		return event, err
	}

	{
		_, err := Cognito.AdminUpdateUserAttributes(ctx, &cognitoidentityprovider.AdminUpdateUserAttributesInput{
			UserPoolId: &event.UserPoolID,
			Username:   new(event.Request.UserAttributes["sub"]),
			UserAttributes: []cognitoTypes.AttributeType{
				{Name: new("userId"), Value: new(id)},
			},
		})
		if err != nil {
			Logger.Error("Error updating cognito id", slog.Any("error", err))
		}
	}

	{
		event, err := identityevents.BuildUserCreatedEvent(EventBusName, Namespace, &identityevents.UserUpserted{
			Id:           id,
			Correlations: []*common.Correlation{},
			Version:      1,
			CreatedTime:  timestamppb.New(timestamp),
			ModifiedTime: timestamppb.New(timestamp),

			ChangedFields: []identityevents.UserUpserted_ChangedField{
				identityevents.UserUpserted_CHANGED_FIELD_STATUS,
				identityevents.UserUpserted_CHANGED_FIELD_EMAIL_ADDRESS,
				identityevents.UserUpserted_CHANGED_FIELD_NAME,
				identityevents.UserUpserted_CHANGED_FIELD_LANGUAGE,
				identityevents.UserUpserted_CHANGED_FIELD_NOTIFICATION_SETTINGS,
			},

			Status:       identity.UserStatus_USER_STATUS_ACTIVE,
			EmailAddress: event.Request.UserAttributes["email"],
			Name:         event.Request.UserAttributes["fullname"],
			Language:     common.Language_LANGUAGE_UNSPECIFIED,
		})
		if err != nil {
			Logger.Error("Error building event", slog.Any("error", err))
		} else {
			_, err = EventBus.PutEvents(ctx, &eventbridge.PutEventsInput{
				EndpointId: EventBusEndpointId,
				Entries:    []eventTypes.PutEventsRequestEntry{event},
			})
			if err != nil {
				Logger.Error("Error emitting event", slog.Any("error", err))
			}
		}
	}

	return event, nil
}

func main() {
	lambda.Start(handleRequest)
}
