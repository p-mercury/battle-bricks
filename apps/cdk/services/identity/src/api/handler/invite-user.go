package main

import (
	"connectkit"
	"context"
	"contracts/dist/common/v1"
	identityevents "contracts/dist/identity/events"
	"contracts/dist/identity/v1"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/auth0/go-auth0/management"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	eventTypes "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Handler) InviteUser(
	ctx context.Context,
	req *connect.Request[identity.InviteUserRequest],
) (*connect.Response[identity.InviteUserResponse], error) {
	logger := connectkit.GetLogger(ctx)
	authCtx := connectkit.GetAuthContext(ctx)
	timestamp := time.Now().Truncate(time.Millisecond)

	if authCtx.Lambda != nil {
		resp, err := authCtx.Evaluate(ctx, "identity_customer", "invite_user",
			&identity.InviteUserContext{
				Request: req.Msg,
			},
		)
		if err != nil {
			logger.Error("Error evaluating policy", slog.Any("error", err))
			return nil, connectkit.NewUnexpected()
		} else if !resp.Authz {
			return nil, connectkit.NewUnauthorized()
		}
	}

	emailAddress := strings.ToLower(req.Msg.EmailAddress)

	{
		dynamoResponse, err := DynamoWrite.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: TableName,
			Key: map[string]dynamoTypes.AttributeValue{
				"pk": &dynamoTypes.AttributeValueMemberS{Value: "ORGANISATIONS"},
				"sk": &dynamoTypes.AttributeValueMemberS{Value: "ORGANISATION#" + req.Msg.OrganisationId},
			},
		})
		if err != nil {
			logger.Error("Error getting organistaion from dynamoDB", slog.Any("error", err))
			return nil, connectkit.NewUnexpected()
		}
		if len(dynamoResponse.Item) < 1 {
			return nil, connectkit.NewInvalidArgument("organisation_id")
		}
	}

	pass, err := uuid.NewRandom()
	if err != nil {
		logger.Error("Error generating random password", slog.Any("error", err))
		return nil, connectkit.NewUnexpected()
	}

	var id string
	{
		for a := range 4 {
			nid, err := connectkit.NewBase62Id("c-", 12)
			if err != nil {
				logger.Error("Error generating base62 id", slog.Any("error", err))
				return nil, connectkit.NewUnexpected()
			}

			response, err := DynamoRead.GetItem(ctx, &dynamodb.GetItemInput{
				TableName: TableName,
				Key: map[string]dynamoTypes.AttributeValue{
					"pk": &dynamoTypes.AttributeValueMemberS{Value: "USER#" + nid},
					"sk": &dynamoTypes.AttributeValueMemberS{Value: "USER"},
				},
				ProjectionExpression: new("pk"),
			})
			if err != nil {
				logger.Error("Error getting user", slog.Any("error", err))
				return nil, connectkit.NewUnexpected()
			}

			if response.Item == nil {
				id = nid
				break
			}

			if a > 2 {
				logger.Error("Can't find free user id")
				return nil, connectkit.NewUnexpected()
			}
		}
	}

	auth0User := &management.User{
		Connection:    Auth0ConnnectinName,
		ID:            new(id),
		Email:         new(emailAddress),
		Password:      new(pass.String()),
		EmailVerified: new(false),
		VerifyEmail:   new(false),
		AppMetadata: &map[string]any{
			"userId":         id,
			"organisationId": req.Msg.OrganisationId,
		},
	}

	err = Auth0.User.Create(ctx, auth0User)
	if err != nil {
		logger.Error("Error creating auth0 user", slog.Any("error", err))
		return nil, connectkit.NewUnexpected()
	}

	{
		transactItems := []types.TransactWriteItem{}
		item := map[string]dynamoTypes.AttributeValue{
			"pk":     &dynamoTypes.AttributeValueMemberS{Value: "USER#" + id},
			"sk":     &dynamoTypes.AttributeValueMemberS{Value: "USER"},
			"gsi1pk": &dynamoTypes.AttributeValueMemberS{Value: "ORGANISATION#" + req.Msg.OrganisationId},
			"gsi1sk": &dynamoTypes.AttributeValueMemberS{Value: "USER#" + id},
			"type":   &dynamoTypes.AttributeValueMemberS{Value: "USER"},

			"id":           &dynamoTypes.AttributeValueMemberS{Value: id},
			"version":      &dynamoTypes.AttributeValueMemberN{Value: "1"},
			"createdTime":  &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(timestamp.UnixMilli(), 10)},
			"modifiedTime": &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(timestamp.UnixMilli(), 10)},

			"organisationId": &dynamoTypes.AttributeValueMemberS{Value: req.Msg.OrganisationId},
			"status":         &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(int64(identity.UserStatus_USER_STATUS_INVITED), 10)},
			"emailAddress":   &dynamoTypes.AttributeValueMemberS{Value: emailAddress},
			"auth0Id":        &dynamoTypes.AttributeValueMemberS{Value: *auth0User.ID},
			"language":       &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(int64(common.Language_LANGUAGE_UNSPECIFIED), 10)},
			"notificationSettings": &dynamoTypes.AttributeValueMemberM{Value: map[string]dynamoTypes.AttributeValue{
				"newsletter": &dynamoTypes.AttributeValueMemberM{Value: map[string]dynamoTypes.AttributeValue{
					"enabled": &dynamoTypes.AttributeValueMemberBOOL{Value: true},
					"methods": &dynamoTypes.AttributeValueMemberNS{Value: []string{"1", "2", "3"}},
				}},
				"purchasingQuotationPriced": &dynamoTypes.AttributeValueMemberM{Value: map[string]dynamoTypes.AttributeValue{
					"enabled": &dynamoTypes.AttributeValueMemberBOOL{Value: true},
					"methods": &dynamoTypes.AttributeValueMemberNS{Value: []string{"1", "2", "3"}},
				}},
				"purchasingOrderConfirmed": &dynamoTypes.AttributeValueMemberM{Value: map[string]dynamoTypes.AttributeValue{
					"enabled": &dynamoTypes.AttributeValueMemberBOOL{Value: true},
					"methods": &dynamoTypes.AttributeValueMemberNS{Value: []string{"1", "2", "3"}},
				}},
			}},
		}

		{
			correlations := map[string]dynamoTypes.AttributeValue{}
			for _, correlation := range req.Msg.Correlations {
				if _, ok := correlations[correlation.Provider+"#"+correlation.Kind]; !ok {
					correlations[correlation.Provider+"#"+correlation.Kind] = &dynamoTypes.AttributeValueMemberM{
						Value: map[string]dynamoTypes.AttributeValue{
							"provider": &dynamoTypes.AttributeValueMemberS{Value: correlation.Provider},
							"kind":     &dynamoTypes.AttributeValueMemberS{Value: correlation.Kind},
							"id":       &dynamoTypes.AttributeValueMemberS{Value: correlation.Id},
						},
					}
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
								"targetId": &dynamoTypes.AttributeValueMemberS{Value: id},
							},
							ConditionExpression: new("attribute_not_exists(pk)"),
						},
					})
				}
			}
			item["correlations"] = &dynamoTypes.AttributeValueMemberM{Value: correlations}
		}

		transactItems = append(transactItems, types.TransactWriteItem{
			Put: &types.Put{
				TableName:           TableName,
				Item:                item,
				ConditionExpression: new("attribute_not_exists(sk)"),
			},
		})

		if _, err := DynamoWrite.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{TransactItems: transactItems}); err != nil {
			logger.Error("Error creating user", slog.Any("error", err))
			deleteUser(ctx, id, *auth0User.ID)
			return nil, connectkit.NewUnexpected()
		}
	}

	notificationSettings := &identity.UserNotificationSettings{
		Newsletter: &identity.UserNotificationSettings_Topic{
			Enabled: true,
			Methods: []identity.UserNotificationSettings_Topic_Method{
				identity.UserNotificationSettings_Topic_METHOD_EMAIL,
				identity.UserNotificationSettings_Topic_METHOD_IN_APP,
				identity.UserNotificationSettings_Topic_METHOD_PUSH,
			},
		},
		PurchasingQuotationPriced: &identity.UserNotificationSettings_Topic{
			Enabled: true,
			Methods: []identity.UserNotificationSettings_Topic_Method{
				identity.UserNotificationSettings_Topic_METHOD_EMAIL,
				identity.UserNotificationSettings_Topic_METHOD_IN_APP,
				identity.UserNotificationSettings_Topic_METHOD_PUSH,
			},
		},
		PurchasingOrderConfirmed: &identity.UserNotificationSettings_Topic{
			Enabled: true,
			Methods: []identity.UserNotificationSettings_Topic_Method{
				identity.UserNotificationSettings_Topic_METHOD_EMAIL,
				identity.UserNotificationSettings_Topic_METHOD_IN_APP,
				identity.UserNotificationSettings_Topic_METHOD_PUSH,
			},
		},
	}

	{
		event, err := identityevents.BuildUserCreatedEvent(EventBusName, Namespace, &identityevents.UserUpserted{
			Id:           id,
			Correlations: req.Msg.Correlations,
			Version:      1,
			CreatedTime:  timestamppb.New(timestamp),
			ModifiedTime: timestamppb.New(timestamp),

			ChangedFields: []identityevents.UserUpserted_ChangedField{
				identityevents.UserUpserted_CHANGED_FIELD_ORGANISATION_ID,
				identityevents.UserUpserted_CHANGED_FIELD_EMAIL_ADDRESS,
				identityevents.UserUpserted_CHANGED_FIELD_STATUS,
				identityevents.UserUpserted_CHANGED_FIELD_NAME,
				identityevents.UserUpserted_CHANGED_FIELD_JOB_TITLE,
				identityevents.UserUpserted_CHANGED_FIELD_LANGUAGE,
				identityevents.UserUpserted_CHANGED_FIELD_NOTIFICATION_SETTINGS,
			},

			OrganisationId:       req.Msg.OrganisationId,
			EmailAddress:         emailAddress,
			Status:               identity.UserStatus_USER_STATUS_INVITED,
			Name:                 emailAddress,
			Language:             common.Language_LANGUAGE_UNSPECIFIED,
			NotificationSettings: notificationSettings,
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

	return connect.NewResponse(&identity.InviteUserResponse{
		User: &identity.User{
			Id:           id,
			Correlations: req.Msg.Correlations,
			Version:      1,
			CreatedTime:  timestamppb.New(timestamp),
			ModifiedTime: timestamppb.New(timestamp),

			OrganisationId:       req.Msg.OrganisationId,
			EmailAddress:         emailAddress,
			Status:               identity.UserStatus_USER_STATUS_INVITED,
			Name:                 emailAddress,
			Language:             common.Language_LANGUAGE_UNSPECIFIED,
			NotificationSettings: notificationSettings,
		},
	}), nil
}

func deleteUser(ctx context.Context, id string, auth0Id string) {
	logger := connectkit.GetLogger(ctx)

	err := Auth0.User.Delete(ctx, auth0Id)
	if err != nil {
		logger.Error("Error deleting auth0 user", slog.Any("error", err))
	}

	_, err = DynamoWrite.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: TableName,
		Key: map[string]dynamoTypes.AttributeValue{
			"pk": &dynamoTypes.AttributeValueMemberS{Value: "USERS"},
			"sk": &dynamoTypes.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		logger.Error("Error deleting dynamo user", slog.Any("error", err))
	}
}
