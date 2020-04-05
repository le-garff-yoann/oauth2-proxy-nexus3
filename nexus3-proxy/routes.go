package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"
)

func setupRoutes(router *mux.Router) {
	nexusUserUpstream, err := url.Parse(envConfig["N3P_NEXUS3_UPSTREAM"])
	if err != nil {
		log.Fatal(err)
	}

	limiter := make(map[string]interface{})

	router.
		PathPrefix("/").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writeErrCb := func(msg string, code int) {
				log.Printf(`Return with HTTP %d: %s`, code, msg)

				http.Error(w, msg, code)
			}

			nexusUserInfo := nexusUserInfo{
				Username: r.Header.Get(envConfig["N3P_NEXUS3_RUT_USER_HEADER"]),
				Email:    r.Header.Get(envConfig["N3P_NEXUS3_RUT_EMAIL_HEADER"]),
			}

			if nexusUserInfo.Username == "" || nexusUserInfo.Email == "" {
				writeErrCb(fmt.Sprintf("Header %s or %s value is null",
					envConfig["N3P_NEXUS3_RUT_USER_HEADER"], envConfig["N3P_NEXUS3_RUT_EMAIL_HEADER"]), http.StatusBadRequest)
			}

			if _, ok := limiter[nexusUserInfo.Username]; !ok {
				if err = createNexusUser(&nexusUserInfo); err != nil {
					writeErrCb(err.Error(), http.StatusInternalServerError)

					return
				}

				limiter[nexusUserInfo.Username] = nil
			}

			httputil.NewSingleHostReverseProxy(nexusUserUpstream).ServeHTTP(w, r)
		})
}
