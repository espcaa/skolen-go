package skolengo

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/espcaa/skolen-go/types"
)

type Client struct {
	HTTP       *http.Client    `json:"-"`
	TokenSet   TokenSet        `json:"tokenSet"`
	School     School          `json:"school"`
	BaseURL    string          `json:"-"`
	OIDCClient OIDCClient      `json:"-"`
	UserInfo   *types.UserInfo `json:"userInfo,omitempty"`
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

func (t *TokenSet) UnmarshalJSON(data []byte) error {
	type Alias TokenSet
	aux := &struct {
		ExpiresAt int64 `json:"expires_at"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	t.ExpiresAt = time.Unix(aux.ExpiresAt, 0)
	return nil
}

func NewClientFromJSON(data []byte) (*Client, error) {
	var client Client
	if err := json.Unmarshal(data, &client); err != nil {
		return nil, err
	}

	client.HTTP = &http.Client{
		Timeout: 10 * time.Second,
	}

	client.BaseURL = "https://api.skolengo.com/api/v1/bff-sko-app"

	client.OIDCClient.ClientID = SkolenGoConstants.OIDCClientID
	client.OIDCClient.ClientSecret = SkolenGoConstants.OIDCClientSecret
	client.OIDCClient.RedirectURI = SkolenGoConstants.RedirectURI

	// get userinfo

	var userInfo, err = client.GetBasicUserInfo()
	if err != nil {
		return nil, err
	}
	client.UserInfo = userInfo

	return &client, nil
}
