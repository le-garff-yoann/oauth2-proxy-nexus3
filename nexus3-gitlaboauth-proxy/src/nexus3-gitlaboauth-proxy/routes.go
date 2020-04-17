package main

import (
	"log"
	"net/http"
	"net/http/httputil"

	"nexus3-gitlaboauth-proxy/gitlab"
	"nexus3-gitlaboauth-proxy/nexus"

	"github.com/gorilla/mux"
)

func setupRoutes(router *mux.Router) {
	nexusConn := nexus.Conn{
		BaseURL:  cfg.NexusURL,
		Username: cfg.NexusAdminUser,
		Password: cfg.NexusAdminPassword,
	}

	router.
		PathPrefix("/").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				writeErrCb = func(msg string, code int) {
					log.Printf(`Return with HTTP %d: %s`, code, msg)

					http.Error(w, msg, code)
				}

				accessToken = r.Header.Get(cfg.GitlabAccessTokenHeader)
			)

			if accessToken == "" {
				writeErrCb("header "+cfg.GitlabAccessTokenHeader+" value is null", http.StatusBadRequest)

				return
			}

			gitlabOAuthConn := gitlab.OAuthConn{URL: cfg.GitlabURL}

			gitlabOAuthUserInfo, err := gitlabOAuthConn.GetUserInfo(accessToken)
			if err != nil {
				writeErrCb(err.Error(), http.StatusInternalServerError)

				return
			}

			if err = nexusConn.SyncUser(gitlabOAuthUserInfo.Username, gitlabOAuthUserInfo.Email, gitlabOAuthUserInfo.Groups); err != nil {
				writeErrCb(err.Error(), http.StatusInternalServerError)

				return
			}

			r.Header.Set(cfg.NexusRutHeader, gitlabOAuthUserInfo.Username)

			httputil.NewSingleHostReverseProxy(nexusConn.BaseURL).ServeHTTP(w, r)
		})
}
