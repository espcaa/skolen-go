package skolengo

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
)

func RefreshAccessToken(c *Client) error {
	if c.TokenSet.RefreshToken == "" {
		return errors.New("no refresh token available")
	}

	provider, err := oidc.NewProvider(context.Background(), c.School.EmsOIDCWellKnownURL)
	if err != nil {
		return err
	}
	tokenURL := provider.Endpoint().TokenURL

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", c.TokenSet.RefreshToken)
	data.Set("client_id", c.OIDCClient.ClientID)
	if c.OIDCClient.ClientSecret != "" {
		data.Set("client_secret", c.OIDCClient.ClientSecret)
	}

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("refresh failed with status " + resp.Status)
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int64  `json:"expires_in"`
		TokenType    string `json:"token_type"`
		Scope        string `json:"scope"`
		IDToken      string `json:"id_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return err
	}

	c.TokenSet.AccessToken = tokenResp.AccessToken
	if tokenResp.RefreshToken != "" {
		c.TokenSet.RefreshToken = tokenResp.RefreshToken
	}
	if tokenResp.ExpiresIn > 0 {
		c.TokenSet.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	}
	c.TokenSet.TokenType = tokenResp.TokenType
	c.TokenSet.Scope = tokenResp.Scope
	c.TokenSet.IDToken = tokenResp.IDToken

	return nil
}
