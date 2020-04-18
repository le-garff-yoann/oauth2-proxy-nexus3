package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"nexus3-gitlaboauth-proxy/reverseproxy"

	env "github.com/caarlos0/env/v6"
)

var cfg = config{}

func main() {
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to parse the configuration: %s", err)
	}

	var (
		reverseProxy = reverseproxy.New(
			cfg.NexusURL, cfg.GitlabURL, cfg.NexusURL,
			cfg.GitlabAccessTokenHeader,
			cfg.NexusAdminUser, cfg.NexusAdminPassword, cfg.NexusRutHeader,
		)

		server = http.Server{Addr: cfg.ListenOn, Handler: reverseProxy.Router}
	)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: cfg.SSLInsecureSkipVerify}

	log.Println("Starting the proxy")

	log.Panicln(server.ListenAndServe())
}
