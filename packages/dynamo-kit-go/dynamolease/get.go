package dynamolease

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type GetInput struct {
	Key string
}

func (client *Client) Get(ctx context.Context, params *GetInput) (*Lease, error) {
	lease := Lease{
		client:  client,
		key:     params.Key,
		active:  false,
		expires: 0,
	}

	response, err := client.dynamo.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &lease.client.tableName,
		Key:       client.BuildKey(params.Key),
	})
	if err != nil {
		return nil, fmt.Errorf("Error getting mapping from dynamo, %v", err)
	}

	var item struct {
		TargetId *string `dynamodbav:"targetId"`
	}
	if err = attributevalue.UnmarshalMap(response.Item, &item); err != nil {
		return nil, fmt.Errorf("Error parsing lease attributes, %v", err)
	}

	if context, ok := response.Item["context"]; ok {
		lease.Context = context
	}

	lease.TargetId = item.TargetId

	return &lease, nil
}
