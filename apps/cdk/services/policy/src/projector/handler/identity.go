package main

import (
	"context"
	identityevents "contracts/dist/identity/events"
	"errors"
	"log/slog"
	"policy/src/types/dynamo"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

		var organisation dynamo.Organisation
		{
			response, err := DynamoRead.GetItem(ctx, &dynamodb.GetItemInput{
				TableName: TableName,
				Key: map[string]dynamoTypes.AttributeValue{
					"pk": &dynamoTypes.AttributeValueMemberS{Value: details.OrganisationId},
				}})
			if err != nil {
				Logger.Error("Error getting organisation", slog.Any("error", err))
				return err
			}
			if response.Item == nil {
				Logger.Warn("Organisation not found", slog.String("id", details.OrganisationId))
				return err
			}
			if err = attributevalue.UnmarshalMap(response.Item, &organisation); err != nil {
				Logger.Warn("Error parsing organisation", slog.Any("item", response.Item), slog.Any("error", err))
				return err
			}
		}

		_, err := DynamoWrite.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
			TransactItems: []types.TransactWriteItem{
				{
					ConditionCheck: &types.ConditionCheck{
						TableName: TableName,
						Key: map[string]dynamoTypes.AttributeValue{
							"pk": &dynamoTypes.AttributeValueMemberS{Value: organisation.Id},
						},
						ExpressionAttributeNames: map[string]string{
							"#ver": "version",
						},
						ExpressionAttributeValues: map[string]types.AttributeValue{
							":ver": &types.AttributeValueMemberN{Value: strconv.FormatInt(int64(organisation.Version), 10)},
						},
						ConditionExpression: new("#ver = :ver"),
					},
				},
				{
					Update: &types.Update{
						TableName: TableName,
						Key: map[string]dynamoTypes.AttributeValue{
							"pk": &dynamoTypes.AttributeValueMemberS{Value: details.Id},
						},
						UpdateExpression: new("SET #gsi1pk = :gsi1pk, #gsi1sk = :gsi1sk, #typ = :typ, #sta = :sta, #ver = :ver, #org = :org"),
						ExpressionAttributeNames: map[string]string{
							"#gsi1pk": "gsi1pk",
							"#gsi1sk": "gsi1sk",
							"#typ":    "type",
							"#sta":    "status",
							"#ver":    "version",
							"#org":    "organisation",
							"#tol":    "ttl",
						},
						ExpressionAttributeValues: map[string]dynamoTypes.AttributeValue{
							":gsi1pk": &dynamoTypes.AttributeValueMemberS{Value: details.OrganisationId},
							":gsi1sk": &dynamoTypes.AttributeValueMemberS{Value: details.Id},
							":typ":    &dynamoTypes.AttributeValueMemberS{Value: "USER"},
							":sta":    &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(int64(details.Status), 10)},
							":ver":    &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatUint(details.Version, 10)},
							":org": &dynamoTypes.AttributeValueMemberM{Value: map[string]dynamoTypes.AttributeValue{
								"entitlements": &dynamoTypes.AttributeValueMemberM{Value: map[string]dynamoTypes.AttributeValue{
									"maxUsers": &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(int64(organisation.Entitlements.MaxUsers), 10)},
									"platforms": &dynamoTypes.AttributeValueMemberM{Value: map[string]dynamoTypes.AttributeValue{
										"sales": &dynamoTypes.AttributeValueMemberM{Value: map[string]dynamoTypes.AttributeValue{
											"status":              &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(int64(organisation.Entitlements.Platforms.Sales.Status), 10)},
											"integration":         &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(int64(organisation.Entitlements.Platforms.Sales.Integration), 10)},
											"allowPublicListings": &dynamoTypes.AttributeValueMemberBOOL{Value: organisation.Entitlements.Platforms.Sales.AllowPublicListings},
											"maxProducts":         &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(int64(organisation.Entitlements.Platforms.Sales.MaxProducts), 10)},
										}},
										"purchasing": &dynamoTypes.AttributeValueMemberM{Value: map[string]dynamoTypes.AttributeValue{
											"status":      &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(int64(organisation.Entitlements.Platforms.Purchasing.Status), 10)},
											"integration": &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(int64(organisation.Entitlements.Platforms.Purchasing.Integration), 10)},
										}},
									}},
								}},
								"version": &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(int64(organisation.Version), 10)},
							}},
						},
						ConditionExpression: new("attribute_not_exists(#tol) AND (attribute_not_exists(#ver) OR #ver < :ver)"),
					},
				},
			},
		})
		if err != nil && !errors.As(err, new(*types.TransactionCanceledException)) {
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
		if err != nil && !errors.As(err, new(*types.ConditionalCheckFailedException)) {
			Logger.Error("Error deleting user")
			return err
		}
	}

	return nil
}
