//go:build authprovider_gitlab
// +build authprovider_gitlab

package reverseproxy

import (
	"net/url"

	"github.com/le-garff-yoann/oauth2-proxy-nexus3/authprovider"
	"github.com/le-garff-yoann/oauth2-proxy-nexus3/authprovider/gitlab"
)

func newAuthproviderClient(URL *url.URL) authprovider.Client {
	return &gitlab.Client{URL: URL}
}
