package skolengo

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/espcaa/skolen-go/types"
)

func (c *Client) GetBasicUserInfo() (*types.UserInfo, error) {
	if c.TokenSet.AccessToken == "" {
		return nil, fmt.Errorf("access token not set")
	}

	req, err := http.NewRequest("GET", c.BaseURL+"/users-info", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.TokenSet.AccessToken)
	req.Header.Set("x-skolengo-ems-code", c.School.EmsOIDCWellKnownURL)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var payload struct {
		Data struct {
			ID         string `json:"id"`
			Attributes struct {
				FirstName string `json:"firstName"`
				LastName  string `json:"lastName"`
			} `json:"attributes"`
			Relationships struct {
				School struct {
					Data struct {
						ID string `json:"id"`
					} `json:"data"`
				} `json:"school"`
			} `json:"relationships"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	user := &types.UserInfo{
		UserID:   payload.Data.ID,
		FullName: payload.Data.Attributes.FirstName + " " + payload.Data.Attributes.LastName,
		SchoolID: payload.Data.Relationships.School.Data.ID,
		EMSCode:  c.School.EmsCode,
	}

	return user, nil
}
