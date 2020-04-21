package main

import "net/url"

type config struct {
	ListenOn                      string   `env:"O2PN3_LISTEN_ON" envDefault:"0.0.0.0:8080"`
	SSLInsecureSkipVerify         bool     `env:"O2PN3_SSL_INSECURE_SKIP_VERIFY" envDefault:"false"`
	AuthproviderURL               *url.URL `env:"O2PN3_AP_URL,required"`
	AuthproviderAccessTokenHeader string   `env:"O2PN3_AP_ACCESS_TOKEN_HEADER" envDefault:"X-Forwarded-Access-Token"`
	NexusURL                      *url.URL `env:"O2PN3_NEXUS3_URL,required"`
	NexusAdminUser                string   `env:"O2PN3_NEXUS3_ADMIN_USER,required"`
	NexusAdminPassword            string   `env:"O2PN3_NEXUS3_ADMIN_PASSWORD,required"`
	NexusRutHeader                string   `env:"O2PN3_NEXUS3_RUT_HEADER" envDefault:"X-Forwarded-User"`
}
