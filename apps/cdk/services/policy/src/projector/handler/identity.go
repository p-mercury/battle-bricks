package main

import (
	"context"
	identityevents "contracts/dist/identity/events"
	"errors"
	"log/slog"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"google.golang.org/protobuf/encoding/protojson"
)

func handlerIdentityEvent(ctx context.Context, event events.EventBridgeEvent) error {
	switch event.DetailType {
	case identityevents.UserCreatedDetailType,
		identityevents.UserUpdatedDetailType:
		var details identityevents.UserUpserted
		if err := protojson.Unmarshal(event.Detail, &details); err != nil {
			Logger.Error("Error parsing UserUpserted event", slog.Any("error", err))
			return err
		}

		_, err := DynamoWrite.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: TableName,
			Item: map[string]dynamoTypes.AttributeValue{
				"pk":      &dynamoTypes.AttributeValueMemberS{Value: details.Id},
				"type":    &dynamoTypes.AttributeValueMemberS{Value: "USER"},
				"status":  &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(int64(details.Status), 10)},
				"version": &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatUint(details.Version, 10)},
			},
			ExpressionAttributeNames: map[string]string{
				"#ver": "version",
				"#tol": "ttl",
			},
			ExpressionAttributeValues: map[string]dynamoTypes.AttributeValue{
				":ver": &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatUint(details.Version, 10)},
			},
			ConditionExpression: new("attribute_not_exists(#tol) AND (attribute_not_exists(#ver) OR #ver < :ver)"),
		})
		if err != nil && !errors.As(err, new(*dynamoTypes.ConditionalCheckFailedException)) {
			Logger.Error("Error updating user")
			return err
		}

	case identityevents.UserDeletedDetailType:
		var details identityevents.UserDeleted
		if err := protojson.Unmarshal(event.Detail, &details); err != nil {
			Logger.Error("Error parsing UserDeleted event", slog.Any("error", err))
			return err
		}

		_, err := DynamoWrite.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: TableName,
			Item: map[string]dynamoTypes.AttributeValue{
				"pk":  &dynamoTypes.AttributeValueMemberS{Value: details.Id},
				"ttl": &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(event.Time.Add(time.Hour*200).Unix(), 10)},
			},
			ExpressionAttributeNames: map[string]string{
				"#tol": "ttl",
			},
			ConditionExpression: new("attribute_not_exists(#tol)"),
		})
		if err != nil && !errors.As(err, new(*dynamoTypes.ConditionalCheckFailedException)) {
			Logger.Error("Error deleting user")
			return err
		}
	}

	return nil
}
