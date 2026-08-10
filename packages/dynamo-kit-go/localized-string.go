package dynamokit

import (
	"contracts/dist/common/v1"
	"strconv"

	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type LocalizedString struct {
	Translations []struct {
		Language common.Language
		Value    string
	}
}

func MarshalLocalizedString(v *common.LocalizedString) *dynamoTypes.AttributeValueMemberM {
	translations := make([]dynamoTypes.AttributeValue, 0, len(v.Translations))
	for _, translation := range v.Translations {
		translations = append(translations, &dynamoTypes.AttributeValueMemberM{Value: map[string]dynamoTypes.AttributeValue{
			"language": &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(int64(translation.Language), 10)},
			"value":    &dynamoTypes.AttributeValueMemberS{Value: translation.Value},
		}})
	}
	return &dynamoTypes.AttributeValueMemberM{Value: map[string]dynamoTypes.AttributeValue{
		"translations": &dynamoTypes.AttributeValueMemberL{Value: translations},
	}}
}

func UnmarshalLocalizedString(v *LocalizedString) *common.LocalizedString {
	var value *common.LocalizedString
	if v != nil {
		value = &common.LocalizedString{Translations: make([]*common.Translation, 0, len(v.Translations))}
		for _, translation := range v.Translations {
			value.Translations = append(value.Translations, &common.Translation{
				Language: translation.Language,
				Value:    translation.Value,
			})
		}
	}
	return value
}

func GetDefault(v *common.LocalizedString) string {
	if v == nil || len(v.Translations) == 0 {
		return ""
	}

	var best *common.Translation
	for _, t := range v.Translations {
		if t.Language == common.Language_LANGUAGE_EN {
			return t.Value
		}
		if t.Language <= 0 {
			continue
		}
		if best == nil || t.Language < best.Language {
			best = t
		}
	}
	if best == nil {
		return ""
	}
	return best.Value
}
