package internal

import (
	"github.com/ktrysmt/go-bitbucket"
	"github.com/mitchellh/mapstructure"
)

type Members struct {
	Values []Account `mapstructure:"values"`
}

func (c Client) GetWorkspaceMembers(workspace string) (*Members, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)

	response, err := client.Teams.Members(workspace)
	if err != nil {
		return nil, err
	}
	members := Members{}
	mapstructure.Decode(response, &members)
	return &members, nil
}
