package dynamolease

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (lease *Lease) Unlink(ctx context.Context) error {
	if lease == nil {
		return nil
	}

	_, err := lease.client.dynamo.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:        &lease.client.tableName,
		Key:              lease.client.BuildKey(lease.key),
		UpdateExpression: new("REMOVE #les, #tai"),
		ExpressionAttributeNames: map[string]string{
			"#les": "lease",
			"#tai": "targetId",
			"#tol": "ttl",
		},
		ExpressionAttributeValues: map[string]dynamoTypes.AttributeValue{
			":exp": &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(lease.expires, 10)},
		},
		ConditionExpression: new("attribute_not_exists(#tol) AND attribute_exists(#les) AND #les.expires = :exp"),
	})
	if err != nil {
		return fmt.Errorf("Error releasing unlinking, %v", err)
	}

	lease.active = false
	return nil
}
