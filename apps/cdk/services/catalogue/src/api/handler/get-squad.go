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

func (s *Handler) GetSquad(
	ctx context.Context,
	req *connect.Request[catalogue.GetSquadRequest],
) (*connect.Response[catalogue.GetSquadResponse], error) {
	logger := connectkit.GetLogger(ctx)
	authCtx := connectkit.GetAuthContext(ctx)

	var squad *catalogue.Squad
	{
		response, err := DynamoWrite.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: TableName,
			Key: map[string]dynamoTypes.AttributeValue{
				"pk": &dynamoTypes.AttributeValueMemberS{Value: "SQUAD#" + req.Msg.Id},
				"sk": &dynamoTypes.AttributeValueMemberS{Value: "SQUAD"},
			},
		})
		if err != nil {
			logger.Error("Error getting squad", slog.Any("error", err))
			return nil, connectkit.NewUnexpected()
		}
		if len(response.Item) > 0 {
			var item dynamo.Squad
			if err := attributevalue.UnmarshalMap(response.Item, &item); err != nil {
				logger.Error("Error parsing squad", slog.Any("error", err))
				return nil, connectkit.NewUnexpected()
			}

			ls := make([]*catalogue.Loadout, len(item.Loadouts))
			for i, loadout := range item.Loadouts {
				ls[i] = Loadouts[loadout]
			}

			squad = &catalogue.Squad{
				Id:      item.Id,
				Version: item.Version,

				Name:     item.Name,
				Faction:  item.Faction,
				Loadouts: ls,
			}
		}
	}

	if authCtx.Lambda != nil {
		resp, err := authCtx.Evaluate(ctx, "catalogue", "get_squad",
			&catalogue.GetSquadContext{
				Request: req.Msg,
				Subject: squad,
			},
		)
		if err != nil {
			logger.Error("Error evaluating policy", slog.Any("error", err))
			return nil, connectkit.NewUnexpected()
		} else if !resp.Authz {
			return nil, connectkit.NewUnauthorized()
		}
	}

	if squad == nil {
		return nil, connectkit.NewNotFound()
	}

	return connect.NewResponse(&catalogue.GetSquadResponse{Squad: squad}), nil
}
