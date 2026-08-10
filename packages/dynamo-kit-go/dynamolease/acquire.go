package dynamolease

import (
	"context"
	"dynamokit"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Lease struct {
	client   *Client
	key      string
	active   bool
	expires  int64
	TargetId *string
	Context  dynamoTypes.AttributeValue
}

type AcquireInput struct {
	Key      string
	Version  any
	TargetId *string
	Ttl      *time.Time
}

func (client *Client) Acquire(ctx context.Context, params *AcquireInput) (*Lease, error) {
	deadline, ok := ctx.Deadline()
	if !ok {
		return nil, errors.New("Missing context deadline")
	}

	lease := Lease{
		client:  client,
		key:     params.Key,
		active:  true,
		expires: deadline.UnixMilli(),
	}

	conditionExpression := "(attribute_not_exists(lease) OR lease.expires < :now)"
	var updateBuilder dynamokit.UpdateBuilder
	dynamokit.BindValue(&updateBuilder, ":now", time.Now().UnixMilli())

	updateBuilder.SetString("type", "LEASE")

	updateBuilder.SetMap("lease", map[string]dynamoTypes.AttributeValue{
		"expires": &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(lease.expires, 10)},
	})

	if params.Version != nil {
		var v string
		switch vx := params.Version.(type) {
		case int32:
			v = strconv.FormatUint(uint64(vx), 10)
		case int64:
			v = strconv.FormatUint(uint64(vx), 10)
		case uint32:
			v = strconv.FormatUint(uint64(vx), 10)
		case uint64:
			v = strconv.FormatUint(vx, 10)
		default:
			return nil, errors.New("Unsupported version type")
		}

		updateBuilder.AddAttributeName("#ver", "version")
		dynamokit.BindValue(&updateBuilder, ":ver", v)
		conditionExpression += " AND (attribute_not_exists(#ver) OR #ver < :ver)"
	}

	if params.TargetId != nil {
		updateBuilder.SetString("targetId", *params.TargetId)
	}

	if params.Ttl != nil {
		if lease.client.ttlKey != nil {
			updateBuilder.SetInt64(*lease.client.ttlKey, params.Ttl.Unix())
		} else {
			return nil, errors.New("Attempted to set TTL without configured TTL key")
		}
	} else {
		if lease.client.ttlKey != nil {
			updateBuilder.Remove(*lease.client.ttlKey)
		}
	}

	response, err := client.dynamo.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                           &lease.client.tableName,
		Key:                                 lease.client.BuildKey(lease.key),
		UpdateExpression:                    updateBuilder.GetUpdateExpression(),
		ExpressionAttributeNames:            updateBuilder.GetAttributeNames(),
		ExpressionAttributeValues:           updateBuilder.GetAttributeValues(),
		ConditionExpression:                 new(conditionExpression),
		ReturnValues:                        dynamoTypes.ReturnValueAllNew,
		ReturnValuesOnConditionCheckFailure: dynamoTypes.ReturnValuesOnConditionCheckFailureAllOld,
	})
	if err != nil {
		var cfe *dynamoTypes.ConditionalCheckFailedException
		if errors.As(err, &cfe) {
			var item struct {
				Lease *struct {
					Version *uint64 `dynamodbav:"version"`
				} `dynamodbav:"lease"`
			}
			if err := attributevalue.UnmarshalMap(cfe.Item, &item); err != nil {
				return nil, fmt.Errorf("Error parsing old item, %v", err)
			}

			if item.Lease == nil {
				return nil, nil
			} else {
				return nil, errors.New("Error acquiring lease, active lease in place")
			}
		} else {
			return nil, fmt.Errorf("Error acquiring lease, %v", err)
		}
	}

	var item struct {
		TargetId *string `dynamodbav:"targetId"`
	}
	if err = attributevalue.UnmarshalMap(response.Attributes, &item); err != nil {
		return nil, fmt.Errorf("Error parsing lease attributes, %v", err)
	}

	if context, ok := response.Attributes["context"]; ok {
		lease.Context = context
	}

	lease.TargetId = item.TargetId

	return &lease, nil
}
