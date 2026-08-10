package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-xray-sdk-go/xray"
)

type UpdateObjectInput struct {
	id            string
	status        string
	contentType   *string
	contentLength *int64
	checksum      *string
}

func updateObject(ctx context.Context, input UpdateObjectInput) error {
	updateExpression := "SET #sta = :sta"
	expressionAttributeNames := map[string]string{
		"#sta": "status",
	}
	expressionAttributeValues := map[string]dynamoTypes.AttributeValue{
		":sta": &dynamoTypes.AttributeValueMemberN{Value: input.status},
	}

	if input.contentType != nil {
		updateExpression += ", #cnt = :cnt"
		expressionAttributeNames["#cnt"] = "contentType"
		expressionAttributeValues[":cnt"] = &dynamoTypes.AttributeValueMemberS{Value: *input.contentType}
	}

	if input.contentLength != nil {
		updateExpression += ", #cnl = :cnl"
		expressionAttributeNames["#cnl"] = "contentLength"
		expressionAttributeValues[":cnl"] = &dynamoTypes.AttributeValueMemberN{Value: strconv.FormatInt(*input.contentLength, 10)}
	}

	if input.checksum != nil {
		updateExpression += ", #chs = :chs"
		expressionAttributeNames["#chs"] = "checksum"
		expressionAttributeValues[":chs"] = &dynamoTypes.AttributeValueMemberS{Value: *input.checksum}
	}

	_, err := Dynamo.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: TableName,
		Key: map[string]dynamoTypes.AttributeValue{
			"pk": &dynamoTypes.AttributeValueMemberS{Value: input.id},
		},
		UpdateExpression:          &updateExpression,
		ExpressionAttributeNames:  expressionAttributeNames,
		ExpressionAttributeValues: expressionAttributeValues,
	})
	if err != nil {
		Logger.Error("Error updating status", slog.Any("error", err))
	}
	return err
}

func handleMessage(ctx context.Context, message events.SQSMessage) (err error) {
	var segment *xray.Segment
	ctx, segment = xray.BeginSubsegment(ctx, "HandleMessage")
	segment.AddAnnotation("MessageId", message.MessageId)
	defer func() { segment.Close(err) }()

	Logger.Info("HandleMessage", slog.Any("message", message))

	var event events.S3Event
	if err := json.Unmarshal([]byte(message.Body), &event); err != nil {
		return err
	}

	Logger.Info("handleRecord", slog.Any("event", event))

	for _, record := range event.Records {
		var id string
		{
			l := strings.SplitN(record.S3.Object.Key, "/", 3)
			if len(l) == 3 {
				return nil
			} else if len(l) < 2 {
				return errors.New("invalid key length")
			}
			id = l[0]
		}

		tags, err := S3.GetObjectTagging(ctx, &s3.GetObjectTaggingInput{
			Bucket: BucketName,
			Key:    &record.S3.Object.Key,
		})
		if err != nil {
			return err
		}

		var status *string
		for _, t := range tags.TagSet {
			if aws.ToString(t.Key) == "GuardDutyMalwareScanStatus" {
				status = t.Value
				break
			}
		}

		if status == nil {
			updateObject(ctx, UpdateObjectInput{
				id:     id,
				status: "1",
			})
		} else if *status == "NO_THREATS_FOUND" {
			head, err := S3.HeadObject(ctx, &s3.HeadObjectInput{
				Bucket:       BucketName,
				Key:          &record.S3.Object.Key,
				ChecksumMode: s3Types.ChecksumModeEnabled,
			})
			if err != nil {
				return err
			}

			if head.ContentType == nil || head.ContentLength == nil || head.ChecksumSHA256 == nil {
				updateObject(ctx, UpdateObjectInput{
					id:     id,
					status: "4",
				})
				return nil
			}

			updateObject(ctx, UpdateObjectInput{
				id:            id,
				status:        "3",
				contentType:   head.ContentType,
				contentLength: head.ContentLength,
			})
		} else {
			updateObject(ctx, UpdateObjectInput{
				id:     id,
				status: "4",
			})
		}
	}

	return nil
}

func handleRequest(ctx context.Context, event events.SQSEvent) (_ events.SQSEventResponse, err error) {
	var segment *xray.Segment
	if xray.GetSegment(ctx) != nil {
		ctx, segment = xray.BeginSubsegment(ctx, "HandleRequest")
	} else {
		ctx, segment = xray.BeginSegment(ctx, "HandleRequest")
	}
	segment.AddAnnotation("Stack", StackName)
	defer func() { segment.Close(err) }()

	Logger.Info("handleRequest", slog.Any("event", event))

	resp := events.SQSEventResponse{BatchItemFailures: []events.SQSBatchItemFailure{}}
	for _, message := range event.Records {
		if err := handleMessage(ctx, message); err != nil {
			Logger.Error("Failed to proccess message", slog.Any("message", message), slog.Any("error", err))
			resp.BatchItemFailures = append(
				resp.BatchItemFailures,
				events.SQSBatchItemFailure{ItemIdentifier: message.MessageId},
			)
		}
	}

	return resp, nil
}

func main() {
	lambda.Start(handleRequest)
}
