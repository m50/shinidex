package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	oidcc "github.com/coreos/go-oidc/v3/oidc"
	"github.com/gookit/slog"
	"github.com/labstack/echo/v4"
	"github.com/m50/shinidex/pkg/context"
	"github.com/m50/shinidex/pkg/database"
	"github.com/m50/shinidex/pkg/oidc"
	"github.com/m50/shinidex/pkg/types"
	"github.com/m50/shinidex/pkg/views"
	"github.com/m50/shinidex/pkg/web/session"
	"github.com/spf13/viper"
)

const (
	cookieState = "oidcstate"
	cookieNonce = "oidcnonce"
)

func oidcRouter(g *echo.Group) {
	if viper.GetString(oidc.KeyClientID) == "" {
		return
	}
	group := g.Group("/oidc")

	group.GET("/login", oidcLogin).Name = oidc.PathNameOIDCLogin
	group.GET("/callback", oidcRedirect).Name = oidc.PathNameOIDCCallback
}

func oidcLogin(c echo.Context) error {
	oauth2Config := oidc.OAuth2Config
	oauth2Config.RedirectURL = fmt.Sprintf("%s://%s%s", c.Scheme(), c.Request().Host, c.Echo().Reverse(oidc.PathNameOIDCCallback))
	state, err := randString(32)
	if err != nil {
		return views.RenderError(c, err)
	}
	nonce, err := randString(16)
	if err != nil {
		return views.RenderError(c, err)
	}
	c.SetCookie(&http.Cookie{
		Name:     cookieState,
		Value:    state,
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   c.Request().TLS != nil,
		HttpOnly: true,
	})
	c.SetCookie(&http.Cookie{
		Name:     cookieNonce,
		Value:    nonce,
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   c.Request().TLS != nil,
		HttpOnly: true,
	})
	return c.Redirect(http.StatusFound, oauth2Config.AuthCodeURL(state, oidcc.Nonce(nonce)))
}

func oidcRedirect(c echo.Context) error {
	ctx := context.FromEcho(c)
	log := slog.WithContext(ctx).AddData(slog.M{"oidc-state": "callback"})
	cookie, err := c.Cookie(cookieState)
	if err != nil || cookie == nil {
		log.Errorf("error getting oicdstate cookie: %v", err)
		return views.RenderView(c, http.StatusUnauthorized, LoginForm(c),
			views.Error(err))
	}
	if cookie.Value != c.QueryParam("state") {
		log.Errorf("state does not match")
		return views.RenderView(c, http.StatusUnauthorized, LoginForm(c),
			views.Error(errors.New("state does not match")))
	}

	oauth2Token, err := oidc.OAuth2Config.Exchange(ctx, c.QueryParam("code"))
	if err != nil {
		log.Error(err)
		return views.RenderView(c, http.StatusUnauthorized, LoginForm(c),
			views.Error(err))
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return views.RenderView(c, http.StatusUnauthorized, LoginForm(c),
			views.Error(errors.New("oauth2 token not found")))
	}

	verifier := oidc.Provider.VerifierContext(ctx, &oidcc.Config{
		ClientID: oidc.OAuth2Config.ClientID,
	})
	// Parse and verify ID Token payload.
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		log.Error(err)
		return views.RenderView(c, http.StatusForbidden, LoginForm(c),
			views.Error(err))
	}

	nonce, err := c.Cookie(cookieNonce)
	if err != nil {
		log.Error(err)
		return views.RenderView(c, http.StatusForbidden, LoginForm(c),
			views.Error(err))
	}
	if idToken.Nonce != nonce.Value {
		log.Error("nonce does not match")
		return views.RenderView(c, http.StatusForbidden, LoginForm(c),
			views.Error(errors.New("nonce did not match")))
	}

	var claims struct {
		Email    string `json:"email"`
		Verified bool   `json:"email_verified"`
	}
	if err := idToken.Claims(&claims); err != nil {
		log.Error(err)
		return views.RenderView(c, http.StatusForbidden, LoginForm(c),
			views.Error(err))
	}
	db := c.(database.DBContext).DB()
	user, err := db.Users().FindOrMake(ctx, types.User{
		Email:   claims.Email,
		Managed: true,
	})
	if !user.Managed {
		user.Managed = true
		if err = db.Users().Update(ctx, user); err != nil {
			slog.WithContext(ctx).Error(err)
		}
	}
	if err != nil {
		log.Error(err)
		return views.RenderView(c, http.StatusForbidden, LoginForm(c),
			views.Error(fmt.Errorf("no account found for %s", claims.Email)))
	}
	if err := session.New(c, user); err != nil {
		log.Error(err)
		return views.RenderView(c, http.StatusInternalServerError, LoginForm(c),
			views.Error(err))
	}

	log.Info("successful login")
	return c.Redirect(http.StatusMovedPermanently, "/")
}

func randString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
