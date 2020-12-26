package client

import (
	"path/filepath"

	"github.com/ktrysmt/go-bitbucket"
	"github.com/mitchellh/mapstructure"
)

type Downloads struct {
	Values []Download `mapstructure:"values"`
}

type Download struct {
	Name      string          `mapstructure:"name"`
	Links     map[string]Link `mapstructure:"links"`
	Downloads int             `mapstructure:"downloads"`
	CreatedOn string          `mapstructure:"created_on"`
	User      Account         `mapstructure:"user"`
	Type      string          `mapstructure:"type"`
	Size      int             `mapstructure:"size"`
}

func (c Client) GetDownloads(repoOrga string, repoSlug string) (*Downloads, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)

	opt := &bitbucket.DownloadsOptions{
		Owner:    repoOrga,
		RepoSlug: repoSlug,
	}

	res, err := client.Repositories.Downloads.List(opt)
	if err != nil {
		return nil, err
	}

	downloads := &Downloads{}
	err = mapstructure.Decode(res, downloads)
	if err != nil {
		return nil, err
	}

	return downloads, nil
}

func (c Client) UploadDownload(repoOrga string, repoSlug string, fpath string) (*Download, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)

	opt := &bitbucket.DownloadsOptions{
		Owner:    repoOrga,
		RepoSlug: repoSlug,
		FilePath: fpath,
		FileName: filepath.Base(fpath),
	}

	res, err := client.Repositories.Downloads.Create(opt)

	if err != nil {
		return nil, err
	}

	download := &Download{}
	err = mapstructure.Decode(res, download)
	if err != nil {
		return nil, err
	}

	return download, nil
}
