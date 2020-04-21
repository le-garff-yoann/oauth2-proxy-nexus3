package gitlab

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetUserInfo(t *testing.T) {
	t.Parallel()

	var (
		expectedAccessToken = "expectedToken"
		expectedUserInfo    = UserInfo{
			Nickname: "foo",
			Email:    "foo@test.bar",
			Groups:   []string{"bar"},
		}

		srv = NewTestServer(expectedAccessToken, &expectedUserInfo)

		srvURL, _ = url.Parse(srv.URL)
		client    = Client{URL: srvURL}
	)

	defer srv.Close()

	userInfo, err := client.GetUserInfo(expectedAccessToken)
	require.NoError(t, err)
	require.Equal(t, &expectedUserInfo, userInfo)

	_, err = client.GetUserInfo("unexpectedToken")
	require.Error(t, err)

	badURL, _ := url.Parse("/bad-url")
	client = Client{URL: badURL}

	_, err = client.GetUserInfo(expectedAccessToken)
	require.Error(t, err)
}
