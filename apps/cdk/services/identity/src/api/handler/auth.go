package main

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

func NewAuthHandler(ctx context.Context, oauth2Url string, issuerUrl string, clientId string, clientSecret string) (*OAuthHandler, error) {
	handler := &OAuthHandler{
		clientId:     clientId,
		clientSecret: clientSecret,
		oauth2Url:    oauth2Url,
	}

	{
		provider, err := oidc.NewProvider(ctx, issuerUrl)
		if err != nil {
			Logger.Error("Error creating provider", slog.Any("error", err))
			return nil, err
		}
		handler.authenticator = oauth2.Config{
			ClientID:     handler.clientId,
			ClientSecret: handler.clientSecret,
			Endpoint:     provider.Endpoint(),
			Scopes:       []string{oidc.ScopeOpenID},
		}
	}

	return handler, nil
}

type OAuthHandler struct {
	clientId      string
	clientSecret  string
	authenticator oauth2.Config
	oauth2Url     string
}

func (handler OAuthHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx, segment := xray.BeginSubsegment(req.Context(), "/auth_service.Service")
	if segment != nil {
		segment.AddAnnotation("Procedure", "/identity.v1.service/"+req.Method+req.URL.Path)
		defer segment.Close(nil)
	}

	logger := Logger.With(slog.String("procedure", "/identity.v1.service/"+req.Method+req.URL.Path))
	logger.Info("Request")

	switch path.Base(req.URL.Path) {

	case "authorize":
		authorizeUrl, err := url.Parse(*UserPoolUrl + "/oauth2/authorize")
		if err != nil {
			logger.Error("Error creating authorize url", slog.Any("error", err))
			if segment != nil {
				segment.Error = true
			}
			http.Error(w, "Internal error", 500)
			return
		}

		query := authorizeUrl.Query()
		query.Set("client_id", handler.clientId)
		query.Set("redirect_uri", req.URL.Query().Get("redirect_uri"))
		query.Set("response_type", req.URL.Query().Get("response_type"))
		query.Set("scope", req.URL.Query().Get("scope"))
		query.Set("state", req.URL.Query().Get("state"))
		authorizeUrl.RawQuery = query.Encode()

		w.Header().Add("set-cookie", "accessToken=;HttpOnly;Secure;SameSite=Strict;Path=/;Max-Age=0")
		w.Header().Add("set-cookie", "refreshToken=;HttpOnly;Secure;SameSite=Strict;Path=/;Max-Age=0")
		http.Redirect(w, req, authorizeUrl.String(), 302)
		return

	case "code":
		code := req.URL.Query().Get("code")
		if code == "" {
			logger.Warn("Missing code search parameters")
			if segment != nil {
				segment.Error = true
			}
			http.Error(w, "Missing code search parameters", 500)
			return
		}

		redirectUri := req.URL.Query().Get("redirectUri")
		if redirectUri == "" {
			logger.Warn("Missing redirectUri search parameters")
			if segment != nil {
				segment.Error = true
			}
			http.Error(w, "Missing redirectUri search parameters", 500)
			return
		}

		tokens, err := handler.authenticator.Exchange(ctx, code, oauth2.SetAuthURLParam("redirect_uri", redirectUri))
		if err != nil {
			logger.Warn("Failed to exchange an authorization code for a token", slog.Any("error", err))
			if segment != nil {
				segment.Error = true
			}
			http.Error(w, "Internal error", 500)
			return
		}

		w.Header().Add("set-cookie", "accessToken="+tokens.AccessToken+";HttpOnly;Secure;SameSite=Strict;Path=/")
		w.Header().Add("set-cookie", "refreshToken="+tokens.RefreshToken+";HttpOnly;Secure;SameSite=Strict;Path=/;Max-Age=604800")
		w.WriteHeader(http.StatusOK)
		return

	case "refresh":
		var refreshToken = ""
		{
			refreshTokenCookie, err := req.Cookie("refreshToken")
			if err != nil {
				logger.Warn("Missing refreshToken cookie", slog.Any("error", err))
				if segment != nil {
					segment.Error = true
				}
				http.Error(w, "Missing refreshToken cookie", 500)
				return
			}
			refreshToken = refreshTokenCookie.Value

			if refreshToken == "" {
				logger.Warn("Missing refreshToken")
				if segment != nil {
					segment.Error = true
				}
				http.Error(w, "Missing refreshToken", 500)
				return
			}
		}

		payload := strings.NewReader("grant_type=refresh_token&client_id=" + handler.clientId + "&client_secret=" + handler.clientSecret + "&refresh_token=" + refreshToken)

		request, err := http.NewRequest("POST", handler.oauth2Url+"/oauth2/token", payload)
		if err != nil {
			logger.Warn("Response error", slog.Any("error", err))
			if segment != nil {
				segment.Error = true
			}
			http.Error(w, "Internal error", 500)
			return
		}

		request.Header.Add("content-type", "application/x-www-form-urlencoded")

		response, err := Http.Do(request.WithContext(ctx))
		if err != nil {
			logger.Warn("Response error", slog.Any("error", err))
			if segment != nil {
				segment.Error = true
			}
			http.Error(w, "Internal error", 500)
			return
		}

		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			logger.Warn("Response error", slog.Any("error", err))
			if segment != nil {
				segment.Error = true
			}
			http.Error(w, "Internal error", 500)
			return
		}

		type RefreshResponse struct {
			AccessToken  *string `json:"access_token"`
			RefreshToken *string `json:"refresh_token"`
			ExpiresIn    *int32  `json:"expires_in"`
		}

		var refreshResponse RefreshResponse
		err = json.Unmarshal(body, &refreshResponse)
		if err != nil {
			logger.Warn("Error parsin reponse", slog.Any("error", err))
			if segment != nil {
				segment.Error = true
			}
			http.Error(w, "Internal error", 500)
			return
		}

		if refreshResponse.AccessToken == nil {
			logger.Warn("Missing access token in refresh response", slog.Any("error", err))
			if segment != nil {
				segment.Error = true
			}
			http.Error(w, "Internal error", 500)
			return
		}

		if refreshResponse.RefreshToken != nil {
			w.Header().Add("set-cookie", "refreshToken="+*refreshResponse.RefreshToken+";HttpOnly;Secure;SameSite=Strict;Path=/;Max-Age=604800")
		}

		w.Header().Add("set-cookie", "accessToken="+*refreshResponse.AccessToken+";HttpOnly;Secure;SameSite=Strict;Path=/")
		w.WriteHeader(http.StatusOK)
		return

	case "signout":
		logoutUrl, err := url.Parse(handler.oauth2Url + "/logout")
		if err != nil {
			logger.Error("Error creating logout url", slog.Any("error", err))
			if segment != nil {
				segment.Error = true
			}
			http.Error(w, "Internal error", 500)
			return
		}

		query := logoutUrl.Query()
		query.Set("client_id", handler.clientId)
		query.Set("logout_uri", req.URL.Query().Get("redirect"))
		logoutUrl.RawQuery = query.Encode()

		w.Header().Add("set-cookie", "accessToken=;HttpOnly;Secure;SameSite=Strict;Path=/;Max-Age=0")
		w.Header().Add("set-cookie", "refreshToken=;HttpOnly;Secure;SameSite=Strict;Path=/;Max-Age=0")
		http.Redirect(w, req, logoutUrl.String(), 302)
		return

	default:
		logger.Warn("Unknow endpoint")
		if segment != nil {
			segment.Error = true
		}
		http.Error(w, "Unknow endpoint", 404)
	}
}
