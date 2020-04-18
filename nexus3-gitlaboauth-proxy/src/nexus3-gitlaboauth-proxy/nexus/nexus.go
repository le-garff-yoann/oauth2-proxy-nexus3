package nexus

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
	anonymousRoleID       = "nx-anonymous"
	defaultSourceName     = "default"
	userStatusActiveValue = "active"

	userEndpointPath = "/service/rest/beta/security/users"
	roleEndpointPath = "/service/rest/beta/security/roles"
)

// Conn represents a connexion to Nexus 3.
type Conn struct {
	BaseURL *url.URL

	Username, Password string
}

// User represents a Nexus 3 user.
type User struct {
	UserID       string   `json:"userId"`
	FirstName    string   `json:"firstName"`
	LastName     string   `json:"lastName"`
	EmailAddress string   `json:"emailAddress"`
	Password     string   `json:"password"`
	Status       string   `json:"status"`
	RoleIDs      []string `json:"roles"`
}

// Role partially represents a Nexus 3 role.
type Role struct {
	ID string `json:"id"`
}

// UserModifier partially represents a
// Nexus 3 user "modifier" (used in PUT requests).
type UserModifier struct {
	User

	Source string `json:"source"`
}

func (s *Conn) newUserEndpointURL() *url.URL {
	url, _ := url.Parse(s.BaseURL.String() + userEndpointPath)

	return url
}

func (s *Conn) getUser(userID string) (*User, error) {
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
	if err != nil {
		return nil, fmt.Errorf("failed to request the Nexus 3 GET user endpoint on %s: %s", endpoint.String(), err)
	}
	defer res.Body.Close()

	var users []User
	if err := json.NewDecoder(res.Body).Decode(&users); err != nil {
		return nil, fmt.Errorf("failed to decode the Nexus 3 GET user response: %s", err)
	}

	if c := len(users); c == 1 {
		return &users[0], nil
	} else if c > 1 {
		return nil, errors.New(endpoint.String() + ": one user should have been returned by the Nexus 3 GET user request")
	}

	return nil, nil
}

func (s *Conn) createUser(user *User) error {
	endpoint := s.newUserEndpointURL()

	reqBody, err := json.Marshal(user)
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
	if err != nil {
		return fmt.Errorf("failed to request the Nexus 3 POST user endpoint on %s: %s", endpoint.String(), err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		if resBody, err := ioutil.ReadAll(res.Body); err == nil {
			return fmt.Errorf(`cannot create the user "%s" on %s: %s`, user.UserID, s.BaseURL.String(), resBody)
		}

		return fmt.Errorf("failed to read the Nexus 3 POST user error response: %s", err)
	}

	return nil
}

func (s *Conn) modifyUser(username string, userModifier *UserModifier) error {
	endpoint := s.newUserEndpointURL()
	endpoint.Path += "/" + username

	reqBody, err := json.Marshal(userModifier)
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
	if err != nil {
		return fmt.Errorf("failed to request the Nexus 3 PUT user endpoint on %s: %s", endpoint.String(), err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		if resBody, err := ioutil.ReadAll(res.Body); err == nil {
			return fmt.Errorf(`cannot modify the user "%s" on %s (HTTP %d): %s`, userModifier.UserID, s.BaseURL.String(), res.StatusCode, resBody)
		}

		return fmt.Errorf("failed to read the Nexus 3 PUT user error response: %s", err)
	}

	return nil
}

func (s *Conn) getRoles() ([]Role, error) {
	endpoint, _ := url.Parse(s.BaseURL.String() + roleEndpointPath)

	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create the Nexus 3 GET roles request: %s", err)
	}

	req.SetBasicAuth(s.Username, s.Password)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request the Nexus 3 GET roles endpoint on %s: %s", endpoint.String(), err)
	}
	defer res.Body.Close()

	var roles []Role
	if err := json.NewDecoder(res.Body).Decode(&roles); err != nil {
		return nil, fmt.Errorf("failed to decode the Nexus 3 GET roles response: %s", err)
	}

	return roles, nil
}

func (s *Conn) userModifier(oldUser, newUser *User, existingRoles []Role) (bool, *UserModifier) {
	var (
		userModifier = UserModifier{
			User: User{
				UserID:       newUser.UserID,
				FirstName:    newUser.FirstName,
				LastName:     newUser.LastName,
				EmailAddress: newUser.EmailAddress,
				Password:     oldUser.Password,
				Status:       newUser.Status,
			},
			Source: defaultSourceName,
		}

		shouldModify = oldUser.UserID != newUser.UserID ||
			oldUser.FirstName != newUser.FirstName ||
			oldUser.LastName != newUser.LastName ||
			oldUser.EmailAddress != newUser.EmailAddress ||
			oldUser.Status != newUser.Status
	)

	var existingRoleIDs []string
	for _, role := range existingRoles {
		existingRoleIDs = append(existingRoleIDs, role.ID)
	}

	for _, roleID := range newUser.RoleIDs {
		if funk.ContainsString(existingRoleIDs, roleID) {
			userModifier.RoleIDs = append(userModifier.RoleIDs, roleID)
		}
	}

	if len(userModifier.RoleIDs) == 0 {
		shouldModify = true

		userModifier.RoleIDs = append(userModifier.RoleIDs, anonymousRoleID)
	} else {
		userModifier.RoleIDs = funk.UniqString(userModifier.RoleIDs)

		if !shouldModify {
			oldRoleIDsDiff, newRoleIDsDiff := funk.DifferenceString(oldUser.RoleIDs, userModifier.RoleIDs)

			shouldModify = len(oldRoleIDsDiff) > 0 && len(newRoleIDsDiff) > 0
		}
	}

	return shouldModify, &userModifier
}

// SyncUser "synchronizes" the user on Nexus 3
// based on the parameters passed to this method.
func (s *Conn) SyncUser(username, email string, roleIDs []string) error {
	user := &User{
		UserID:       username,
		FirstName:    username,
		LastName:     username,
		EmailAddress: email,
		Status:       userStatusActiveValue,
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
			user.RoleIDs = []string{anonymousRoleID}
		} else {
			user.RoleIDs = userModifier.RoleIDs
		}

		user.Password = uuid.New().String()

		return s.createUser(user)
	}

	if shouldModify, userModifier := s.userModifier(originalUser, user, existingRoles); shouldModify {
		return s.modifyUser(user.UserID, userModifier)
	}

	return nil
}
