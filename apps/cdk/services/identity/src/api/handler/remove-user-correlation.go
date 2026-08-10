package main

import (
	"connectkit"
	"context"
	"contracts/dist/common/v1"
	identityevents "contracts/dist/identity/events"
	"contracts/dist/identity/v1"
	"errors"
	"identity/src/types/dynamo"
	"log/slog"
	"strconv"
	"time"

	"connectrpc.com/connect"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	eventTypes "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Handler) RemoveUserCorrelation(
	ctx context.Context,
	req *connect.Request[identity.RemoveUserCorrelationRequest],
) (*connect.Response[identity.RemoveUserCorrelationResponse], error) {
	logger := connectkit.GetLogger(ctx)
	authCtx := connectkit.GetAuthContext(ctx)
	timestamp := time.Now().Truncate(time.Millisecond)

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

			user = parseDynamoUser(&item)
		}
	}

	if authCtx.Lambda != nil {
		resp, err := authCtx.Evaluate(ctx, "identity_customer", "remove_user_correlation",
			&identity.RemoveUserCorrelationContext{
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

	_, err := DynamoWrite.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: TableName,
		Key: map[string]dynamoTypes.AttributeValue{
			"pk": &dynamoTypes.AttributeValueMemberS{Value: "USER#" + req.Msg.Id},
			"sk": &dynamoTypes.AttributeValueMemberS{Value: "USER"},
		},
		UpdateExpression: new("SET #ver = :ver, #mdt = :mdt REMOVE #cor.#coy"),
		ExpressionAttributeNames: map[string]string{
			"#ver": "version",
			"#mdt": "modifiedTime",
			"#sta": "status",
			"#cor": "correlations",
			"#coy": req.Msg.Provider + "#" + req.Msg.Kind,
		},
		ExpressionAttributeValues: map[string]dynamoTypes.AttributeValue{
			":ver": &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatUint(user.Version+1, 10)},
			":mdt": &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(timestamp.UnixMilli(), 10)},
			":dst": &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(int64(identity.UserStatus_USER_STATUS_DELETED), 10)},
		},
		ConditionExpression:                 new("attribute_exists(pk) AND #ver < :ver AND #sta <> :dst AND attribute_exists(#cor.#coy)"),
		ReturnValuesOnConditionCheckFailure: dynamoTypes.ReturnValuesOnConditionCheckFailureAllOld,
	})
	if err != nil {
		var cfe *dynamoTypes.ConditionalCheckFailedException
		if errors.As(err, &cfe) {
			var u dynamo.User
			if len(cfe.Item) > 0 {
				if err := attributevalue.UnmarshalMap(cfe.Item, &u); err != nil {
					logger.Error("Error parsing old user", slog.Any("error", err))
					return nil, connectkit.NewUnexpected()
				}
				if u.Version == user.Version && u.Status != identity.UserStatus_USER_STATUS_DELETED {
					return connect.NewResponse(&identity.RemoveUserCorrelationResponse{User: user}), nil
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

	user.Version++
	user.ModifiedTime = timestamppb.New(timestamp)

	{
		newCorrelations := []*common.Correlation{}
		for _, n := range user.Correlations {
			if req.Msg.Provider != n.Provider || req.Msg.Kind != n.Kind {
				newCorrelations = append(newCorrelations, n)
			}
		}
		user.Correlations = newCorrelations
	}

	{
		event, err := identityevents.BuildUserUpdatedEvent(EventBusName, Namespace, &identityevents.UserUpserted{
			Id:           user.Id,
			Correlations: user.Correlations,
			Version:      user.Version,
			CreatedTime:  user.CreatedTime,
			ModifiedTime: user.ModifiedTime,

			Status:       user.Status,
			EmailAddress: user.EmailAddress,
			Language:     user.Language,
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

	return connect.NewResponse(&identity.RemoveUserCorrelationResponse{User: user}), nil
}
