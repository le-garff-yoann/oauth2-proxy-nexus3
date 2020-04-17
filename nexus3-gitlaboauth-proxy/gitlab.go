package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type gitlabOAuthConn struct {
	URL string
}

type gitlabUserInfo struct {
	Username string   `json:"nickname"`
	Email    string   `json:"email"`
	Groups   []string `json:"groups"`
}

func (s *gitlabOAuthConn) getUserInfo(accessToken string) (*gitlabUserInfo, error) {
	endpoint, err := url.Parse(fmt.Sprintf(s.URL + "/oauth/userinfo"))
	if err != nil {
		log.Fatalf("Failed to parse the GitLab URL: %s", err)
	}

	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create the GitLab GET userinfo request: %s", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to request the GitLab GET userinfo endpoint on %s: %s", s.URL, err)
	}

	if res.StatusCode != http.StatusOK {
		if resBody, err := ioutil.ReadAll(res.Body); err == nil {
			return nil, fmt.Errorf("failed to request the GitLab GET userinfo endpoint on %s: %s", s.URL, resBody)
		}

		return nil, fmt.Errorf("failed to read the GitLab GET userinfo error response: %s", err)
	}

	var userInfo gitlabUserInfo
	if err := json.NewDecoder(res.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode the GitLab GET userinfo responses: %s", err)
	}

	return &userInfo, nil
}
