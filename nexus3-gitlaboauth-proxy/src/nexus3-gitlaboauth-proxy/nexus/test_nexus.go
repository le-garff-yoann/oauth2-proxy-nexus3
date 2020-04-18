package nexus

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
)

// NewTestServer returns an `httptest.Server` that partially implements the Nexus 3 API.
func NewTestServer(userModifiers []UserModifier, roles *[]Role) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeWithDataCb := func(v interface{}) {
			payload, _ := json.Marshal(v)

			w.WriteHeader(http.StatusOK)
			w.Write(payload)
		}

		if r.URL.Path == userEndpointPath && r.Method == http.MethodGet {
			if userID := r.URL.Query().Get("userId"); userID != "" {
				for _, userModifier := range userModifiers {
					if userModifier.UserID == userID {
						writeWithDataCb([]User{userModifier.User})

						return
					}
				}
			}

			writeWithDataCb(userModifiers)

			return
		} else if r.URL.Path == userEndpointPath && r.Method == http.MethodPost {
			var data User
			if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
				w.WriteHeader(http.StatusBadRequest)

				return
			}

			userModifiers = append(userModifiers, UserModifier{
				User:   data,
				Source: defaultSourceName,
			})

			return
		} else if strings.HasPrefix(r.URL.Path, userEndpointPath) && r.Method == http.MethodPut {
			var data UserModifier
			if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
				w.WriteHeader(http.StatusBadRequest)

				return
			}

			userID := strings.TrimPrefix(r.URL.Path, userEndpointPath+"/")

			for i := range userModifiers {
				if userModifiers[i].UserID == userID {
					userModifiers[i] = data

					w.WriteHeader(http.StatusNoContent)

					return
				}
			}
		} else if r.URL.Path == roleEndpointPath {
			writeWithDataCb(&roles)

			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
}
