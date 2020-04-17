package main

import "net/url"

type config struct {
	ListenOn                string   `env:"N3GOP_LISTEN_ON" envDefault:"0.0.0.0:8080"`
	SSLInsecureSkipVerify   bool     `env:"N3GOP_SSL_INSECURE_SKIP_VERIFY" envDefault:"false"`
	GitlabURL               *url.URL `env:"N3GOP_GITLAB_URL,required"`
	GitlabAccessTokenHeader string   `env:"N3GOP_GITLAB_ACCESS_TOKEN_HEADER" envDefault:"X-Forwarded-Access-Token"`
	NexusURL                *url.URL `env:"N3GOP_NEXUS3_URL,required"`
	NexusAdminUser          string   `env:"N3GOP_NEXUS3_ADMIN_USER,required"`
	NexusAdminPassword      string   `env:"N3GOP_NEXUS3_ADMIN_PASSWORD,required"`
	NexusRutHeader          string   `env:"N3GOP_NEXUS3_RUT_HEADER" envDefault:"X-Forwarded-User"`
}
