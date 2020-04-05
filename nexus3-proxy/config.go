package main

import (
	"log"
	"os"
)

var envConfig = make(map[string]string)

func setConfig() {
	for _, envName := range []string{
		"N3P_NEXUS3_UPSTREAM", "N3P_NEXUS3_ADMIN_USER", "N3P_NEXUS3_ADMIN_PASSWORD",
		"N3P_NEXUS3_RUT_USER_HEADER", "N3P_NEXUS3_RUT_EMAIL_HEADER",
	} {
		envVal := os.Getenv(envName)
		if envVal == "" {
			log.Fatalf("The %s variable must be set to a non-empty value", envName)
		}

		envConfig[envName] = envVal
	}
}
