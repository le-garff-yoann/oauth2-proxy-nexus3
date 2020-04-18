package reverseproxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"nexus3-gitlaboauth-proxy/gitlab"
	"nexus3-gitlaboauth-proxy/nexus"

	"github.com/gorilla/mux"
)

const routeName = "main"

// ReverseProxy It represents the reverse proxy which
// is the "glue" between oauth2-proxy, GitLab and Nexus 3.
type ReverseProxy struct {
	Router *mux.Router
}

// New initializes and returns a new `ReverseProxy`.
func New(
	upstreamURL, gitlabURL, nexusURL *url.URL,
	gitlabAccessTokenHeader, nexusAdminUser, nexusAdminPassword, nexusRutHeader string,
) *ReverseProxy {
	s := ReverseProxy{
		Router: mux.NewRouter().StrictSlash(true),
	}

	nexusConn := nexus.Conn{
		BaseURL:  nexusURL,
		Username: nexusAdminUser,
		Password: nexusAdminPassword,
	}

	s.Router.
		PathPrefix("/").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				writeErrCb = func(msg string, code int) {
					log.Printf(`Return with HTTP %d: %s`, code, msg)

					http.Error(w, msg, code)
				}

				accessToken = r.Header.Get(gitlabAccessTokenHeader)
			)

			if accessToken == "" {
				writeErrCb("header "+gitlabAccessTokenHeader+" value is null", http.StatusBadRequest)

				return
			}

			gitlabOAuthConn := gitlab.OAuthConn{URL: gitlabURL}

			gitlabOAuthUserInfo, err := gitlabOAuthConn.GetUserInfo(accessToken)
			if err != nil {
				writeErrCb(err.Error(), http.StatusInternalServerError)

				return
			}

			if err = nexusConn.SyncUser(
				gitlabOAuthUserInfo.Username,
				gitlabOAuthUserInfo.Email,
				gitlabOAuthUserInfo.Groups,
			); err != nil {
				writeErrCb(err.Error(), http.StatusInternalServerError)

				return
			}

			r.Header.Set(nexusRutHeader, gitlabOAuthUserInfo.Username)

			httputil.NewSingleHostReverseProxy(upstreamURL).ServeHTTP(w, r)
		}).
		Name(routeName)

	return &s
}
