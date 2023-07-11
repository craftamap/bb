package client

import (
	"github.com/ktrysmt/go-bitbucket"
	"github.com/mitchellh/mapstructure"
)

type Workspace struct {
	CreatedOn string                 `mapstructure:"created_on"`
	Links     map[string]interface{} `mapstructure:"links"`
	Name      string                 `mapstructure:"name"`
	Slug      string                 `mapstructure:"slug"`
	IsPrivate bool                   `mapstructure:"is_private"`
	Type      string                 `mapstructure:"type"`
	UUID      string                 `mapstructure:"uuid"`
}

type WorkspaceMembership struct {
	User      Account   `mapstructure:"user"`
	Workspace Workspace `mapstructure:"workspace"`
}

type Workspaces struct {
	Workspaces []Workspace `mapstructure:"workspaces"`
}

type Members struct {
	Values []WorkspaceMembership `mapstructure:"values"`
}

func (c Client) GetWorkspaceMembers(workspace string) (*Members, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)

	response, err := client.Workspaces.Members(workspace)
	if err != nil {
		return nil, err
	}
	members := Members{}
	err = mapstructure.Decode(response, &members)
	return &members, err
}

func (c Client) GetWorkspaces() (*Workspaces, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)

	response, err := client.Workspaces.List()
	if err != nil {
		return nil, err
	}

	var workspaces *Workspaces
	err = mapstructure.Decode(response, &workspaces)
	return workspaces, err
}
