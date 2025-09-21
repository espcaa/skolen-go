package skolengo

import (
	"context"
	"errors"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

func RefreshAccessToken(c *Client) error {
	if c.TokenSet.RefreshToken == "" {
		return errors.New("no refresh token available")
	}

	provider, err := oidc.NewProvider(context.Background(), c.School.EmsOIDCWellKnownURL)
	if err != nil {
		return err
	}

	oauth2Config := &oauth2.Config{
		ClientID:     c.OIDCClient.ClientID,
		ClientSecret: c.OIDCClient.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  c.OIDCClient.RedirectURI,
		Scopes:       []string{"openid", "profile"},
	}

	tokenSource := oauth2Config.TokenSource(context.Background(), &oauth2.Token{
		RefreshToken: c.TokenSet.RefreshToken,
	})

	newToken, err := tokenSource.Token()
	if err != nil {
		return err
	}

	c.TokenSet.AccessToken = newToken.AccessToken
	if newToken.RefreshToken != "" {
		c.TokenSet.RefreshToken = newToken.RefreshToken
	}
	c.TokenSet.ExpiresAt = newToken.Expiry
	c.TokenSet.TokenType = newToken.TokenType

	if scope, ok := newToken.Extra("scope").(string); ok {
		c.TokenSet.Scope = scope
	}

	if idToken, ok := newToken.Extra("id_token").(string); ok {
		c.TokenSet.IDToken = idToken
	}

	return nil
}
