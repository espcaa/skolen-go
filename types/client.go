package types

import (
	"net/http"
	"time"
)

type Client struct {
	HTTP       *http.Client `json:"-"`
	TokenSet   TokenSet     `json:"tokenSet"`
	School     School       `json:"school"`
	BaseURL    string       `json:"-"`
	OIDCClient OIDCClient   `json:"-"`
}

type OIDCClient struct {
	ClientID     string `json:"-"`
	ClientSecret string `json:"-"`
	RedirectURI  string `json:"-"`
}

type School struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	EmsOIDCWellKnownURL string `json:"emsOIDCWellKnownUrl"`
	EmsCode             string `json:"emsCode"`
	HomePageURL         string `json:"homePageUrl"`
}

type TokenSet struct {
	AccessToken  string    `json:"access_token"`
	IDToken      string    `json:"id_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresAt    time.Time `json:"-"`
	RawExpiresAt int64     `json:"expires_at"`
	Scope        string    `json:"scope"`
}
