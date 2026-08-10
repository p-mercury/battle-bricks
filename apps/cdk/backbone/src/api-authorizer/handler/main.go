package main

import (
	"context"
	"contracts/dist/identity/v1"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/golang-jwt/jwt/v5"
)

var (
	jwksCache = struct {
		keys map[string]*rsa.PublicKey
		time time.Time
	}{keys: make(map[string]*rsa.PublicKey)}
)

func fetchJwks(ctx context.Context,
) error {
	jwksCache.time = time.Now()

	var (
		mu      sync.Mutex
		newKeys = make(map[string]*rsa.PublicKey)
	)

	resp, err := IdentityService.ListJwks(ctx, connect.NewRequest(&identity.ListJwksRequest{}))
	if err != nil {
		return nil
	}

	for _, jwk := range resp.Msg.Jwks {
		if pubKey, err := jwkToRsa(jwk); err == nil {
			mu.Lock()
			newKeys[jwk.Kid] = pubKey
			mu.Unlock()
		}
	}

	jwksCache.keys = newKeys
	jwksCache.time = time.Now()
	return nil
}

type Jwk interface {
	GetN() string
	GetE() string
}

func jwkToRsa(jwk Jwk) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.GetN())
	if err != nil {
		return nil, err
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.GetE())
	if err != nil {
		return nil, err
	}

	eInt := 0
	for _, b := range eBytes {
		eInt = (eInt << 8) | int(b)
	}

	pubKey := &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: eInt,
	}
	return pubKey, nil
}

func extractAccessToken(cookies []string) (string, error) {
	for _, c := range cookies {
		k, v, ok := strings.Cut(c, "=")
		if ok && k == "accessToken" {
			return v, nil
		}
	}
	return "", errors.New("accessToken cookie not found")
}

func verifyJwt(ctx context.Context, tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("Unexpected signing method")
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("Missing kid in token header")
		}

		if time.Since(jwksCache.time) > 10*time.Minute {
			fetchJwks(ctx)
		}

		pubKey, exists := jwksCache.keys[kid]

		if !exists {
			if time.Since(jwksCache.time) > time.Minute {
				if err := fetchJwks(ctx); err == nil {
					pubKey, exists = jwksCache.keys[kid]
				}
			} else {
				return nil, errors.New("No kid match in JWKS")
			}
		}

		if !exists {
			return nil, errors.New("No kid match in JWKS")
		}

		return pubKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}

func handleRequest(ctx context.Context, event events.APIGatewayV2CustomAuthorizerV2Request) (events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {
	ctx, segment := xray.BeginSubsegment(ctx, "Handler")
	segment.AddAnnotation("Stack", StackName)
	defer segment.Close(nil)

	token, err := extractAccessToken(event.Cookies)
	if err != nil {
		return events.APIGatewayV2CustomAuthorizerSimpleResponse{}, errors.New("Unauthorized")
	}

	claims, err := verifyJwt(ctx, token)
	if err != nil {
		return events.APIGatewayV2CustomAuthorizerSimpleResponse{}, errors.New("Unauthorized")
	}

	return events.APIGatewayV2CustomAuthorizerSimpleResponse{IsAuthorized: true, Context: claims}, nil
}

func main() {
	lambda.Start(handleRequest)
}
