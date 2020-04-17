package main

import (
	"crypto/tls"
	"log"
	"net/http"

	env "github.com/caarlos0/env/v6"
	"github.com/gorilla/mux"
)

var cfg = config{}

func main() {
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to parse the configuration: %s", err)
	}

	var (
		router = mux.NewRouter().StrictSlash(true)

		server = http.Server{Addr: cfg.ListenOn, Handler: router}
	)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: cfg.SSLInsecureSkipVerify}

	log.Println("Starting the proxy")

	setupRoutes(router)
	log.Panicln(server.ListenAndServe())
}
