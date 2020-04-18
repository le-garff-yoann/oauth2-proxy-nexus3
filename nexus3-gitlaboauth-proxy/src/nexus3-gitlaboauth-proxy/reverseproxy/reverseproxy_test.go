package reverseproxy

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"nexus3-gitlaboauth-proxy/gitlab"
	"nexus3-gitlaboauth-proxy/nexus"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	var (
		oauthAccessToken = "token"
		nexusUser        = nexus.User{
			UserID:       "foo",
			EmailAddress: "foo@test.bar",
			RoleIDs:      []string{"bar"},
		}
		nexusAvailablesRoles = []nexus.Role{{ID: nexusUser.RoleIDs[0]}}

		gitLabOAuthTestSrv = gitlab.NewTestServer(oauthAccessToken, &gitlab.OAuthUserInfo{
			Username: nexusUser.UserID,
			Email:    nexusUser.EmailAddress,
			Groups:   nexusUser.RoleIDs,
		})
		gitLabOAuthTestSrvURL, _ = url.Parse(gitLabOAuthTestSrv.URL)

		nexusTestSrv = nexus.NewTestServer(
			[]nexus.UserModifier{{User: nexusUser}},
			&nexusAvailablesRoles,
		)
		nexusTestSrvURL, _ = url.Parse(nexusTestSrv.URL)

		rproxyAccessTokenHeader = "X-Forwarded-Access-Token"
		rproxy                  = New(
			nexusTestSrvURL, gitLabOAuthTestSrvURL, nexusTestSrvURL,
			rproxyAccessTokenHeader,
			"null", "null", "X-Forwarded-User",
		)

		rProxySrv = httptest.NewServer(rproxy.Router.GetRoute(routeName).GetHandler())
	)

	defer gitLabOAuthTestSrv.Close()
	defer nexusTestSrv.Close()
	defer rProxySrv.Close()

	res, err := rProxySrv.Client().Get(rProxySrv.URL)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	sucessfulReq, _ := http.NewRequest("GET", rProxySrv.URL, nil)
	sucessfulReq.Header.Add(rproxyAccessTokenHeader, oauthAccessToken)

	res, err = rProxySrv.Client().Do(sucessfulReq)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
}
