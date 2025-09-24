package skolengo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/espcaa/skolen-go/types"
	"github.com/golang-jwt/jwt/v4"
)

func (c *Client) GetBasicUserInfo() (*types.UserInfo, error) {
	if c.TokenSet.AccessToken == "" {
		return nil, fmt.Errorf("access token not set")
	}

	token, _, err := new(jwt.Parser).ParseUnverified(c.TokenSet.AccessToken, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("could not parse claims")
	}
	userID, ok := claims["sub"].(string)
	if !ok {
		return nil, fmt.Errorf("sub claim missing from token")
	}

	baseUrl := strings.TrimRight(c.BaseURL, "/")
	endpoint := fmt.Sprintf("%s/users-info/%s", baseUrl, userID)

	fmt.Println(endpoint)

	query := url.Values{}
	query.Set("include", "school,students,students.school,schools,prioritySchool")
	query.Set("fields[userInfo]", "lastName,firstName,photoUrl,externalMail,mobilePhone,audienceId,permissions")
	query.Set("fields[school]", "name,timeZone,subscribedServices,city,schoolAudience,administrativeId")
	query.Set("fields[legalRepresentativeUserInfo]", "addressLines,postalCode,city,country,students")
	query.Set("fields[studentUserInfo]", "className,dateOfBirth,regime,school")
	query.Set("fields[teacherUserInfo]", "schools,prioritySchool")
	query.Set("fields[localAuthorityStaffUserInfo]", "schools,prioritySchool")
	query.Set("fields[nonTeachingStaffUserInfo]", "schools,prioritySchool")
	query.Set("fields[otherPersonUserInfo]", "schools,prioritySchool")
	query.Set("fields[student]", "firstName,lastName,photoUrl,className,dateOfBirth,regime,school")

	req, err := http.NewRequest("GET", endpoint+"?"+query.Encode(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.TokenSet.AccessToken)
	req.Header.Set("x-skolengo-ems-code", c.School.EmsCode)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// Read the body
		return nil, fmt.Errorf("unexpected status %d:", resp.StatusCode)
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
