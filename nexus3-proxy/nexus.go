package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
)

type nexusUserInfo struct {
	Username, Email string
}

func createNexusUser(userInfo *nexusUserInfo) error {
	nexusUserEndpoint, err := url.Parse(fmt.Sprintf("%s/%s", envConfig["N3P_NEXUS3_UPSTREAM"], "service/rest/beta/security/users"))
	if err != nil {
		log.Fatal(err)
	}

	dataAsJSON, err := json.Marshal(map[string]interface{}{
		"userId":       userInfo.Username,
		"firstName":    userInfo.Username,
		"lastName":     userInfo.Username,
		"emailAddress": userInfo.Email,
		"password":     uuid.New().String(),
		"status":       "active",
		"roles":        []string{"nx-admin"}, // TODO: group:role mapping if the role exists.
	})
	if err != nil {
		log.Panic(err)
	}

	req, err := http.NewRequest("POST", nexusUserEndpoint.String(), bytes.NewBuffer(dataAsJSON))
	if err != nil {
		return err
	}

	req.SetBasicAuth(envConfig["N3P_NEXUS3_ADMIN_USER"], envConfig["N3P_NEXUS3_ADMIN_PASSWORD"])
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK { // && res.StatusCode != http.StatusConflict // It doesn't work... doesn't look like a RESTful API.
		if body, err := ioutil.ReadAll(res.Body); err == nil {
			if !strings.Contains(string(body), fmt.Sprintf("found duplicated key '%s'", userInfo.Username)) {
				return fmt.Errorf(`Cannot create the user "%s" on %s: %s`, userInfo.Username, envConfig["N3P_NEXUS3_UPSTREAM"], body)
			}
		} else {
			return err
		}
	}

	return nil
}
