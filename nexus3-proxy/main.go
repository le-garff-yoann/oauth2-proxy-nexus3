package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	setConfig()

	envConfig["N3P_LISTEN_ON"] = os.Getenv("N3P_LISTEN_ON")
	envConfig["N3P_SSL_INSECURE_SKIP_VERIFY"] = os.Getenv("N3P_SSL_INSECURE_SKIP_VERIFY")

	var (
		router = mux.NewRouter().StrictSlash(true)

		server = http.Server{Addr: envConfig["N3P_LISTEN_ON"], Handler: router}
	)

	if envConfig["N3P_SSL_INSECURE_SKIP_VERIFY"] == "true" {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	log.Println("Starting the proxy")

	setupRoutes(router)
	log.Panicln(server.ListenAndServe())
}
