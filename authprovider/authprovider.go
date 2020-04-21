package authprovider

// Client is an Auth provider client.
type Client interface {
	GetUserInfo(accessToken string) (UserInfo, error)
}

// UserInfo is the partial representation
// of a OIDC */userinfo* response.
type UserInfo interface {
	Username() string
	EmailAddress() string
	Roles() []string
}
