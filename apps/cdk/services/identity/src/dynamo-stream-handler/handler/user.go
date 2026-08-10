package main

import (
	"context"
	"identity/src/types/dynamo"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func handlerUser(ctx context.Context, record events.DynamoDBEventRecord) error {

	var oldImage *dynamo.User
	{
		if record.Change.OldImage != nil {
			image, err := FromDynamoDBEventAttributeValueMap(record.Change.OldImage)
			if err != nil {
				Logger.Error("Unable to map old image", slog.Any("error", err))
				return err
			}
			if err := attributevalue.UnmarshalMap(image, &oldImage); err != nil {
				Logger.Error("Unable to parse old image", slog.Any("error", err))
				return err
			}
		}
	}

	var newImage *dynamo.User
	{
		if record.Change.NewImage != nil {
			image, err := FromDynamoDBEventAttributeValueMap(record.Change.NewImage)
			if err != nil {
				Logger.Error("Unable to map new image", slog.Any("error", err))
				return err
			}
			if err := attributevalue.UnmarshalMap(image, &newImage); err != nil {
				Logger.Error("Unable to parse new image", slog.Any("error", err))
				return err
			}
		}
	}

	switch record.EventName {
	case "INSERT":
	case "MODIFY":
		for _, oc := range oldImage.Correlations {
			delete := true
			if newImage.Correlations != nil {
				for _, nc := range newImage.Correlations {
					if oc.Provider == nc.Provider &&
						oc.Kind == nc.Kind &&
						oc.Id == nc.Id {
						delete = false
						break
					}
				}
			}
			if delete {
				DynamoWrite.DeleteItem(ctx, &dynamodb.DeleteItemInput{
					TableName: TableName,
					Key: map[string]dynamoTypes.AttributeValue{
						"pk": &dynamoTypes.AttributeValueMemberS{Value: "CORRELATION#USER#" + oc.Provider + "#" + oc.Kind + "#" + oc.Id},
						"sk": &dynamoTypes.AttributeValueMemberS{Value: "CORRELATION"},
					},
				})
			}
		}

	case "REMOVE":
		for _, oc := range oldImage.Correlations {
			DynamoWrite.DeleteItem(ctx, &dynamodb.DeleteItemInput{
				TableName: TableName,
				Key: map[string]dynamoTypes.AttributeValue{
					"pk": &dynamoTypes.AttributeValueMemberS{Value: "CORRELATION#USER#" + oc.Provider + "#" + oc.Kind + "#" + oc.Id},
					"sk": &dynamoTypes.AttributeValueMemberS{Value: "CORRELATION"},
				},
			})
		}
	}

	return nil
}
