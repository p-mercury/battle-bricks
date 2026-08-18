package main

import (
	"catalogue/src/types/dynamo"
	"connectkit"
	"context"
	"contracts/dist/catalogue/v1"
	"dynamokit"
	"log/slog"
	"slices"
	"time"

	"connectrpc.com/connect"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Handler) UpdateSquad(
	ctx context.Context,
	req *connect.Request[catalogue.UpdateSquadRequest],
) (*connect.Response[catalogue.UpdateSquadResponse], error) {
	logger := connectkit.GetLogger(ctx)
	authCtx := connectkit.GetAuthContext(ctx)
	timestamp := time.Now().Truncate(time.Millisecond)

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

			ls := make([]*catalogue.Unit, len(item.Loadouts))
			for i, loadout := range item.Loadouts {
				ls[i] = Units[loadout]
			}

			squad = &catalogue.Squad{
				Id:      item.Id,
				Version: item.Version,

				Name:    item.Name,
				Faction: item.Faction,
				Units:   ls,
			}
		}
	}

	if authCtx.Lambda != nil {
		resp, err := authCtx.Evaluate(ctx, "catalogue", "update_squad",
			&catalogue.UpdateSquadContext{
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

	if req.Msg.Version != nil && squad.Version != *req.Msg.Version {
		return nil, connectkit.NewConflict("Version mismatch")
	}

	{
		transactItems := []types.TransactWriteItem{}
		var updateBuilder dynamokit.UpdateBuilder

		updateBuilder.Set("#ver", ":ver")
		updateBuilder.AddAttributeName("#ver", "version")
		dynamokit.BindValue(&updateBuilder, ":ver", squad.Version+1)

		updateBuilder.Set("#mdt", ":mdt")
		updateBuilder.AddAttributeName("#mdt", "modifiedTime")
		dynamokit.BindValue(&updateBuilder, ":mdt", timestamp.UnixMilli())

		// if slices.Contains(req.Msg.UpdateMask, "correlations") {
		// 	if len(req.Msg.Correlations) < 1 {
		// 		return nil, connectkit.NewInvalidArgument("correlations")
		// 	}

		// 	changed := false
		// 	for i, correlation := range req.Msg.Correlations {
		// 		{
		// 			exists := false
		// 			for n := range squad.Correlations {
		// 				if correlation.Provider == squad.Correlations[n].Provider && correlation.Kind == squad.Correlations[n].Kind {
		// 					exists = true
		// 					if correlation.Id != squad.Correlations[n].Id {
		// 						squad.Correlations[n].Id = correlation.Id
		// 						changed = true
		// 					}
		// 					break
		// 				}
		// 			}
		// 			if !exists {
		// 				squad.Correlations = append(squad.Correlations, correlation)
		// 				changed = true
		// 			}
		// 		}

		// 		if changed {
		// 			transactItems = append(transactItems, types.TransactWriteItem{
		// 				Put: &types.Put{
		// 					TableName: TableName,
		// 					Item: map[string]dynamoTypes.AttributeValue{
		// 						"pk":   &dynamoTypes.AttributeValueMemberS{Value: "CORRELATION#PRODUCT#" + correlation.Provider + "#" + correlation.Kind + "#" + correlation.Id},
		// 						"sk":   &dynamoTypes.AttributeValueMemberS{Value: "CORRELATION"},
		// 						"type": &dynamoTypes.AttributeValueMemberS{Value: "CORRELATION"},

		// 						"provider": &dynamoTypes.AttributeValueMemberS{Value: correlation.Provider},
		// 						"kind":     &dynamoTypes.AttributeValueMemberS{Value: correlation.Kind},
		// 						"id":       &dynamoTypes.AttributeValueMemberS{Value: correlation.Id},
		// 						"targetId": &dynamoTypes.AttributeValueMemberS{Value: req.Msg.Id},
		// 					},
		// 					ExpressionAttributeNames: map[string]string{
		// 						"#tid": "targetId",
		// 					},
		// 					ExpressionAttributeValues: map[string]dynamoTypes.AttributeValue{
		// 						":tid": &dynamoTypes.AttributeValueMemberS{Value: req.Msg.Id},
		// 					},
		// 					ConditionExpression: new("attribute_not_exists(pk) OR #tid = :tid"),
		// 				},
		// 			})

		// 			updateBuilder.AddAttributeName("#cor", "correlations")
		// 			updateBuilder.Set(fmt.Sprintf("#cor.#co%d", i), fmt.Sprintf(":co%d", i))
		// 			updateBuilder.AddAttributeName(fmt.Sprintf("#co%d", i), correlation.Provider+"#"+correlation.Kind)
		// 			dynamokit.BindValue(&updateBuilder, fmt.Sprintf(":co%d", i), map[string]dynamoTypes.AttributeValue{
		// 				"provider": &dynamoTypes.AttributeValueMemberS{Value: correlation.Provider},
		// 				"kind":     &dynamoTypes.AttributeValueMemberS{Value: correlation.Kind},
		// 				"id":       &dynamoTypes.AttributeValueMemberS{Value: correlation.Id},
		// 			})
		// 		}
		// 	}
		// 	if changed {
		// 		changedFields = append(changedFields, catalogueevents.ProductUpserted_CHANGED_FIELD_CORRELATIONS)
		// 	}
		// }

		if slices.Contains(req.Msg.UpdateMask, "name") {
			if req.Msg.Name == nil {
				return nil, connectkit.NewInvalidArgument("name")
			}
			if squad.Name != *req.Msg.Name {
				squad.Name = *req.Msg.Name
				updateBuilder.SetString("name", squad.Name)
			}
		}

		if slices.Contains(req.Msg.UpdateMask, "faction") {
			if req.Msg.Faction == nil {
				return nil, connectkit.NewInvalidArgument("faction")
			}
			if squad.Faction != *req.Msg.Faction {
				squad.Faction = *req.Msg.Faction
				updateBuilder.SetInt64("name", int64(squad.Faction))
			}
		}

		if slices.Contains(req.Msg.UpdateMask, "units") {
			ls := make([]*catalogue.Unit, len(req.Msg.Units))
			for i, loadout := range req.Msg.Units {
				ls[i] = Units[loadout]
			}

			if !slices.Equal(squad.Units, ls) {
				squad.Units = ls

				loadouts := make([]dynamoTypes.AttributeValue, len(req.Msg.Units))
				for i, loadout := range req.Msg.Units {
					loadouts[i] = &dynamoTypes.AttributeValueMemberS{Value: loadout}
				}
				updateBuilder.SetList("loadouts", loadouts)
			}
		}

		transactItems = append(transactItems, types.TransactWriteItem{
			Update: &types.Update{
				TableName: TableName,
				Key: map[string]dynamoTypes.AttributeValue{
					"pk": &dynamoTypes.AttributeValueMemberS{Value: "SQUAD#" + req.Msg.Id},
					"sk": &dynamoTypes.AttributeValueMemberS{Value: "SQUAD"},
				},
				UpdateExpression:                    updateBuilder.GetUpdateExpression(),
				ExpressionAttributeNames:            updateBuilder.GetAttributeNames(),
				ExpressionAttributeValues:           updateBuilder.GetAttributeValues(),
				ConditionExpression:                 new("attribute_exists(pk) AND #ver < :ver"),
				ReturnValuesOnConditionCheckFailure: dynamoTypes.ReturnValuesOnConditionCheckFailureAllOld,
			},
		})
		if _, err := DynamoWrite.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{TransactItems: transactItems}); err != nil {
			logger.Error("Error updating product", slog.Any("error", err))
			return nil, connectkit.NewUnexpected()
		}
	}

	squad.Version++
	squad.ModifiedTime = timestamppb.New(timestamp)

	return connect.NewResponse(&catalogue.UpdateSquadResponse{Squad: squad}), nil
}
