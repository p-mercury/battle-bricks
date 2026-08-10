package dynamokit

import (
	"strconv"

	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"google.golang.org/genproto/googleapis/type/money"
)

type Money struct {
	CurrencyCode string
	Units        int64
	Nanos        int32
}

func MarshalMoney(m *money.Money) *dynamoTypes.AttributeValueMemberM {
	return &dynamoTypes.AttributeValueMemberM{Value: map[string]dynamoTypes.AttributeValue{
		"currencyCode": &dynamoTypes.AttributeValueMemberS{Value: m.CurrencyCode},
		"units":        &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(m.Units, 10)},
		"nanos":        &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(int64(m.Nanos), 10)},
	}}
}

func UnmarshalMoney(m *Money) *money.Money {
	var out *money.Money
	if m != nil {
		out = &money.Money{
			CurrencyCode: m.CurrencyCode,
			Units:        m.Units,
			Nanos:        m.Nanos,
		}
	}
	return out
}
