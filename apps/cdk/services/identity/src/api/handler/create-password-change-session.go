package main

import (
	"context"
	dynamo "identity/src/types/dynamo"
	"log/slog"

	"connectkit"
	"contracts/dist/identity/v1"

	"connectrpc.com/connect"
	"github.com/auth0/go-auth0/management"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (s *Handler) CreatePasswordChangeSession(
	ctx context.Context,
	req *connect.Request[identity.CreatePasswordChangeSessionRequest],
) (*connect.Response[identity.CreatePasswordChangeSessionResponse], error) {
	logger := connectkit.GetLogger(ctx)
	authCtx := connectkit.GetAuthContext(ctx)

	if authCtx.Iam == nil {
		return nil, connectkit.NewUnauthorized()
	}

	var user dynamo.User
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
		if len(response.Item) < 1 {
			return nil, connectkit.NewNotFound()
		}

		if err = attributevalue.UnmarshalMap(response.Item, &user); err != nil {
			logger.Error("Error parsing user", slog.Any("error", err))
			return nil, connectkit.NewUnexpected()
		}
	}

	auth0Ticket := &management.Ticket{
		ClientID:               Auth0ClientId,
		MarkEmailAsVerified:    new(true),
		IncludeEmailInRedirect: new(true),
		UserID:                 &user.Auth0Id,
	}

	if err := Auth0.Ticket.ChangePassword(ctx, auth0Ticket); err != nil {
		logger.Error("Error creating change password ticket", slog.Any("error", err))
		return nil, connectkit.NewUnexpected()
	}

	return connect.NewResponse(&identity.CreatePasswordChangeSessionResponse{
		Session: &identity.PasswordChangeSession{
			Url: *auth0Ticket.Ticket},
	}), nil
}
