package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	funk "github.com/thoas/go-funk"
)

const (
	nexusAnonymousRoleID       = "nx-anonymous"
	nexusDefaultSourceName     = "default"
	nexusUserStatusActiveValue = "active"
)

type nexusConn struct {
	BaseURL *url.URL

	Username, Password string
}

type nexusUser struct {
	UserID       string   `json:"userId"`
	FirstName    string   `json:"firstName"`
	LastName     string   `json:"lastName"`
	EmailAddress string   `json:"emailAddress"`
	Password     string   `json:"password"`
	Status       string   `json:"status"`
	RoleIDs      []string `json:"roles"`
}

type nexusRole struct {
	ID string `json:"id"`
}

type nexusUserModifier struct {
	nexusUser

	Source string `json:"source"`
}

func newNexusConn(baseURL, Username, Password string) (*nexusConn, error) {
	url, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the Nexus 3 base URL: %s", err)
	}

	return &nexusConn{url, Username, Password}, nil
}

func (s *nexusConn) newUserEndpointURL() *url.URL {
	url, _ := url.Parse(s.BaseURL.String() + "/service/rest/beta/security/users")

	return url
}

func (s *nexusConn) getUser(userID string) (*nexusUser, error) {
	endpoint := s.newUserEndpointURL()

	endpointQuery, _ := url.ParseQuery(endpoint.RawQuery)
	endpointQuery.Add("userId", userID)

	endpoint.RawQuery = endpointQuery.Encode()

	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create the Nexus 3 GET user request: %s", err)
	}

	req.SetBasicAuth(s.Username, s.Password)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to request the Nexus 3 GET user endpoint on %s: %s", endpoint.String(), err)
	}

	var users []nexusUser
	if err := json.NewDecoder(res.Body).Decode(&users); err != nil {
		return nil, fmt.Errorf("failed to decode the Nexus 3 GET user response: %d", err)
	}

	if c := len(users); c == 1 {
		return &users[0], nil
	} else if c > 1 {
		return nil, errors.New(endpoint.String() + ": one user should have been returned by the Nexus 3 GET user request")
	}

	return nil, nil
}

func (s *nexusConn) createUser(user nexusUser) error {
	endpoint := s.newUserEndpointURL()

	reqBody, err := json.Marshal(&user)
	if err != nil {
		return fmt.Errorf(`failed to encode %#v to JSON`, user)
	}

	req, err := http.NewRequest("POST", endpoint.String(), bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create the Nexus 3 POST user request: %s", err)
	}

	req.SetBasicAuth(s.Username, s.Password)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to request the Nexus 3 POST user endpoint on %s: %s", endpoint.String(), err)
	}

	if res.StatusCode != http.StatusOK {
		if resBody, err := ioutil.ReadAll(res.Body); err == nil {
			return fmt.Errorf(`cannot create the user "%s" on %s: %s`, user.UserID, s.BaseURL.String(), resBody)
		}

		return fmt.Errorf("failed to read the Nexus 3 POST user error response: %s", err)
	}

	return nil
}

func (s *nexusConn) modifyUser(userModifier nexusUserModifier) error {
	endpoint := s.newUserEndpointURL()

	reqBody, err := json.Marshal(&userModifier)
	if err != nil {
		return fmt.Errorf(`failed to encode %#v to JSON`, userModifier)
	}

	req, err := http.NewRequest("PUT", endpoint.String(), bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create the Nexus 3 PUT user request: %s", err)
	}

	req.SetBasicAuth(s.Username, s.Password)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to request the Nexus 3 PUT user endpoint on %s: %s", endpoint.String(), err)
	}

	if res.StatusCode != http.StatusOK {
		if resBody, err := ioutil.ReadAll(res.Body); err == nil {
			return fmt.Errorf(`cannot modify the user "%s" on %s: %s`, userModifier.UserID, s.BaseURL.String(), resBody)
		}

		return fmt.Errorf("failed to read the Nexus 3 PUT user error response: %s", err)
	}

	return nil
}

func (s *nexusConn) getRoles() ([]nexusRole, error) {
	endpoint, _ := url.Parse(s.BaseURL.String() + "/service/rest/beta/security/roles")

	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create the Nexus 3 GET roles request: %s", err)
	}

	req.SetBasicAuth(s.Username, s.Password)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to request the Nexus 3 GET roles endpoint on %s: %s", endpoint.String(), err)
	}

	var roles []nexusRole
	if err := json.NewDecoder(res.Body).Decode(&roles); err != nil {
		return nil, fmt.Errorf("failed to decode the Nexus 3 GET roles response: %d", err)
	}

	return roles, nil
}

func (s *nexusConn) userModifier(oldUser, newUser nexusUser, existingRoles []nexusRole) (bool, *nexusUserModifier) {
	var (
		ok = false

		userModifier = nexusUserModifier{}
	)

	userModifier.Password = oldUser.Password
	userModifier.Source = nexusDefaultSourceName

	if oldUser.UserID != newUser.UserID {
		ok = true

		userModifier.UserID = newUser.UserID
	}
	if oldUser.FirstName != newUser.FirstName {
		ok = true

		userModifier.FirstName = newUser.FirstName
	}
	if oldUser.LastName != newUser.LastName {
		ok = true

		userModifier.LastName = newUser.LastName
	}
	if oldUser.Status != newUser.Status {
		ok = true

		userModifier.Status = newUser.Status
	}

	for _, roleID := range oldUser.RoleIDs {
		userModifier.RoleIDs = append(userModifier.RoleIDs, roleID)
	}

	var existingRoleIDs []string
	for _, role := range existingRoles {
		existingRoleIDs = append(existingRoleIDs, role.ID)
	}

	for _, roleID := range newUser.RoleIDs {
		if funk.ContainsString(existingRoleIDs, roleID) {
			userModifier.RoleIDs = append(userModifier.RoleIDs, roleID)
		}
	}

	userModifier.RoleIDs = funk.UniqString(userModifier.RoleIDs)

	return ok || len(userModifier.RoleIDs) > len(oldUser.RoleIDs), &userModifier
}

func (s *nexusConn) syncUser(username, email string, roleIDs []string) error {
	user := nexusUser{
		UserID:       username,
		FirstName:    username,
		LastName:     username,
		EmailAddress: email,
		Status:       nexusUserStatusActiveValue,
		RoleIDs:      roleIDs,
	}

	originalUser, err := s.getUser(user.UserID)
	if err != nil {
		return err
	}

	existingRoles, err := s.getRoles()
	if err != nil {
		return err
	}

	if originalUser == nil {
		_, userModifier := s.userModifier(user, user, existingRoles)

		if len(userModifier.RoleIDs) == 0 {
			user.RoleIDs = []string{nexusAnonymousRoleID}
		} else {
			user.RoleIDs = userModifier.RoleIDs
		}

		user.Password = uuid.New().String()

		return s.createUser(user)
	}

	if shouldModify, userModifier := s.userModifier(*originalUser, user, existingRoles); shouldModify {
		return s.modifyUser(*userModifier)
	}

	return nil
}
