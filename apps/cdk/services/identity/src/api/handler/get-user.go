package main

import (
	"connectkit"
	"context"
	"contracts/dist/identity/v1"
	"log/slog"

	"connectrpc.com/connect"

	dynamo "identity/src/types/dynamo"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (s *Handler) GetUser(
	ctx context.Context,
	req *connect.Request[identity.GetUserRequest],
) (*connect.Response[identity.GetUserResponse], error) {
	logger := connectkit.GetLogger(ctx)
	authCtx := connectkit.GetAuthContext(ctx)

	if req.Msg.Id == nil {
		if authCtx.Lambda != nil {
			req.Msg.Id = &authCtx.Lambda.UserId
		} else {
			return nil, connectkit.NewInvalidArgument("userId")
		}
	}

	var user *identity.User
	{
		response, err := DynamoRead.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: TableName,
			Key: map[string]dynamoTypes.AttributeValue{
				"pk": &dynamoTypes.AttributeValueMemberS{Value: "USER#" + *req.Msg.Id},
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
		resp, err := authCtx.Evaluate(ctx, "identity", "get_user",
			&identity.GetUserContext{
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

	return connect.NewResponse(&identity.GetUserResponse{User: user}), nil
}
