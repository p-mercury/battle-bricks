package main

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-xray-sdk-go/xray"
)

func handleMessage(ctx context.Context, message events.SQSMessage) (err error) {
	var segment *xray.Segment
	ctx, segment = xray.BeginSubsegment(ctx, "HandleMessage")
	segment.AddAnnotation("MessageId", message.MessageId)
	defer func() { segment.Close(err) }()

	Logger.Info("HandleMessage", slog.Any("message", message))

	var event events.EventBridgeEvent
	if err := json.Unmarshal([]byte(message.Body), &event); err != nil {
		return err
	}

	switch event.Source {
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
