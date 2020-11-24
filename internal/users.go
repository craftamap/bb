package internal

import "github.com/ktrysmt/go-bitbucket"

func (c Client) GetCurrentUser() (*bitbucket.User, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)

	return client.User.Profile()
}
