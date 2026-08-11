package main

import (
	"context"
	dynamo "identity/src/types/dynamo"
	"log/slog"
	"strconv"
	"time"

	"connectkit"
	identityevents "contracts/dist/identity/events"
	"contracts/dist/identity/v1"

	"connectrpc.com/connect"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	eventTypes "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Handler) DeleteUser(
	ctx context.Context,
	req *connect.Request[identity.DeleteUserRequest],
) (*connect.Response[identity.DeleteUserResponse], error) {
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
		resp, err := authCtx.Evaluate(ctx, "identity_customer", "delete_user",
			&identity.DeleteUserContext{
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

	{
		_, err := DynamoWrite.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: TableName,
			Key: map[string]dynamoTypes.AttributeValue{
				"pk": &dynamoTypes.AttributeValueMemberS{Value: "USER#" + req.Msg.Id},
				"sk": &dynamoTypes.AttributeValueMemberS{Value: "USER"},
			},
			UpdateExpression: new("SET #ver = :ver, #mdt = :mdt, #sta = :sta REMOVE #g1p, #g1s, #cor, #ead, #nam, #jot, #aid, #lag, #nos"),
			ExpressionAttributeNames: map[string]string{
				"#ver": "version",
				"#mdt": "modifiedTime",
				"#sta": "status",
				"#g1p": "gsi1pk",
				"#g1s": "gsi1sk",
				"#cor": "correlations",
				"#ead": "emailAddress",
				"#nam": "name",
				"#jot": "jobTitle",
				"#aid": "auth0Id",
				"#lag": "language",
				"#nos": "notificationSettings",
			},
			ExpressionAttributeValues: map[string]dynamoTypes.AttributeValue{
				":ver": &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatUint(user.Version+1, 10)},
				":mdt": &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(timestamp.UnixMilli(), 10)},
				":sta": &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(int64(identity.UserStatus_USER_STATUS_DELETED), 10)},
			},
			ConditionExpression: new("attribute_exists(pk) AND #ver < :ver AND #sta <> :sta"),
		})
		if err != nil {
			logger.Error("Error deleting dynamo user", slog.Any("error", err))
			return nil, connectkit.NewUnexpected()
		}
	}

	_, err := Cognito.AdminDeleteUser(ctx, &cognitoidentityprovider.AdminDeleteUserInput{
		UserPoolId: UserPoolId,
		Username:   new(cognitoId),
	})
	if err != nil {
		logger.Error("Error deleting cognito user", slog.Any("error", err))
	}

	user.Version++
	user.ModifiedTime = timestamppb.New(timestamp)

	{
		event, err := identityevents.BuildUserDeletedEvent(EventBusName, Namespace, &identityevents.UserDeleted{
			Id:           user.Id,
			Correlations: user.Correlations,
			Version:      user.Version,
			CreatedTime:  user.CreatedTime,
			ModifiedTime: user.ModifiedTime,
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

	return connect.NewResponse(&identity.DeleteUserResponse{
		Id:           user.Id,
		Correlations: user.Correlations,
		Version:      user.Version,
		CreatedTime:  user.CreatedTime,
		ModifiedTime: user.ModifiedTime,
	}), nil
}
