package gitlab

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

// NewTestServer returns an `httptest.Server` that partially implements the GitLab OIDC API.
func NewTestServer(accessToken string, userInfo *UserInfo) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == userInfoEndpointPath {
			if r.Header["Authorization"][0] == "Bearer "+accessToken {
				payload, _ := json.Marshal(&userInfo)

				w.WriteHeader(200)
				w.Write(payload)
			} else {
				w.WriteHeader(401)
			}
		} else {
			w.WriteHeader(404)
		}
	}))
}
