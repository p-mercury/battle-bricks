package main

import (
	"connectkit"
	"context"
	"encoding/json"
	"log/slog"

	"contracts/dist/identity/v1"

	"connectrpc.com/connect"
)

type JWKS struct {
	Keys []identity.Jwk
}

func (s *Handler) ListJwks(
	ctx context.Context,
	req *connect.Request[identity.ListJwksRequest],
) (*connect.Response[identity.ListJwksResponse], error) {
	logger := ctx.Value("logger").(*slog.Logger)

	rep, err := Http.Get(*UserPoolProviderUrl + "/.well-known/jwks.json")
	if err != nil {
		logger.Error("Error retriving jwks", slog.Any("error", err))
		return nil, connectkit.NewUnexpected()
	}
	defer rep.Body.Close()

	var jwks JWKS
	if err := json.NewDecoder(rep.Body).Decode(&jwks); err != nil {
		logger.Error("Error decoding JWKS", slog.Any("error", err))
		return nil, connectkit.NewUnexpected()
	}

	resp := &identity.ListJwksResponse{
		Jwks: make([]*identity.Jwk, len(jwks.Keys)),
	}

	for i := range jwks.Keys {
		resp.Jwks[i] = &jwks.Keys[i]
	}

	return connect.NewResponse(resp), nil
}
