package main

import (
	"connectkit"
	"context"
	"contracts/dist/filestaging/v1"
	"log/slog"
	"strconv"
	"time"

	"connectrpc.com/connect"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Handler) CreateUploadSession(
	ctx context.Context,
	req *connect.Request[filestaging.CreateUploadSessionRequest],
) (*connect.Response[filestaging.CreateUploadSessionResponse], error) {
	logger := connectkit.GetLogger(ctx)
	authCtx := connectkit.GetAuthContext(ctx)
	timestamp := time.Now().Truncate(time.Millisecond)

	if authCtx.Lambda != nil {
		resp, err := authCtx.Evaluate(ctx, "file_staging", "create_upload_session",
			&filestaging.CreateUploadSessionContext{
				Request: req.Msg,
			},
		)
		if err != nil {
			logger.Error("Error evaluating policy", slog.Any("error", err))
			return nil, connectkit.NewUnexpected()
		} else if !resp.Authz {
			return nil, connectkit.NewUnauthorized()
		}
	}

	id, err := connectkit.NewBase62Id("fs-", 40)
	if err != nil {
		logger.Error("Error generating id", slog.Any("error", err))
		return nil, connectkit.NewUnexpected()
	}

	key := id + "/" + req.Msg.Name

	presigned, err := S3Presign.PresignPutObject(ctx,
		&s3.PutObjectInput{
			Bucket:         BucketName,
			Key:            &key,
			ACL:            s3Types.ObjectCannedACLBucketOwnerFullControl,
			ChecksumSHA256: &req.Msg.Checksum,
			IfNoneMatch:    new("*"),
		},
		s3.WithPresignExpires(time.Minute),
		func(opt *s3.PresignOptions) {
			opt.Presigner = v4.NewSigner(func(so *v4.SignerOptions) {
				so.DisableURIPathEscaping = true
				so.DisableHeaderHoisting = true
			})
		},
	)
	if err != nil {
		logger.Error("Error generating presigned post", slog.Any("error", err))
		return nil, connectkit.NewUnexpected()
	}

	headers := map[string]string{}
	for key, value := range presigned.SignedHeader {
		headers[key] = value[0]
	}

	if authCtx.Lambda != nil {
		_, err = DynamoWrite.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
			TransactItems: []dynamoTypes.TransactWriteItem{
				{
					Put: &dynamoTypes.Put{
						TableName: TableName,
						Item: map[string]dynamoTypes.AttributeValue{
							"pk":       &dynamoTypes.AttributeValueMemberS{Value: id},
							"type":     &dynamoTypes.AttributeValueMemberS{Value: "FILE"},
							"userId":   &dynamoTypes.AttributeValueMemberS{Value: authCtx.Lambda.UserId},
							"bucket":   &dynamoTypes.AttributeValueMemberS{Value: *BucketName},
							"key":      &dynamoTypes.AttributeValueMemberS{Value: key},
							"name":     &dynamoTypes.AttributeValueMemberS{Value: req.Msg.Name},
							"status":   &dynamoTypes.AttributeValueMemberN{Value: "1"},
							"checksum": &dynamoTypes.AttributeValueMemberS{Value: req.Msg.Checksum},
							"ttl":      &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(timestamp.Add(time.Hour*720).Unix(), 10)},
						},
						ConditionExpression: new("attribute_not_exists(pk)"),
					},
				},
				{
					Update: &dynamoTypes.Update{
						TableName: TableName,
						Key: map[string]dynamoTypes.AttributeValue{
							"pk": &dynamoTypes.AttributeValueMemberS{Value: authCtx.Lambda.UserId + "#" + strconv.FormatInt(timestamp.Truncate(time.Hour).Unix(), 10)},
						},
						UpdateExpression: new("SET #upl = if_not_exists(#upl, :start) + :inc, #tol = :tol, #typ = :typ"),
						ExpressionAttributeNames: map[string]string{
							"#upl": "uploads",
							"#tol": "ttl",
							"#typ": "type",
						},
						ExpressionAttributeValues: map[string]dynamoTypes.AttributeValue{
							":start": &dynamoTypes.AttributeValueMemberN{Value: "0"},
							":inc":   &dynamoTypes.AttributeValueMemberN{Value: "1"},
							":max":   &dynamoTypes.AttributeValueMemberN{Value: "40"},
							":tol":   &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(timestamp.Truncate(time.Hour).Add(time.Hour*2).Unix(), 10)},
							":typ":   &dynamoTypes.AttributeValueMemberS{Value: "UPLOAD_COUNTER"},
						},
						ConditionExpression: new("attribute_not_exists(#upl) OR #upl <= :max"),
					},
				},
			},
		})
		if err != nil {
			logger.Error("Error creating file in DynamoDB", slog.Any("error", err))
			return nil, connectkit.NewUnexpected()
		}
	} else {
		_, err = DynamoWrite.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: TableName,
			Item: map[string]dynamoTypes.AttributeValue{
				"pk":       &dynamoTypes.AttributeValueMemberS{Value: id},
				"bucket":   &dynamoTypes.AttributeValueMemberS{Value: *BucketName},
				"key":      &dynamoTypes.AttributeValueMemberS{Value: key},
				"name":     &dynamoTypes.AttributeValueMemberS{Value: req.Msg.Name},
				"checksum": &dynamoTypes.AttributeValueMemberS{Value: req.Msg.Checksum},
				"ttl":      &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(timestamp.Add(time.Hour*720).Unix(), 10)},
			},
			ConditionExpression: new("attribute_not_exists(pk)"),
		})
		if err != nil {
			logger.Error("Error creating file in DynamoDB", slog.Any("error", err))
			return nil, connectkit.NewUnexpected()
		}
	}

	return connect.NewResponse(&filestaging.CreateUploadSessionResponse{
		Session: &filestaging.UploadSession{
			Id:          id,
			Url:         presigned.URL,
			Headers:     headers,
			ExpiresTime: timestamppb.New(timestamp.Add(time.Hour * 720).Truncate(time.Second)),
		},
	}), nil
}
