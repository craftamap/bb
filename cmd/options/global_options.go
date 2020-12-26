package options

import (
	bbgit "github.com/craftamap/bb/git"
	"github.com/craftamap/bb/client"
)

type GlobalOptions struct {
	BitbucketRepo *bbgit.BitbucketRepo
	Client        *client.Client
}
