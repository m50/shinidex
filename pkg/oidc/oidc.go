package oidc

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/spf13/pflag"
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

func Flags(flags *pflag.FlagSet) {
	flags.String("auth.key", string(make([]byte, 32)), "they key to be used for signing sessions (default: Regens every run)")
	flags.Bool("auth.disable-registration", false, "disable registration (default: false)")
	flags.String(KeyName, "OIDC", "the name for the OIDC login (default: OIDC)")
	flags.Bool(KeyCreateUsers, false, "automatically create users that sign in via OIDC (default: false)")
	flags.Bool(KeyDisablePassword, false, "disable password authentication")
	flags.String(KeyClientID, "", "client-id for OIDC")
	flags.String(KeyClientSecret, "", "client-secret for OIDC")
	flags.String(KeyClientSecretFile, "", "path to file containing client-secret for OIDC")
	flags.StringSlice(KeyScopes, []string{"profile", "email"}, "scopes for OIDC (default: profile,email)")
	flags.String(KeyAuthURL, "", "auth URL for OIDC")
	flags.String(KeyTokenURL, "", "token URL for OIDC")
	flags.String(KeyUserInfoURL, "", "user info URL for OIDC")
	flags.String(KeyCertificatesURL, "", "url to renew certificates for OIDC")
	flags.String(KeyBaseURL, "", "OIDC discovery URL, if provided other URLs will be ignored")
}

func SecretFile() error {
	oidcClientSecretFile := viper.GetString(KeyClientSecretFile)
	if oidcClientSecretFile == "" {
		return nil
	}
	oidcClientSecret, err := os.ReadFile(oidcClientSecretFile)
	if err != nil {
		return err
	}
	viper.Set(KeyClientSecret, oidcClientSecret)
	return nil
}

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
