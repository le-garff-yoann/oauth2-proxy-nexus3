package main

import (
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/gorilla/mux"
)

func setupRoutes(router *mux.Router) {
	nexusConn, err := newNexusConn(envConfig["N3GOP_NEXUS3_URL"], envConfig["N3GOP_NEXUS3_ADMIN_USER"], envConfig["N3GOP_NEXUS3_ADMIN_PASSWORD"])
	if err != nil {
		log.Fatalf("Failed to initialize the Nexus 3 connection: %s", err)
	}

	router.
		PathPrefix("/").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				writeErrCb = func(msg string, code int) {
					log.Printf(`Return with HTTP %d: %s`, code, msg)

					http.Error(w, msg, code)
				}

				accessToken = r.Header.Get(envConfig["N3GOP_GITLAB_ACCESS_TOKEN_HEADER"])
			)

			if accessToken == "" {
				writeErrCb("header "+envConfig["N3GOP_GITLAB_ACCESS_TOKEN_HEADER"]+" value is null", http.StatusBadRequest)

				return
			}

			gitlabOAuthConn := gitlabOAuthConn{URL: envConfig["N3GOP_GITLAB_URL"]}

			gitlabUserInfo, err := gitlabOAuthConn.getUserInfo(accessToken)
			if err != nil {
				writeErrCb(err.Error(), http.StatusInternalServerError)

				return
			}

			if err = nexusConn.syncUser(gitlabUserInfo.Username, gitlabUserInfo.Email, gitlabUserInfo.Groups); err != nil {
				writeErrCb(err.Error(), http.StatusInternalServerError)

				return
			}

			r.Header.Set(envConfig["N3GOP_NEXUS3_RUT_USER_HEADER"], gitlabUserInfo.Username)
			r.Header.Set(envConfig["N3GOP_NEXUS3_RUT_EMAIL_HEADER"], gitlabUserInfo.Email)

			httputil.NewSingleHostReverseProxy(nexusConn.BaseURL).ServeHTTP(w, r)
		})
}
