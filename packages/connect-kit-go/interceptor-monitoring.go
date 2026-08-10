package connectkit

import (
	"context"
	"log/slog"
	"os"

	"connectrpc.com/connect"
	"github.com/aws/aws-xray-sdk-go/xray"
)

type contextKey int

const monitoringInterceptorContextKey contextKey = 1

type NewMonitoringInterceptorInput struct {
	Logger   *slog.Logger
	Metadata map[string]any
}

func NewMonitoringInterceptor(params NewMonitoringInterceptorInput) connect.UnaryInterceptorFunc {

	metadata := map[string]any{}
	if params.Metadata != nil {
		metadata = params.Metadata
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if params.Logger != nil {
		logger = params.Logger
	}

	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			ctx = context.WithValue(ctx, "x-amzn-trace-id", req.Header().Get("x-amzn-trace-id"))
			ctx, segment := xray.BeginSubsegment(ctx, req.Spec().Procedure)
			if segment != nil {
				defer segment.Close(nil)
				segment.AddAnnotation("Procedure", req.Spec().Procedure)
				for key, value := range metadata {
					segment.AddAnnotation(key, value)
				}
			}

			ctx = context.WithValue(ctx, monitoringInterceptorContextKey, logger)

			logger.Info("Request",
				slog.String("procedure", req.Spec().Procedure),
				slog.String("traceId", xray.TraceID(ctx)),
				slog.Any("metadata", metadata))

			rep, err := next(ctx, req)

			if err != nil {
				logger.Error("Response",
					slog.String("status", "ERROR"),
					slog.String("procedure", req.Spec().Procedure),
					slog.Any("error", err))
				if segment != nil {
					segment.AddError(err)
				} else {
					xray.AddError(ctx, err)
				}
			} else {
				logger.Info("Response", slog.String("status", "SUCCESS"))
			}

			return rep, err
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}

func GetLogger(ctx context.Context) *slog.Logger {
	if v, ok := ctx.Value(monitoringInterceptorContextKey).(*slog.Logger); ok {
		return v
	}
	return slog.Default()
}
