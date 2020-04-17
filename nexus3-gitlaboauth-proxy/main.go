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

	var (
		router = mux.NewRouter().StrictSlash(true)

		server = http.Server{Addr: os.Getenv("N3GOP_LISTEN_ON"), Handler: router}
	)

	if os.Getenv("N3GOP_SSL_INSECURE_SKIP_VERIFY") == "true" {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	log.Println("Starting the proxy")

	setupRoutes(router)
	log.Panicln(server.ListenAndServe())
}
