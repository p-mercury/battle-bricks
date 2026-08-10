package main

import (
	"connectkit"
	"context"
	"errors"
	"log/slog"

	"contracts/dist/identity/v1"
	"contracts/dist/policy/v1"
	"policy/src/types/dynamo"

	"connectrpc.com/connect"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (s *Handler) GetUser(
	ctx context.Context,
	req *connect.Request[policy.GetUserRequest],
) (*connect.Response[policy.GetUserResponse], error) {
	logger := connectkit.GetLogger(ctx)
	authCtx := connectkit.GetAuthContext(ctx)

	if authCtx.Lambda == nil {
		logger.Error("Missin authorizer")
		return nil, connectkit.NewUnexpected()
	}

	resp := &policy.GetUserResponse{}
	{
		response, err := DynamoRead.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: TableName,
			Key: map[string]dynamoTypes.AttributeValue{
				"pk": &dynamoTypes.AttributeValueMemberS{Value: authCtx.Lambda.UserId},
			}})
		if err != nil {
			logger.Error("Error getting user", slog.String("id", authCtx.Lambda.UserId), slog.Any("error", err))
			return nil, connectkit.NewUnexpected()
		}
		if response.Item == nil {
			logger.Warn("User not found", slog.String("id", authCtx.Lambda.UserId))
			return nil, connect.NewError(
				connect.CodeNotFound,
				errors.New("User not found"))
		}
		if _, deleted := response.Item["ttl"]; deleted {
			logger.Warn("User is deleted", slog.String("id", authCtx.Lambda.UserId))
			return nil, connect.NewError(
				connect.CodeNotFound,
				errors.New("User not found"))
		}

		var u *dynamo.CustomerUser
		if err = attributevalue.UnmarshalMap(response.Item, &u); err != nil {
			logger.Error("Error parsing item of type 'USER'", slog.Any("item", response.Item), slog.Any("error", err))
			return nil, connectkit.NewUnexpected()
		}

		resp.User = &policy.User{
			Id:     u.Pk,
			Status: identity.UserStatus(u.Status),
		}
	}

	return connect.NewResponse(resp), nil
}
