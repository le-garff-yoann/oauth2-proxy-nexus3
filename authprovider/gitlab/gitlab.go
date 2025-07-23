package gitlab

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/le-garff-yoann/oauth2-proxy-nexus3/authprovider"
)

const userInfoEndpointPath = "/oauth/userinfo"

// Client implements `authprovider.Client`.
type Client struct {
	URL *url.URL
}

// GetUserInfo implements `authprovider.Client`.
func (s *Client) GetUserInfo(accessToken string) (authprovider.UserInfo, error) {
	endpoint, err := url.Parse(fmt.Sprintf(s.URL.String() + userInfoEndpointPath))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse the GitLab URL: %s", err)
	}

	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create the GitLab GET userinfo request: %s", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request the GitLab GET userinfo endpoint on %s: %s", s.URL, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		if resBody, err := ioutil.ReadAll(res.Body); err == nil {
			return nil, fmt.Errorf("failed to request the GitLab GET userinfo endpoint on %s: %s", s.URL, resBody)
		}

		return nil, fmt.Errorf("failed to read the GitLab GET userinfo error response: %s", err)
	}

	var userInfo UserInfo
	if err := json.NewDecoder(res.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode the GitLab GET userinfo responses: %s", err)
	}

	return &userInfo, nil
}

// UserInfo implements `authprovider.UserInfo`.
type UserInfo struct {
	Nickname string   `json:"nickname"`
	Email    string   `json:"email"`
	Groups   []string `json:"groups"`
}

// Username implements `authprovider.UserInfo`.
func (s *UserInfo) Username() string {
	return s.Nickname
}

// EmailAddress implements `authprovider.UserInfo`.
func (s *UserInfo) EmailAddress() string {
	return s.Email
}

// Roles implements `authprovider.UserInfo`.
func (s *UserInfo) Roles() []string {
	return s.Groups
}
