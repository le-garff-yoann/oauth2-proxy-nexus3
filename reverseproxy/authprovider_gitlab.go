// +build authprovider_gitlab

package reverseproxy

import (
	"net/url"
	"oauth2-proxy-nexus3/authprovider"
	"oauth2-proxy-nexus3/authprovider/gitlab"
)

func newAuthproviderClient(URL *url.URL) authprovider.Client {
	return &gitlab.Client{URL: URL}
}
