package main

import (
	"catalogue/src/types/dynamo"
	"connectkit"
	"context"
	"contracts/dist/catalogue/v1"
	"log/slog"

	"connectrpc.com/connect"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (s *Handler) ListSquads(
	ctx context.Context,
	req *connect.Request[catalogue.ListSquadsRequest],
) (*connect.Response[catalogue.ListSquadsResponse], error) {
	logger := connectkit.GetLogger(ctx)
	authCtx := connectkit.GetAuthContext(ctx)

	if authCtx.Lambda == nil {
		return nil, connectkit.NewUnauthorized()
	}

	if authCtx.Lambda != nil {
		resp, err := authCtx.Evaluate(ctx, "catalogue", "list_squads",
			&catalogue.ListSquadsContext{
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

	var squads []dynamo.Squad
	{
		response, err := DynamoRead.Query(ctx, &dynamodb.QueryInput{
			TableName:              TableName,
			IndexName:              new("gsi1"),
			KeyConditionExpression: new("#g1p = :g1p AND begins_with(#g1s, :prefix)"),
			ExpressionAttributeNames: map[string]string{
				"#g1p": "gsi1pk",
				"#g1s": "gsi1sk",
			},
			ExpressionAttributeValues: map[string]dynamoTypes.AttributeValue{
				":g1p":    &dynamoTypes.AttributeValueMemberS{Value: "USER#" + authCtx.Lambda.UserId},
				":prefix": &dynamoTypes.AttributeValueMemberS{Value: "SQUAD#"},
			},
		})
		if err != nil {
			logger.Error("Error retriving squads", slog.Any("error", err))
			return nil, connectkit.NewUnexpected()
		}
		if err := attributevalue.UnmarshalListOfMaps(response.Items, &squads); err != nil {
			logger.Warn("Error parsing squads", slog.Any("items", response.Items), slog.Any("error", err))
			return nil, connectkit.NewUnexpected()
		}
	}

	ss := make([]*catalogue.Squad, len(squads))
	for i, squad := range squads {
		ss[i] = &catalogue.Squad{
			Id:      squad.Id,
			Version: squad.Version,

			Name:     squad.Name,
			Faction:  squad.Faction,
			Loadouts: make([]*catalogue.Loadout, len(squad.Loadouts)),
		}
		for n, loadout := range squad.Loadouts {
			ss[i].Loadouts[n] = Loadouts[loadout]
		}
	}

	return connect.NewResponse(&catalogue.ListSquadsResponse{Squads: ss}), nil
}
