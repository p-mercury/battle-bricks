package dynamokit

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func MarshalPageToken(lastEvaluatedKey map[string]dynamoTypes.AttributeValue) (*string, error) {
	item := make(map[string]any, len(lastEvaluatedKey))
	for key, value := range lastEvaluatedKey {
		switch av := value.(type) {
		case *dynamoTypes.AttributeValueMemberS:
			item[key] = map[string]string{"S": av.Value}
		case *dynamoTypes.AttributeValueMemberN:
			item[key] = map[string]string{"N": av.Value}
		case *dynamoTypes.AttributeValueMemberB:
			item[key] = map[string][]byte{"B": av.Value}
		default:
			return nil, errors.New("invalid attribute value type")
		}
	}

	b, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}

	return new(base64.StdEncoding.EncodeToString(b)), nil
}

func UnmarshalPageToken(token string) (map[string]dynamoTypes.AttributeValue, error) {
	b, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}

	var raw map[string]map[string]string
	if err := json.Unmarshal(b, &raw); err != nil {
		return nil, err
	}

	result := make(map[string]dynamoTypes.AttributeValue, len(raw))
	for k, typed := range raw {
		if s, ok := typed["S"]; ok {
			result[k] = &dynamoTypes.AttributeValueMemberS{Value: s}
		} else if n, ok := typed["N"]; ok {
			result[k] = &dynamoTypes.AttributeValueMemberN{Value: n}
		} else if b, ok := typed["B"]; ok {
			result[k] = &dynamoTypes.AttributeValueMemberB{Value: []byte(b)}
		} else {
			return nil, errors.New("invalid attribute value type")
		}
	}

	return result, nil
}
