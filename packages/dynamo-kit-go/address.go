package dynamokit

import (
	"contracts/dist/common/v1"

	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Address struct {
	RegionCode         string
	AddressLines       []string
	PostalCode         *string
	LanguageCode       *string
	SortingCode        *string
	AdministrativeArea *string
	Locality           *string
	Sublocality        *string
}

func MarshalAddress(a *common.Address) *dynamoTypes.AttributeValueMemberM {
	address := map[string]dynamoTypes.AttributeValue{
		"regionCode": &dynamoTypes.AttributeValueMemberS{Value: a.RegionCode},
	}

	{
		lines := make([]dynamoTypes.AttributeValue, len(a.AddressLines))
		for i, line := range a.AddressLines {
			lines[i] = &dynamoTypes.AttributeValueMemberS{Value: line}
		}
		address["addressLines"] = &dynamoTypes.AttributeValueMemberL{Value: lines}
	}

	if a.PostalCode != nil {
		address["postalCode"] = &dynamoTypes.AttributeValueMemberS{Value: *a.PostalCode}
	}
	if a.LanguageCode != nil {
		address["languageCode"] = &dynamoTypes.AttributeValueMemberS{Value: *a.LanguageCode}
	}
	if a.SortingCode != nil {
		address["sortingCode"] = &dynamoTypes.AttributeValueMemberS{Value: *a.SortingCode}
	}
	if a.AdministrativeArea != nil {
		address["administrativeArea"] = &dynamoTypes.AttributeValueMemberS{Value: *a.AdministrativeArea}
	}
	if a.Locality != nil {
		address["locality"] = &dynamoTypes.AttributeValueMemberS{Value: *a.Locality}
	}
	if a.Sublocality != nil {
		address["sublocality"] = &dynamoTypes.AttributeValueMemberS{Value: *a.Sublocality}
	}

	return &dynamoTypes.AttributeValueMemberM{Value: address}
}

func UnmarshalAddress(a *Address) *common.Address {
	var address *common.Address
	if a != nil {
		address = &common.Address{
			RegionCode:         a.RegionCode,
			AddressLines:       a.AddressLines,
			PostalCode:         a.PostalCode,
			LanguageCode:       a.LanguageCode,
			SortingCode:        a.SortingCode,
			AdministrativeArea: a.AdministrativeArea,
			Locality:           a.Locality,
			Sublocality:        a.Sublocality,
		}
	}
	return address
}
