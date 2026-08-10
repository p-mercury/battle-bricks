package main

import (
	"connectkit"
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"contracts/dist/identity/v1"
	"contracts/dist/policy/v1"
	"policy/src/types/dynamo"

	"connectrpc.com/connect"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/open-policy-agent/opa/v1/sdk"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

func (s *Handler) Evaluate(
	ctx context.Context,
	req *connect.Request[policy.EvaluateRequest],
) (*connect.Response[policy.EvaluateResponse], error) {
	logger := connectkit.GetLogger(ctx)
	authCtx := connectkit.GetAuthContext(ctx)

	if authCtx.Iam == nil {
		return nil, connectkit.NewUnauthorized()
	}

	context := map[string]any{}
	if req.Msg.Context != nil {
		if err := json.Unmarshal(req.Msg.Context, &context); err != nil {
			logger.Error("Error parsing context", slog.Any("error", err))
			return nil, connectkit.NewUnexpected()
		}
	}

	var user map[string]any
	{
		rawUser := &policy.GetUserResponse{}

		response, err := DynamoRead.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: TableName,
			Key: map[string]dynamoTypes.AttributeValue{
				"pk": &dynamoTypes.AttributeValueMemberS{Value: req.Msg.UserId},
			}})
		if err != nil {
			logger.Error("Error getting user", slog.String("id", req.Msg.UserId), slog.Any("error", err))
			return nil, connectkit.NewUnexpected()
		}
		if response.Item == nil {
			logger.Warn("User not found", slog.String("id", req.Msg.UserId))
			return nil, connect.NewError(
				connect.CodeNotFound,
				errors.New("User not found"))
		}
		if _, deleted := response.Item["ttl"]; deleted {
			logger.Warn("User is deleted", slog.String("id", req.Msg.UserId))
			return nil, connect.NewError(
				connect.CodeNotFound,
				errors.New("User not found"))
		}

		var u *dynamo.CustomerUser
		if err = attributevalue.UnmarshalMap(response.Item, &u); err != nil {
			logger.Error("Error parsing item of type 'USER'", slog.Any("item", response.Item), slog.Any("error", err))
			return nil, connectkit.NewUnexpected()
		}

		rawUser.User = &policy.User{
			Id:     u.Pk,
			Status: identity.UserStatus(u.Status),
		}

		jsonString, err := protojson.Marshal(rawUser)
		if err != nil {
			logger.Error("Error json marshaling user", slog.Any("error", err))
			return nil, connectkit.NewUnexpected()
		}

		json.Unmarshal(jsonString, &user)
	}

	path := "services/" + req.Msg.Service + "/" + req.Msg.Method + "/decision"

	logger = logger.With(slog.Any("input", map[string]any{
		"user":    user,
		"context": context,
	}), slog.String("path", path))

	result, err := Opa.Decision(ctx, sdk.DecisionOptions{
		Path: path,
		Input: map[string]any{
			"user":    user,
			"context": context,
		},
	})
	if err != nil {
		logger.Error("Error evaluating policy", slog.Any("error", err))
		return nil, connectkit.NewUnexpected()
	}

	logger.Info("Decision", slog.String("decisionId", result.ID), slog.Any("output", result.Result))

	outputMask := map[string]*structpb.Value{}
	if rawMask, exists := result.Result.(map[string]any)["output_mask"]; exists {
		maskMap, ok := rawMask.(map[string]any)
		if !ok {
			logger.Warn("output_mask present but not a map", slog.Any("output_mask", rawMask))
		} else {
			for k, v := range maskMap {
				sv, err := structpb.NewValue(v)
				if err != nil {
					logger.Warn("Error converting input mask value", slog.Any("error", err))
					outputMask[k] = structpb.NewBoolValue(false)
				} else {
					outputMask[k] = sv
				}
			}
		}
	}

	return connect.NewResponse(&policy.EvaluateResponse{
		Id:         result.ID,
		Authz:      result.Result.(map[string]any)["authz"].(bool),
		OutputMask: outputMask,
	}), nil
}
