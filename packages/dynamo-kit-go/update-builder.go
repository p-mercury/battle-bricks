package dynamokit

import (
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"

	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type UpdateBuilder struct {
	set             map[string]string
	remove          map[string]any
	add             []string
	delete          []string
	attributeNames  map[string]string
	attributeValues map[string]dynamoTypes.AttributeValue
}

func fieldAlias(field string) string {
	h := fnv.New32a()
	h.Write([]byte(strings.ToUpper(field)))
	return fmt.Sprintf("#f%x", h.Sum32())
}

func valueAlias(path string) string {
	h := fnv.New32a()
	h.Write([]byte(path))
	return fmt.Sprintf(":v%x", h.Sum32())
}

func (builder *UpdateBuilder) registerPath(path string) (exprPath, valueKey string) {
	keys := strings.Split(strings.TrimSpace(path), ".")
	aliases := make([]string, len(keys))
	for i, key := range keys {
		if strings.HasPrefix(key, "#") {
			aliases[i] = key
		} else if IsReservedKey(key) {
			aliases[i] = fieldAlias(key)
			builder.AddAttributeName(aliases[i], key)
		} else {
			aliases[i] = key
		}
	}
	return strings.Join(aliases, "."), valueAlias(path)
}

func (builder *UpdateBuilder) Set(key string, value string) {
	if builder.set == nil {
		builder.set = map[string]string{}
	}
	builder.set[strings.TrimSpace(key)] = strings.TrimSpace(value)
}

func (builder *UpdateBuilder) SetAttributeValue(path string, value dynamoTypes.AttributeValue) {
	exprPath, valueKey := builder.registerPath(path)
	builder.Set(exprPath, valueKey)
	builder.AddAttributeValue(valueKey, value)
}

func (builder *UpdateBuilder) SetNull(path string) {
	builder.SetAttributeValue(path, &dynamoTypes.AttributeValueMemberNULL{Value: true})
}

func (builder *UpdateBuilder) SetString(path string, value string) {
	builder.SetAttributeValue(path, &dynamoTypes.AttributeValueMemberS{Value: value})
}

func (builder *UpdateBuilder) SetUint32(path string, value uint32) {
	builder.SetAttributeValue(path, &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatUint(uint64(value), 10)})
}

func (builder *UpdateBuilder) SetInt32(path string, value int32) {
	builder.SetAttributeValue(path, &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(int64(value), 10)})
}

func (builder *UpdateBuilder) SetUint64(path string, value uint64) {
	builder.SetAttributeValue(path, &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatUint(value, 10)})
}

func (builder *UpdateBuilder) SetInt64(path string, value int64) {
	builder.SetAttributeValue(path, &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(value, 10)})
}

func (builder *UpdateBuilder) SetBool(path string, value bool) {
	builder.SetAttributeValue(path, &dynamoTypes.AttributeValueMemberBOOL{Value: value})
}

func (builder *UpdateBuilder) SetList(path string, value []dynamoTypes.AttributeValue) {
	builder.SetAttributeValue(path, &dynamoTypes.AttributeValueMemberL{Value: value})
}

func (builder *UpdateBuilder) SetMap(path string, value map[string]dynamoTypes.AttributeValue) {
	builder.SetAttributeValue(path, &dynamoTypes.AttributeValueMemberM{Value: value})
}

func (builder *UpdateBuilder) SetStringSet(path string, value []string) {
	builder.SetAttributeValue(path, &dynamoTypes.AttributeValueMemberSS{Value: value})
}

func (builder *UpdateBuilder) SetNumberSet(path string, value []string) {
	builder.SetAttributeValue(path, &dynamoTypes.AttributeValueMemberNS{Value: value})
}

func (builder *UpdateBuilder) Remove(path string) {
	if builder.remove == nil {
		builder.remove = map[string]any{}
	}
	exprPath, _ := builder.registerPath(path)
	builder.remove[exprPath] = nil
}

func (builder *UpdateBuilder) Add(value string) {
	builder.add = append(builder.add, strings.TrimSpace(value))
}

func (builder *UpdateBuilder) Delete(path, subset string) {
	builder.delete = append(builder.delete, strings.TrimSpace(path)+" "+strings.TrimSpace(subset))
}

func (builder *UpdateBuilder) GetUpdateExpression() *string {
	var updateExpression []string

	if len(builder.set) > 0 {
		parts := make([]string, 0, len(builder.set))
		for path, valueKey := range builder.set {
			parts = append(parts, path+" = "+valueKey)
		}
		updateExpression = append(updateExpression, "SET")
		updateExpression = append(updateExpression, strings.Join(parts, ", "))
	}

	if len(builder.remove) > 0 {
		parts := make([]string, 0, len(builder.remove))
		for path := range builder.remove {
			parts = append(parts, path)
		}
		updateExpression = append(updateExpression, "REMOVE")
		updateExpression = append(updateExpression, strings.Join(parts, ", "))
	}

	if len(builder.add) > 0 {
		updateExpression = append(updateExpression, "ADD")
		updateExpression = append(updateExpression, strings.Join(builder.add, ", "))
	}

	if len(builder.delete) > 0 {
		updateExpression = append(updateExpression, "DELETE")
		updateExpression = append(updateExpression, strings.Join(builder.delete, ", "))
	}

	if len(updateExpression) == 0 {
		return nil
	}

	return new(strings.Join(updateExpression, " "))
}

