package main

import (
	"context"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-xray-sdk-go/xray"
)

func handleRecord(ctx context.Context, record events.DynamoDBEventRecord) error {
	Logger.Info("handleRecord", slog.Any("record", record))

	var recordType string
	switch {
	case record.Change.NewImage != nil:
		if attr, ok := record.Change.NewImage["type"]; ok && attr.DataType() == events.DataTypeString {
			recordType = attr.String()
		}
	case record.Change.OldImage != nil:
		if attr, ok := record.Change.OldImage["type"]; ok && attr.DataType() == events.DataTypeString {
			recordType = attr.String()
		}
	}

	switch recordType {
	case "USER":
		return handlerUser(ctx, record)
	}

	return nil
}

func handleRequest(ctx context.Context, event events.DynamoDBEvent) (events.DynamoDBEventResponse, error) {
	ctx, segment := xray.BeginSubsegment(ctx, "Handler")
	segment.AddAnnotation("Stack", StackName)
	defer segment.Close(nil)

	Logger.Info("handleRequest", slog.Any("event", event))

	resp := events.DynamoDBEventResponse{}
	for _, record := range event.Records {
		if err := handleRecord(ctx, record); err != nil {
			Logger.Error("Failed to proccess record", slog.Any("error", err), slog.Any("record", record))
			resp.BatchItemFailures = append(
				resp.BatchItemFailures,
				events.DynamoDBBatchItemFailure{ItemIdentifier: record.Change.SequenceNumber},
			)
		}
	}

	return resp, nil
}

func main() {
	lambda.Start(handleRequest)
}
