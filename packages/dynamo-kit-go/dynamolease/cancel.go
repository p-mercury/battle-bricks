package dynamolease

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (lease *Lease) Cancel(ctx context.Context) error {
	if lease == nil {
		return nil
	}

	if !lease.active {
		return nil
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return errors.New("Missing context deadline")
	}

	_, err := lease.client.dynamo.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:        &lease.client.tableName,
		Key:              lease.client.BuildKey(lease.key),
		UpdateExpression: new("REMOVE #les"),
		ExpressionAttributeNames: map[string]string{
			"#les": "lease",
			"#tol": "ttl",
		},
		ExpressionAttributeValues: map[string]dynamoTypes.AttributeValue{
			":exp": &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(deadline.UnixMilli(), 10)},
		},
		ConditionExpression: new("attribute_not_exists(#tol) AND attribute_exists(#les) AND #les.expires = :exp"),
	})
	if err != nil {
		return fmt.Errorf("Error releasing lease, %v", err)
	}

	lease.active = false

	return nil
}
