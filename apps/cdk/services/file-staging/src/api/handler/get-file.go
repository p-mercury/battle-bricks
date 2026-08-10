package main

import (
	"connectkit"
	"context"
	"contracts/dist/filestaging/v1"
	"log/slog"
	"time"

	"file-staging/src/types/dynamo"

	"connectrpc.com/connect"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Handler) GetFile(
	ctx context.Context,
	req *connect.Request[filestaging.GetFileRequest],
) (*connect.Response[filestaging.GetFileResponse], error) {
	logger := connectkit.GetLogger(ctx)
	authCtx := connectkit.GetAuthContext(ctx)

	var file *filestaging.File
	{
		response, err := DynamoRead.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: TableName,
			Key: map[string]dynamoTypes.AttributeValue{
				"pk": &dynamoTypes.AttributeValueMemberS{Value: req.Msg.Id},
			},
		})
		if err != nil {
			logger.Error("Error getting file", slog.Any("error", err))
			return nil, connectkit.NewUnexpected()
		}
		if len(response.Item) > 0 {
			var r dynamo.File
			if err := attributevalue.UnmarshalMap(response.Item, &r); err != nil {
				logger.Error("Error parsing file", slog.Any("error", err))
				return nil, connectkit.NewUnexpected()
			}

			file = &filestaging.File{
				Id:          r.Pk,
				Status:      filestaging.FileStatus(r.Status),
				Name:        r.Name,
				Checksum:    r.Checksum,
				ExpiresTime: timestamppb.New(time.Unix(r.Ttl, 0)),
				UserId:      r.UserId,
			}

			if file.Status == filestaging.FileStatus_FILE_STATUS_VALID {
				presigned, err := S3Presign.PresignGetObject(ctx,
					&s3.GetObjectInput{
						Bucket:                     new(r.Bucket),
						Key:                        new(r.Key),
						ResponseContentDisposition: new("attachment; filename=\"" + r.Name + "\""),
					},
					s3.WithPresignExpires(time.Minute*2))
				if err != nil {
					logger.Error("Error generating presigned url", slog.Any("error", err))
					return nil, connectkit.NewUnexpected()
				}
				file.Url = &presigned.URL
				file.ContentType = r.ContentType
				file.ContentLength = r.ContentLength
			}
		}
	}

	if authCtx.Lambda != nil {
		resp, err := authCtx.Evaluate(ctx, "file_staging", "get_file",
			&filestaging.GetFileContext{
				Request: req.Msg,
				Subject: file,
			},
		)
		if err != nil {
			logger.Error("Error evaluating policy", slog.Any("error", err))
			return nil, connectkit.NewUnexpected()
		} else if !resp.Authz {
			return nil, connectkit.NewUnauthorized()
		}
	}

	if file == nil {
		return nil, connectkit.NewNotFound()
	}

	return connect.NewResponse(&filestaging.GetFileResponse{File: file}), nil
}
