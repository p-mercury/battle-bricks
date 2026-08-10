package dynamolease

import (
	"context"
	"dynamokit"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type LockLeaseInput struct {
	Version  any
	TargetId *string
	Ttl      *time.Time
	Context  map[string]dynamoTypes.AttributeValue
}

func (lease *Lease) Lock(ctx context.Context, params *LockLeaseInput) error {
	if lease == nil {
		return nil
	}

	var updateBuilder dynamokit.UpdateBuilder
	updateBuilder.Remove("#les")
	updateBuilder.AddAttributeName("#les", "lease")
	dynamokit.BindValue(&updateBuilder, ":exp", lease.expires)

	if params.TargetId != nil {
		updateBuilder.SetString("targetId", *params.TargetId)
	} else {
		updateBuilder.Remove("targetId")
	}

	if params.Ttl != nil {
		if lease.client.ttlKey != nil {
			updateBuilder.SetInt64(*lease.client.ttlKey, params.Ttl.Unix())
		} else {
			return errors.New("Attempted to set TTL without configured TTL key")
		}
	} else {
		if lease.client.ttlKey != nil {
			updateBuilder.Remove(*lease.client.ttlKey)
		}
	}

	if params.Context != nil {
		updateBuilder.SetMap("context", params.Context)
	} else {
		updateBuilder.Remove("context")
	}

	switch vx := params.Version.(type) {
	case int32:
		updateBuilder.SetInt32("version", vx)
	case int64:
		updateBuilder.SetInt64("version", vx)
	case uint32:
		updateBuilder.SetUint32("version", vx)
	case uint64:
		updateBuilder.SetUint64("version", vx)
	default:
		return errors.New("Unsupported version type")
	}

	_, err := lease.client.dynamo.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 &lease.client.tableName,
		Key:                       lease.client.BuildKey(lease.key),
		UpdateExpression:          updateBuilder.GetUpdateExpression(),
		ExpressionAttributeNames:  updateBuilder.GetAttributeNames(),
		ExpressionAttributeValues: updateBuilder.GetAttributeValues(),
		ConditionExpression:       new("attribute_exists(#les) AND #les.expires = :exp"),
	})
	if err != nil {
		return fmt.Errorf("Error releasing lease, %v", err)
	}

	lease.active = false

	return nil
}
