package dynamolease

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Client struct {
	dynamo       *dynamodb.Client
	tableName    string
	partitionKey string
	sortKey      *string
	ttlKey       *string
}

func NewClient(dynamo *dynamodb.Client, tableName string, partitionKey string, sortKey *string, ttlKey *string) *Client {
	return new(Client{
		dynamo:       dynamo,
		tableName:    tableName,
		partitionKey: partitionKey,
		sortKey:      sortKey,
		ttlKey:       ttlKey,
	})
}

func (client *Client) BuildKey(key string) map[string]dynamoTypes.AttributeValue {
	k := map[string]dynamoTypes.AttributeValue{
		client.partitionKey: &dynamoTypes.AttributeValueMemberS{Value: key},
	}

	if client.sortKey != nil {
		k[*client.sortKey] = &dynamoTypes.AttributeValueMemberS{Value: "LEASE"}
	}

	return k
}