type BindableValue interface {
	string | bool |
		int | int32 | int64 |
		uint | uint32 | uint64 |
		float32 | float64 |
		[]dynamoTypes.AttributeValue |
		map[string]dynamoTypes.AttributeValue |
		[]string | []int32 | []int64 |
		[]uint32 | []uint64 |
		[]float32 | []float64
}

func BindValue[T BindableValue](
	builder *UpdateBuilder,
	key string,
	value T,
) {
	var attributeValue dynamoTypes.AttributeValue

	switch value := any(value).(type) {
	case string:
		attributeValue = &dynamoTypes.AttributeValueMemberS{Value: value}
	case bool:
		attributeValue = &dynamoTypes.AttributeValueMemberBOOL{Value: value}
	case int:
		attributeValue = &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(int64(value), 10)}
	case int32:
		attributeValue = &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(int64(value), 10)}
	case int64:
		attributeValue = &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(value, 10)}
	case uint:
		attributeValue = &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatUint(uint64(value), 10)}
	case uint32:
		attributeValue = &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatUint(uint64(value), 10)}
	case uint64:
		attributeValue = &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatUint(value, 10)}
	case float32:
		attributeValue = &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatFloat(float64(value), 'f', -1, 32)}
	case float64:
		attributeValue = &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatFloat(value, 'f', -1, 64)}
	case []dynamoTypes.AttributeValue:
		attributeValue = &dynamoTypes.AttributeValueMemberL{Value: value}
	case map[string]dynamoTypes.AttributeValue:
		attributeValue = &dynamoTypes.AttributeValueMemberM{Value: value}
	case []string:
		attributeValue = &dynamoTypes.AttributeValueMemberSS{Value: value}
	case []int32:
		values := make([]dynamoTypes.AttributeValue, len(value))
		for i, number := range value {
			values[i] = &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(int64(number), 10)}
		}
		attributeValue = &dynamoTypes.AttributeValueMemberL{
			Value: values,
		}
	case []int64:
		values := make([]dynamoTypes.AttributeValue, len(value))
		for i, number := range value {
			values[i] = &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(number, 10)}
		}
		attributeValue = &dynamoTypes.AttributeValueMemberL{
			Value: values,
		}
	case []uint32:
		values := make([]dynamoTypes.AttributeValue, len(value))
		for i, number := range value {
			values[i] = &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatUint(uint64(number), 10)}
		}
		attributeValue = &dynamoTypes.AttributeValueMemberL{
			Value: values,
		}
	case []uint64:
		values := make([]dynamoTypes.AttributeValue, len(value))
		for i, number := range value {
			values[i] = &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatUint(number, 10)}
		}
		attributeValue = &dynamoTypes.AttributeValueMemberL{
			Value: values,
		}
	case []float32:
		values := make([]dynamoTypes.AttributeValue, len(value))
		for i, number := range value {
			values[i] = &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatFloat(float64(number), 'f', -1, 32)}
		}
		attributeValue = &dynamoTypes.AttributeValueMemberL{
			Value: values,
		}
	case []float64:
		values := make([]dynamoTypes.AttributeValue, len(value))
		for i, number := range value {
			values[i] = &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatFloat(number, 'f', -1, 64)}
		}
		attributeValue = &dynamoTypes.AttributeValueMemberL{
			Value: values,
		}
	}

	builder.AddAttributeValue(key, attributeValue)
}

func (builder *UpdateBuilder) AddAttributeName(key string, value string) {
	if builder.attributeNames == nil {
		builder.attributeNames = map[string]string{}
	}

	if strings.HasPrefix(key, "#") {
		builder.attributeNames[key] = value
	} else {
		builder.attributeNames["#"+key] = value
	}
}

func (builder *UpdateBuilder) GetAttributeNames() map[string]string {
	if builder.attributeNames == nil {
		builder.attributeNames = map[string]string{}
	}

	return builder.attributeNames
}

func (builder *UpdateBuilder) AddAttributeValue(key string, value dynamoTypes.AttributeValue) {
	if builder.attributeValues == nil {
		builder.attributeValues = map[string]dynamoTypes.AttributeValue{}
	}

	if strings.HasPrefix(key, ":") {
		builder.attributeValues[key] = value
	} else {
		builder.attributeValues[":"+key] = value
	}
}

func (builder *UpdateBuilder) GetAttributeValues() map[string]dynamoTypes.AttributeValue {
	if builder.attributeValues == nil {
		builder.attributeValues = map[string]dynamoTypes.AttributeValue{}
	}

	return builder.attributeValues
}
