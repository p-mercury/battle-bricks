package dynamolease

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (lease *Lease) Delete(ctx context.Context) error {
	if lease == nil {
		return nil
	}

	_, err := lease.client.dynamo.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &lease.client.tableName,
		Item: map[string]dynamoTypes.AttributeValue{
			"pk":  &dynamoTypes.AttributeValueMemberS{Value: lease.key},
			"sk":  &dynamoTypes.AttributeValueMemberS{Value: "LEASE"},
			"ttl": &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(time.Now().Unix()+604800, 10)},
		},
	})
	if err != nil {
		return fmt.Errorf("Error deleting mapping, %v", err)
	}

	lease.active = false

	return nil
}
