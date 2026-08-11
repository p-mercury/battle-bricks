package main

import (
	"connectkit"
	"context"
	identityevents "contracts/dist/identity/events"
	"contracts/dist/identity/v1"
	"errors"
	"fmt"
	"identity/src/types/dynamo"
	"log/slog"
	"slices"
	"strconv"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	cognitoTypes "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	eventTypes "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Handler) UpdateUser(
	ctx context.Context,
	req *connect.Request[identity.UpdateUserRequest],
) (*connect.Response[identity.UpdateUserResponse], error) {
	logger := connectkit.GetLogger(ctx)
	authCtx := connectkit.GetAuthContext(ctx)
	timestamp := time.Now().Truncate(time.Millisecond)

	var cognitoId string
	var user *identity.User
	{
		response, err := DynamoWrite.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: TableName,
			Key: map[string]dynamoTypes.AttributeValue{
				"pk": &dynamoTypes.AttributeValueMemberS{Value: "USER#" + req.Msg.Id},
				"sk": &dynamoTypes.AttributeValueMemberS{Value: "USER"},
			},
		})
		if err != nil {
			logger.Error("Error getting user", slog.Any("error", err))
			return nil, connectkit.NewUnexpected()
		}
		if len(response.Item) > 0 {
			var item dynamo.User
			if err = attributevalue.UnmarshalMap(response.Item, &item); err != nil {
				logger.Error("Error parsing user", slog.Any("error", err))
				return nil, connectkit.NewUnexpected()
			}

			cognitoId = item.CognitoId
			user = parseDynamoUser(&item)
		}
	}

	if authCtx.Lambda != nil {
		resp, err := authCtx.Evaluate(ctx, "identity_customer", "update_user",
			&identity.UpdateUserContext{
				Request: req.Msg,
				Subject: user,
			},
		)
		if err != nil {
			logger.Error("Error evaluating policy", slog.Any("error", err))
			return nil, connectkit.NewUnexpected()
		} else if !resp.Authz {
			return nil, connectkit.NewUnauthorized()
		}
	}

	if user == nil {
		return nil, connectkit.NewNotFound()
	}

	if req.Msg.Version != nil && user.Version != *req.Msg.Version {
		return nil, connectkit.NewConflict("Version mismatch")
	}

	{
		transactItems := []types.TransactWriteItem{}
		var updateExpression strings.Builder
		updateExpression.WriteString("SET #ver = :ver, #mdt = :mdt")
		expressionAttributeNames := map[string]string{
			"#ver": "version",
			"#mdt": "modifiedTime",
			"#sta": "status",
		}
		expressionAttributeValues := map[string]dynamoTypes.AttributeValue{
			":ver": &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatUint(user.Version+1, 10)},
			":mdt": &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(timestamp.UnixMilli(), 10)},
			":dst": &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(int64(identity.UserStatus_USER_STATUS_DELETED), 10)},
		}
		var conditionExpression strings.Builder
		conditionExpression.WriteString("attribute_exists(pk) AND #ver < :ver AND #sta <> :dst AND (")

		if slices.Contains(req.Msg.UpdateMask, "correlations") {
			if len(req.Msg.Correlations) < 1 {
				return nil, connectkit.NewInvalidArgument("correlations")
			}

			expressionAttributeNames["#cor"] = "correlations"
			for i, correlation := range req.Msg.Correlations {
				transactItems = append(transactItems, types.TransactWriteItem{
					Put: &types.Put{
						TableName: TableName,
						Item: map[string]dynamoTypes.AttributeValue{
							"pk":   &dynamoTypes.AttributeValueMemberS{Value: "CORRELATION#USER#" + correlation.Provider + "#" + correlation.Kind + "#" + correlation.Id},
							"sk":   &dynamoTypes.AttributeValueMemberS{Value: "CORRELATION"},
							"type": &dynamoTypes.AttributeValueMemberS{Value: "CORRELATION"},

							"provider": &dynamoTypes.AttributeValueMemberS{Value: correlation.Provider},
							"kind":     &dynamoTypes.AttributeValueMemberS{Value: correlation.Kind},
							"id":       &dynamoTypes.AttributeValueMemberS{Value: correlation.Id},
							"targetId": &dynamoTypes.AttributeValueMemberS{Value: req.Msg.Id},
						},
						ExpressionAttributeNames: map[string]string{
							"#tid": "targetId",
						},
						ExpressionAttributeValues: map[string]dynamoTypes.AttributeValue{
							":tid": &dynamoTypes.AttributeValueMemberS{Value: req.Msg.Id},
						},
						ConditionExpression: new("attribute_not_exists(pk) OR #tid = :tid"),
					},
				})
				fmt.Fprintf(&updateExpression, ", #cor.#co%d = :co%d", i, i)
				expressionAttributeNames[fmt.Sprintf("#co%d", i)] = correlation.Provider + "#" + correlation.Kind
				expressionAttributeValues[fmt.Sprintf(":co%d", i)] = &dynamoTypes.AttributeValueMemberM{
					Value: map[string]dynamoTypes.AttributeValue{
						"provider": &dynamoTypes.AttributeValueMemberS{Value: correlation.Provider},
						"kind":     &dynamoTypes.AttributeValueMemberS{Value: correlation.Kind},
						"id":       &dynamoTypes.AttributeValueMemberS{Value: correlation.Id},
					},
				}
				fmt.Fprintf(&conditionExpression, " #cor.#co%d <> :co%d OR", i, i)
			}
		}

		if slices.Contains(req.Msg.UpdateMask, "email_address") {
			if req.Msg.EmailAddress == nil {
				return nil, connectkit.NewInvalidArgument("email_address")
			}

			emailAddress := strings.ToLower(*req.Msg.EmailAddress)

			updateExpression.WriteString(", #ead = :ead")
			expressionAttributeNames["#ead"] = "emailAddress"
			expressionAttributeValues[":ead"] = &dynamoTypes.AttributeValueMemberS{Value: emailAddress}
			conditionExpression.WriteString(" #ead <> :ead OR")
		}

		if slices.Contains(req.Msg.UpdateMask, "name") {
			if req.Msg.Name == nil {
				return nil, connectkit.NewInvalidArgument("name")
			}

			updateExpression.WriteString(", #nam = :nam")
			expressionAttributeNames["#nam"] = "name"
			expressionAttributeValues[":nam"] = &dynamoTypes.AttributeValueMemberS{Value: *req.Msg.Name}
			conditionExpression.WriteString(" #nam <> :nam OR")
		}

		if slices.Contains(req.Msg.UpdateMask, "status") {
			if req.Msg.Status == nil {
				return nil, connectkit.NewInvalidArgument("status")
			}

			updateExpression.WriteString(", #sta = :sta")
			expressionAttributeNames["#sta"] = "status"
			expressionAttributeValues[":sta"] = &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(int64(*req.Msg.Status), 10)}
			conditionExpression.WriteString(" #sta <> :sta OR")
		}

		if slices.Contains(req.Msg.UpdateMask, "language") {
			if req.Msg.Language == nil {
				return nil, connectkit.NewInvalidArgument("language")
			}

			updateExpression.WriteString(", #lng = :lng")
			expressionAttributeNames["#lng"] = "language"
			expressionAttributeValues[":lng"] = &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(int64(*req.Msg.Language), 10)}
			conditionExpression.WriteString(" #lng <> :lng OR")
		}

		if slices.Contains(req.Msg.UpdateMask, "notification_settings") {
			if req.Msg.NotificationSettings == nil {
				return nil, connectkit.NewInvalidArgument("notification_settings")
			}

			updateExpression.WriteString(", #nos = :nos")
			expressionAttributeNames["#nos"] = "notificationSettings"
			expressionAttributeValues[":nos"] = &dynamoTypes.AttributeValueMemberM{Value: map[string]dynamoTypes.AttributeValue{
				"newsletter": &dynamoTypes.AttributeValueMemberM{Value: map[string]dynamoTypes.AttributeValue{
					"enabled": &dynamoTypes.AttributeValueMemberBOOL{Value: req.Msg.NotificationSettings.Newsletter.Enabled},
					"methods": &dynamoTypes.AttributeValueMemberL{Value: []dynamoTypes.AttributeValue{
						&dynamoTypes.AttributeValueMemberN{Value: "1"},
						&dynamoTypes.AttributeValueMemberN{Value: "2"},
						&dynamoTypes.AttributeValueMemberN{Value: "3"},
					}},
				}},
				"purchasingQuotationPriced": &dynamoTypes.AttributeValueMemberM{Value: map[string]dynamoTypes.AttributeValue{
					"enabled": &dynamoTypes.AttributeValueMemberBOOL{Value: req.Msg.NotificationSettings.PurchasingQuotationPriced.Enabled},
					"methods": &dynamoTypes.AttributeValueMemberL{Value: []dynamoTypes.AttributeValue{
						&dynamoTypes.AttributeValueMemberN{Value: "1"},
						&dynamoTypes.AttributeValueMemberN{Value: "2"},
						&dynamoTypes.AttributeValueMemberN{Value: "3"},
					}},
				}},
				"purchasingOrderConfirmed": &dynamoTypes.AttributeValueMemberM{Value: map[string]dynamoTypes.AttributeValue{
					"enabled": &dynamoTypes.AttributeValueMemberBOOL{Value: req.Msg.NotificationSettings.PurchasingOrderConfirmed.Enabled},
					"methods": &dynamoTypes.AttributeValueMemberL{Value: []dynamoTypes.AttributeValue{
						&dynamoTypes.AttributeValueMemberN{Value: "1"},
						&dynamoTypes.AttributeValueMemberN{Value: "2"},
						&dynamoTypes.AttributeValueMemberN{Value: "3"},
					}},
				}},
			}}
			conditionExpression.WriteString(" #nos <> :nos OR")
		}

		transactItems = append(transactItems, types.TransactWriteItem{
			Update: &types.Update{
				TableName: TableName,
				Key: map[string]dynamoTypes.AttributeValue{
					"pk": &dynamoTypes.AttributeValueMemberS{Value: "USER#" + req.Msg.Id},
					"sk": &dynamoTypes.AttributeValueMemberS{Value: "USER"},
				},
				UpdateExpression:                    new(updateExpression.String()),
				ExpressionAttributeNames:            expressionAttributeNames,
				ExpressionAttributeValues:           expressionAttributeValues,
				ConditionExpression:                 new(strings.TrimSuffix(conditionExpression.String(), "OR") + ")"),
				ReturnValuesOnConditionCheckFailure: dynamoTypes.ReturnValuesOnConditionCheckFailureAllOld,
			},
		})

		if _, err := DynamoWrite.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{TransactItems: transactItems}); err != nil {
			var cfe *dynamoTypes.TransactionCanceledException
			if errors.As(err, &cfe) {
				if len(cfe.CancellationReasons) > 0 && cfe.CancellationReasons[len(cfe.CancellationReasons)-1].Item != nil {
					var u dynamo.User
					if err := attributevalue.UnmarshalMap(cfe.CancellationReasons[len(cfe.CancellationReasons)-1].Item, &u); err != nil {
						logger.Error("Error parsing old user", slog.Any("error", err))
						return nil, connectkit.NewUnexpected()
					}
					if u.Version == user.Version && u.Status != identity.UserStatus_USER_STATUS_DELETED {
						return connect.NewResponse(&identity.UpdateUserResponse{User: user}), nil
					} else {
						logger.Error("Error updating user", slog.Any("error", err))
						return nil, connectkit.NewUnexpected()
					}
				} else {
					return nil, connectkit.NewNotFound()
				}
			} else {
				logger.Error("Error updating user", slog.Any("error", err))
				return nil, connectkit.NewUnexpected()
			}
		}
	}

	if slices.Contains(req.Msg.UpdateMask, "email_address") {
		if req.Msg.EmailAddress == nil {
			return nil, connectkit.NewInvalidArgument("email_address")
		}

		emailAddress := strings.ToLower(*req.Msg.EmailAddress)

		_, err := Cognito.AdminUpdateUserAttributes(ctx, &cognitoidentityprovider.AdminUpdateUserAttributesInput{
			UserPoolId: UserPoolId,
			Username:   &cognitoId,
			UserAttributes: []cognitoTypes.AttributeType{
				{Name: new("email"), Value: new(emailAddress)},
			},
		})
		if err != nil {
			logger.Error("Error updating cognito email", slog.Any("error", err))
			return nil, connectkit.NewUnexpected()
		}
	}

	user.Version++
	user.ModifiedTime = timestamppb.New(timestamp)

	changedFields := []identityevents.UserUpserted_ChangedField{}

	if slices.Contains(req.Msg.UpdateMask, "correlations") {
		for _, i := range req.Msg.Correlations {
			exists := false
			for _, n := range user.Correlations {
				if i.Provider == n.Provider && i.Kind == n.Kind {
					exists = true
					n.Id = i.Id
					break
				}
			}
			if !exists {
				user.Correlations = append(user.Correlations, i)
			}
		}
	}

	if slices.Contains(req.Msg.UpdateMask, "email_address") {
		user.EmailAddress = *req.Msg.EmailAddress
		if user.EmailAddress != *req.Msg.EmailAddress {
			changedFields = append(changedFields, identityevents.UserUpserted_CHANGED_FIELD_EMAIL_ADDRESS)
		}
	}

	if slices.Contains(req.Msg.UpdateMask, "name") {
		user.Name = *req.Msg.Name
		if user.Name != *req.Msg.Name {
			changedFields = append(changedFields, identityevents.UserUpserted_CHANGED_FIELD_NAME)
		}
	}

	if slices.Contains(req.Msg.UpdateMask, "status") {
		user.Status = *req.Msg.Status
		if user.Status != *req.Msg.Status {
			changedFields = append(changedFields, identityevents.UserUpserted_CHANGED_FIELD_STATUS)
		}
	}

	if slices.Contains(req.Msg.UpdateMask, "language") {
		user.Language = *req.Msg.Language
		if user.Language != *req.Msg.Language {
			changedFields = append(changedFields, identityevents.UserUpserted_CHANGED_FIELD_LANGUAGE)
		}
	}

	if slices.Contains(req.Msg.UpdateMask, "notification_settings") {
		user.NotificationSettings = req.Msg.NotificationSettings
		if user.NotificationSettings != req.Msg.NotificationSettings {
			changedFields = append(changedFields, identityevents.UserUpserted_CHANGED_FIELD_NOTIFICATION_SETTINGS)
		}
	}

	{
		event, err := identityevents.BuildUserUpdatedEvent(EventBusName, Namespace, &identityevents.UserUpserted{
			Id:           user.Id,
			Correlations: user.Correlations,
			Version:      user.Version,
			CreatedTime:  user.CreatedTime,
			ModifiedTime: user.ModifiedTime,

			ChangedFields: changedFields,

			EmailAddress:         user.EmailAddress,
			Status:               user.Status,
			Name:                 user.Name,
			Language:             user.Language,
			NotificationSettings: user.NotificationSettings,
		})
		if err != nil {
			logger.Error("Error building event", slog.Any("error", err))
			return nil, connectkit.NewUnexpected()
		}

		_, err = EventBus.PutEvents(ctx, &eventbridge.PutEventsInput{
			EndpointId: EventBusEndpointId,
			Entries:    []eventTypes.PutEventsRequestEntry{event},
		})
		if err != nil {
			logger.Error("Error emitting event", slog.Any("error", err))
			return nil, connectkit.NewUnexpected()
		}
	}

	return connect.NewResponse(&identity.UpdateUserResponse{User: user}), nil
}
