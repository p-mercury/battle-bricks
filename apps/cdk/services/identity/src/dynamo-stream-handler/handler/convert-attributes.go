package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func FromDynamoDBEventAttributeValueMap(m map[string]events.DynamoDBAttributeValue) (res map[string]types.AttributeValue, err error) {
	res = make(map[string]types.AttributeValue, len(m))
	for k, v := range m {
		res[k], err = FromDynamoDBEventAttributeValue(v)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func FromDynamoDBEventAttributeValueList(l []events.DynamoDBAttributeValue) (res []types.AttributeValue, err error) {
	res = make([]types.AttributeValue, len(l))
	for i, v := range l {
		res[i], err = FromDynamoDBEventAttributeValue(v)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func FromDynamoDBEventAttributeValue(av events.DynamoDBAttributeValue) (types.AttributeValue, error) {
	switch av.DataType() {
	case events.DataTypeNull:
		return &types.AttributeValueMemberNULL{Value: true}, nil
	case events.DataTypeString:
		return &types.AttributeValueMemberS{Value: av.String()}, nil
	case events.DataTypeBoolean:
		return &types.AttributeValueMemberBOOL{Value: av.Boolean()}, nil
	case events.DataTypeNumber:
		return &types.AttributeValueMemberN{Value: av.Number()}, nil
	case events.DataTypeBinary:
		return &types.AttributeValueMemberB{Value: av.Binary()}, nil
	case events.DataTypeStringSet:
		return &types.AttributeValueMemberSS{Value: av.StringSet()}, nil
	case events.DataTypeNumberSet:
		return &types.AttributeValueMemberNS{Value: av.NumberSet()}, nil
	case events.DataTypeBinarySet:
		return &types.AttributeValueMemberBS{Value: av.BinarySet()}, nil
	case events.DataTypeMap:
		m := make(map[string]types.AttributeValue)
		for k, v := range av.Map() {
			o, err := FromDynamoDBEventAttributeValue(v)
			if err != nil {
				return nil, err
			}
			m[k] = o
		}
		return &types.AttributeValueMemberM{Value: m}, nil
	case events.DataTypeList:
		l := make([]types.AttributeValue, len(av.List()))
		for i, v := range av.List() {
			o, err := FromDynamoDBEventAttributeValue(v)
			if err != nil {
				return nil, err
			}
			l[i] = o
		}
		return &types.AttributeValueMemberL{Value: l}, nil
	default:
		return nil, fmt.Errorf("unknown AttributeValue union member, %T", av)
	}
}
