package main

import (
	"connectkit"
	"context"
	"contracts/dist/catalogue/v1"
	"errors"
	"log/slog"
	"strconv"
	"time"

	"connectrpc.com/connect"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (s *Handler) CreateSquad(
	ctx context.Context,
	req *connect.Request[catalogue.CreateSquadRequest],
) (*connect.Response[catalogue.CreateSquadResponse], error) {
	logger := connectkit.GetLogger(ctx)
	authCtx := connectkit.GetAuthContext(ctx)
	timestamp := time.Now().Truncate(time.Millisecond)

	if req.Msg.UserId == nil {
		if authCtx.Lambda != nil {
			req.Msg.UserId = &authCtx.Lambda.UserId
		}
	}

	if authCtx.Lambda != nil {
		resp, err := authCtx.Evaluate(ctx, "catalogue", "create_squad",
			&catalogue.CreateSquadContext{
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

	var id string
	for range 4 {
		nid, err := connectkit.NewBase62Id("", 20)
		if err != nil {
			logger.Error("Error generating base62 id", slog.Any("error", err))
			continue
		}

		transactItems := []types.TransactWriteItem{}
		item := map[string]dynamoTypes.AttributeValue{
			"pk":     &dynamoTypes.AttributeValueMemberS{Value: "SQUAD#" + nid},
			"sk":     &dynamoTypes.AttributeValueMemberS{Value: "SQUAD"},
			"gsi1pk": &dynamoTypes.AttributeValueMemberS{Value: "USER#" + *req.Msg.UserId},
			"gsi1sk": &dynamoTypes.AttributeValueMemberS{Value: "SQUAD#" + nid},
			"type":   &dynamoTypes.AttributeValueMemberS{Value: "SQUAD"},

			"id":           &dynamoTypes.AttributeValueMemberS{Value: nid},
			"version":      &dynamoTypes.AttributeValueMemberN{Value: "1"},
			"createdTime":  &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(timestamp.UnixMilli(), 10)},
			"modifiedTime": &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(timestamp.UnixMilli(), 10)},

			"name":    &dynamoTypes.AttributeValueMemberS{Value: req.Msg.Name},
			"faction": &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(int64(req.Msg.Faction), 10)},
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
								"pk":   &dynamoTypes.AttributeValueMemberS{Value: "CORRELATION#SQUAD#" + correlation.Provider + "#" + correlation.Kind + "#" + correlation.Id},
								"sk":   &dynamoTypes.AttributeValueMemberS{Value: "CORRELATION"},
								"type": &dynamoTypes.AttributeValueMemberS{Value: "CORRELATION"},

								"provider": &dynamoTypes.AttributeValueMemberS{Value: correlation.Provider},
								"kind":     &dynamoTypes.AttributeValueMemberS{Value: correlation.Kind},
								"id":       &dynamoTypes.AttributeValueMemberS{Value: correlation.Id},
								"targetId": &dynamoTypes.AttributeValueMemberS{Value: nid},
							},
							ConditionExpression: new("attribute_not_exists(pk)"),
						},
					})
				}
			}
			item["correlations"] = &dynamoTypes.AttributeValueMemberM{Value: correlations}
		}

		{
			loadouts := make([]dynamoTypes.AttributeValue, len(req.Msg.Loadouts))
			for i, loadout := range req.Msg.Loadouts {
				loadouts[i] = &dynamoTypes.AttributeValueMemberS{Value: loadout}
			}
			item["loadouts"] = &dynamoTypes.AttributeValueMemberL{Value: loadouts}
		}

		transactItems = append(transactItems, types.TransactWriteItem{
			Put: &types.Put{
				TableName:                           TableName,
				Item:                                item,
				ConditionExpression:                 new("attribute_not_exists(pk)"),
				ReturnValuesOnConditionCheckFailure: dynamoTypes.ReturnValuesOnConditionCheckFailureAllOld,
			},
		})

		if _, err := DynamoWrite.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{TransactItems: transactItems}); err != nil {
			var cfe *dynamoTypes.TransactionCanceledException
			if errors.As(err, &cfe) {
				if cfe.CancellationReasons[len(cfe.CancellationReasons)-1].Item != nil {
					continue
				} else {
					return nil, connectkit.NewUnexpected()
				}
			} else {
				logger.Error("Error creating product", slog.Any("error", err))
				return nil, connectkit.NewUnexpected()
			}
		}

		id = nid
		break
	}

	if id == "" {
		logger.Error("Can't find free id")
		return nil, connectkit.NewUnexpected()
	}

	ls := make([]*catalogue.Loadout, len(req.Msg.Loadouts))
	for i, loadout := range req.Msg.Loadouts {
		ls[i] = Loadouts[loadout]
	}

	return connect.NewResponse(&catalogue.CreateSquadResponse{
		Squad: &catalogue.Squad{
			Id:      id,
			Version: 1,

			Name:     req.Msg.Name,
			Faction:  req.Msg.Faction,
			Loadouts: ls,
		},
	}), nil
}
