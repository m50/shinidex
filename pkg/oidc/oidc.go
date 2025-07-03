package oidc

import (
	"context"
	"net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

// config options:
const (
	KeyName             = "auth.oidc.name"
	KeyCreateUsers      = "auth.oidc.create-users"
	KeyDisablePassword  = "auth.oidc.disable-password"
	KeyClientID         = "auth.oidc.client-id"
	KeyClientSecret     = "auth.oidc.client-secret"
	KeyClientSecretFile = "auth.oidc.client-secret-file"
	KeyScopes           = "auth.oidc.scopes"
	KeyAuthURL          = "auth.oidc.auth-url"
	KeyTokenURL         = "auth.oidc.token-url"
	KeyUserInfoURL      = "auth.oidc.user-info-url"
	KeyCertificatesURL  = "auth.oidc.certificates-url"
	KeyBaseURL          = "auth.oidc.base-url"

	PathNameOIDCLogin    = "auth.oidc.login"
	PathNameOIDCCallback = "auth.oidc.redirect"
)

var (
	Provider     *oidc.Provider
	OAuth2Config oauth2.Config
)

func Initialize(ctx context.Context) error {
	if viper.GetString(KeyClientID) == "" {
		return nil
	}
	ctx = oidc.ClientContext(ctx, &http.Client{
		Timeout: 3 * time.Second,
	})
	p, err := NewProviderFromViper(ctx)
	if err != nil {
		return err
	}
	Provider = p
	scopes := viper.GetStringSlice(KeyScopes)
	scopes = append(scopes, oidc.ScopeOpenID)
	OAuth2Config = oauth2.Config{
		ClientID:     viper.GetString(KeyClientID),
		ClientSecret: viper.GetString(KeyClientSecret),
		RedirectURL:  "",
		Endpoint:     p.Endpoint(),
		Scopes:       scopes,
	}

	return nil
}

func NewProviderFromViper(ctx context.Context) (*oidc.Provider, error) {
	baseURL := viper.GetString(KeyBaseURL)
	if baseURL != "" {
		return oidc.NewProvider(ctx, baseURL)
	}

	providerConfig := &oidc.ProviderConfig{
		AuthURL:     viper.GetString(KeyAuthURL),
		TokenURL:    viper.GetString(KeyTokenURL),
		UserInfoURL: viper.GetString(KeyUserInfoURL),
		JWKSURL:     viper.GetString(KeyCertificatesURL),
	}

	return providerConfig.NewProvider(ctx), nil
}
