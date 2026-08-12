package main

import (
	"connectkit"
	"context"
	"contracts/dist/catalogue/v1"
	"log/slog"

	"connectrpc.com/connect"
)

func (s *Handler) ListUnits(
	ctx context.Context,
	req *connect.Request[catalogue.ListUnitsRequest],
) (*connect.Response[catalogue.ListUnitsResponse], error) {
	logger := connectkit.GetLogger(ctx)
	authCtx := connectkit.GetAuthContext(ctx)

	if authCtx.Lambda != nil {
		resp, err := authCtx.Evaluate(ctx, "catalogue", "list_units",
			&catalogue.ListUnitsContext{
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

	return connect.NewResponse(&catalogue.ListUnitsResponse{Units: []*catalogue.Unit{
		&catalogue.Unit{},
	}}), nil
}
