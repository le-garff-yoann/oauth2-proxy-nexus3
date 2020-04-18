package nexus

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUsers(t *testing.T) {
	t.Parallel()

	var (
		expectedUser = &User{
			UserID:       "foo",
			FirstName:    "foo",
			LastName:     "foo",
			EmailAddress: "foo@test.bar",
			Password:     "",
			Status:       userStatusActiveValue,
			RoleIDs:      []string{"bar"},
		}

		srv = NewTestServer([]UserModifier{}, &[]Role{})

		srvURL, _ = url.Parse(srv.URL)
		client    = Conn{BaseURL: srvURL}
	)

	defer srv.Close()

	user, err := client.getUser(expectedUser.UserID)
	require.NoError(t, err)
	require.Nil(t, user)

	require.NoError(t, client.createUser(expectedUser))

	user, err = client.getUser(expectedUser.UserID)
	require.NoError(t, err)
	require.Equal(t, expectedUser, user)

	expectedUser.LastName = "bar"

	require.NoError(t, client.modifyUser(expectedUser.UserID, &UserModifier{User: *expectedUser, Source: defaultSourceName}))

	user, err = client.getUser(expectedUser.UserID)
	require.NoError(t, err)
	require.Equal(t, expectedUser, user)
}

func TestRoles(t *testing.T) {
	t.Parallel()

	var (
		existingRoles = []Role{}

		srv = NewTestServer([]UserModifier{}, &existingRoles)

		srvURL, _ = url.Parse(srv.URL)
		client    = Conn{BaseURL: srvURL}
	)

	defer srv.Close()

	roles, err := client.getRoles()
	require.NoError(t, err)
	require.Empty(t, roles)

	existingRoles = append(existingRoles, Role{ID: "foo"})

	roles, err = client.getRoles()
	require.NoError(t, err)
	require.Equal(t, roles, existingRoles)
}

func TestSyncUser(t *testing.T) {
	t.Parallel()

	var (
		existingUser = &User{
			UserID:       "foo",
			FirstName:    "foo",
			LastName:     "foo",
			EmailAddress: "foo@test.bar",
			Password:     "",
			Status:       userStatusActiveValue,
			RoleIDs:      []string{"bar"},
		}
		existingRoles = []Role{{"bar"}, {"foo"}}

		srv = NewTestServer(
			[]UserModifier{{User: *existingUser, Source: defaultSourceName}},
			&existingRoles,
		)

		srvURL, _ = url.Parse(srv.URL)
		client    = Conn{BaseURL: srvURL}
	)

	defer srv.Close()

	expectedUser, err := client.getUser(existingUser.UserID)
	require.NoError(t, err)
	require.NotNil(t, expectedUser)

	newUserEmail := "foo@test.foo"
	require.NoError(t, client.SyncUser(expectedUser.UserID, newUserEmail, expectedUser.RoleIDs))

	syncedUser, err := client.getUser(existingUser.UserID)
	require.NoError(t, err)
	require.Equal(t, newUserEmail, syncedUser.EmailAddress)

	newRoleIDs := []string{"foo", "unexpectedRoleID"}
	require.NoError(t, client.SyncUser(syncedUser.UserID, syncedUser.EmailAddress, newRoleIDs))

	syncedUser, err = client.getUser(syncedUser.UserID)
	require.NoError(t, err)
	require.Len(t, syncedUser.RoleIDs, 1)
	require.Equal(t, syncedUser.RoleIDs[0], newRoleIDs[0])
}
