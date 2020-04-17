package main

import (
	"log"
	"os"
)

var envConfig = make(map[string]string)

func setConfig() {
	for _, envName := range []string{
		"N3GOP_GITLAB_URL", "N3GOP_GITLAB_ACCESS_TOKEN_HEADER",
		"N3GOP_NEXUS3_URL", "N3GOP_NEXUS3_ADMIN_USER", "N3GOP_NEXUS3_ADMIN_PASSWORD",
		"N3GOP_NEXUS3_RUT_USER_HEADER", "N3GOP_NEXUS3_RUT_EMAIL_HEADER",
	} {
		envVal := os.Getenv(envName)
		if envVal == "" {
			log.Fatalf("The %s variable must be set to a non-empty value", envName)
		}

		envConfig[envName] = envVal
	}
}
