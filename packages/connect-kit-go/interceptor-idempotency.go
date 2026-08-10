package connectkit

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"time"

	"connectrpc.com/connect"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

type IdempotencyEntry struct {
	Pk       string
	Sk       string
	Ttl      int64
	Response *[]byte
}

func NewIdempotencyInterceptor(dynamoRead *dynamodb.Client, dynamoWrite *dynamodb.Client, tableName *string) connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			logger := GetLogger(ctx)
			authCtx := GetAuthContext(ctx)

			var key string
			{
				k := req.Header().Get("idempotency-key")
				if k == "" {
					return next(ctx, req)
				}

				var identity string
				{
					if authCtx.Iam != nil {
						identity = authCtx.Iam.UserArn
					} else if authCtx.Lambda != nil {
						identity = authCtx.Lambda.UserId
					} else {
						return next(ctx, req)
					}
				}

				key = "IDEMPOTENCY#" + req.Spec().Procedure + "#" + identity + "#" + k
			}

			_, err := dynamoWrite.PutItem(ctx, &dynamodb.PutItemInput{
				TableName: tableName,
				Item: map[string]dynamoTypes.AttributeValue{
					"pk":  &dynamoTypes.AttributeValueMemberS{Value: key},
					"sk":  &dynamoTypes.AttributeValueMemberS{Value: "IDEMPOTENCY"},
					"ttl": &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(time.Now().Unix()+20, 10)},
				},
				ConditionExpression: new("attribute_not_exists(pk) OR (attribute_exists(#tol) AND #tol <= :now)"),
				ExpressionAttributeNames: map[string]string{
					"#tol": "ttl",
				},
				ExpressionAttributeValues: map[string]dynamoTypes.AttributeValue{
					":now": &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(time.Now().Unix(), 10)},
				},
			})
			if err != nil && !errors.As(err, new(*dynamoTypes.ConditionalCheckFailedException)) {
				logger.Error("Error creating idempotency lock", slog.Any("error", err))
				return nil, NewUnexpected()
			}

			if err != nil && errors.As(err, new(*dynamoTypes.ConditionalCheckFailedException)) {
				var entry IdempotencyEntry
				var found bool
				for i := range 4 {
					response, err := dynamoRead.Query(ctx, &dynamodb.QueryInput{
						TableName:              tableName,
						KeyConditionExpression: new("pk = :pk"),
						ExpressionAttributeValues: map[string]dynamoTypes.AttributeValue{
							":pk": &dynamoTypes.AttributeValueMemberS{Value: key},
						}})
					if err != nil {
						logger.Error("Error retrieving valid idempotency entry", slog.Any("error", err))
						return nil, NewUnexpected()
					}
					if len(response.Items) > 0 {
						err = attributevalue.UnmarshalMap(response.Items[0], &entry)
						if err != nil {
							logger.Error("Error parsing idempotency entry", slog.Any("item", response.Items[0]), slog.Any("error", err))
							return nil, NewUnexpected()
						}
						if entry.Response != nil && entry.Ttl > time.Now().Unix() {
							found = true
							break
						}
					}

					time.Sleep(time.Second * (1 << i))
				}

				if !found {
					logger.Error("Unable to establish idempotency state")
					return nil, NewUnexpected()
				}

				md, ok := req.Spec().Schema.(protoreflect.MethodDescriptor)
				if !ok {
					logger.Error("spec.Schema is not a MethodDescriptor")
					return nil, NewUnexpected()
				}

				msg := dynamicpb.NewMessage(md.Output())

				if err := proto.Unmarshal(*entry.Response, msg); err != nil {
					logger.Error("Error parsing response", slog.Any("entry", entry), slog.Any("error", err))
					return nil, NewUnexpected()
				}

				return connect.NewResponse(msg), nil
			}

			resp, err := next(ctx, req)
			if err != nil {
				return resp, err
			}

			var respBytes []byte
			if resp != nil && resp.Any() != nil {
				if respMsg, ok := resp.Any().(proto.Message); ok {
					respBytes, _ = proto.Marshal(respMsg)
				}
			}

			_, err = dynamoWrite.PutItem(ctx, &dynamodb.PutItemInput{
				TableName: tableName,
				Item: map[string]dynamoTypes.AttributeValue{
					"pk":       &dynamoTypes.AttributeValueMemberS{Value: key},
					"sk":       &dynamoTypes.AttributeValueMemberS{Value: "IDEMPOTENCY"},
					"ttl":      &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(time.Now().Unix()+172800, 10)},
					"response": &dynamoTypes.AttributeValueMemberB{Value: respBytes},
				},
			})
			if err != nil && !errors.As(err, new(*dynamoTypes.ConditionalCheckFailedException)) {
				logger.Warn("Error creating idempotency entry", slog.Any("error", err))
				return nil, NewUnexpected()
			}

			return resp, err
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
