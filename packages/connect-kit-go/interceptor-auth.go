package connectkit

import (
	"context"
	"contracts/dist/policy/v1"
	"contracts/dist/policy/v1/policyconnect"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const authInterceptorContextKey contextKey = 0

type Iam struct {
	AccessKey string `json:"accessKey"`
	AccountId string `json:"accountId"`
	CallerId  string `json:"callerId"`
	UserArn   string `json:"userArn"`
	UserId    string `json:"userId"`
}

type Lambda struct {
	Sub            string   `json:"sub"`
	UserId         string   `json:"userId"`
	OrganisationId *string  `json:"organisationId"`
	Scope          string   `json:"scope"`
	Aud            []string `json:"aud"`
}

type AuthInterceptorContext struct {
	Iam              *Iam
	Lambda           *Lambda
	EvaluateResponse *policy.EvaluateResponse
	policyService    policyconnect.InternalClient
}

func parseAuthContext(
	ctx context.Context,
	header http.Header,
	policyService policyconnect.InternalClient,
) (*AuthInterceptorContext, error) {
	logger := GetLogger(ctx)

	var requestContext struct {
		Authorizer *struct {
			Iam    *Iam    `json:"iam"`
			Lambda *Lambda `json:"lambda"`
		}
	}

	r := header.Get("x-amzn-request-context")
	if r == "" {
		logger.Error("Missing x-amzn-request-context header")
		return nil, NewUnexpected()
	}
	if err := json.Unmarshal([]byte(r), &requestContext); err != nil {
		logger.Error("Error parsing x-amzn-request-context header", slog.Any("error", err))
		return nil, NewUnexpected()
	}
	if requestContext.Authorizer == nil {
		logger.Error("Missing authorizer in context")
		return nil, NewUnexpected()
	}

	switch {
	case requestContext.Authorizer.Iam != nil:
		return &AuthInterceptorContext{
			Iam:           requestContext.Authorizer.Iam,
			policyService: policyService,
		}, nil
	case requestContext.Authorizer.Lambda != nil:
		return &AuthInterceptorContext{
			Lambda:        requestContext.Authorizer.Lambda,
			policyService: policyService,
		}, nil
	default:
		logger.Error("Invalid authorizer context")
		return nil, NewUnexpected()
	}
}

func (authCtx *AuthInterceptorContext) Evaluate(ctx context.Context, service string, method string, c proto.Message) (*policy.EvaluateResponse, error) {
	logger := GetLogger(ctx)

	if authCtx.Lambda == nil {
		return nil, fmt.Errorf("Can't evaluate policy for non lambda authorizer")
	}

	cJson, err := protojson.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling context: %v", err)
	}

	resp, err := authCtx.policyService.Evaluate(ctx, connect.NewRequest(
		&policy.EvaluateRequest{
			Service: service,
			Method:  method,
			UserId:  authCtx.Lambda.UserId,
			Context: cJson,
		}))
	if err != nil {
		return nil, fmt.Errorf("Error evaluating policy: %v", err)
	}

	authCtx.EvaluateResponse = resp.Msg

	logger.Info("EvaluatePolicy", slog.Any("decisionId", resp.Msg.Id), slog.Any("authz", resp.Msg.Authz))

	return resp.Msg, nil
}

func NewAuthInterceptor(policyService policyconnect.InternalClient) connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			authCtx, err := parseAuthContext(ctx, req.Header(), policyService)
			if err != nil {
				return nil, err
			}
			ctx = context.WithValue(ctx, authInterceptorContextKey, authCtx)

			if authCtx.Iam != nil {
				return next(ctx, req)
			}

			res, err := next(ctx, req)
			if err != nil {
				return nil, err
			}
			if res == nil {
				return res, nil
			}
			if authCtx.EvaluateResponse != nil {
				if !authCtx.EvaluateResponse.Authz {
					return nil, NewUnauthorized()
				} else if pm, ok := res.Any().(protoreflect.ProtoMessage); ok && pm != nil {
					ApplyOutputMask(pm, authCtx.EvaluateResponse.OutputMask)
				}
			}
			return res, nil
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}

func GetAuthContext(ctx context.Context) *AuthInterceptorContext {
	if v, ok := ctx.Value(authInterceptorContextKey).(*AuthInterceptorContext); ok {
		return v
	}
	return nil
}
