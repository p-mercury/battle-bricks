package main

import (
	"connectkit"
	"context"
	"contracts/dist/catalogue/v1"
	"log/slog"

	"connectrpc.com/connect"
)

func (s *Handler) ListItems(
	ctx context.Context,
	req *connect.Request[catalogue.ListItemsRequest],
) (*connect.Response[catalogue.ListItemsResponse], error) {
	logger := connectkit.GetLogger(ctx)
	authCtx := connectkit.GetAuthContext(ctx)

	if authCtx.Lambda != nil {
		resp, err := authCtx.Evaluate(ctx, "catalogue", "list_items",
			&catalogue.ListItemsContext{
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

	return connect.NewResponse(&catalogue.ListItemsResponse{Items: []*catalogue.Item{
		&catalogue.Item{},
	}}), nil
}
