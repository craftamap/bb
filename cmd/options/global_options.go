package options

import (
	"github.com/craftamap/bb/client"
	bbgit "github.com/craftamap/bb/git"
)

type GlobalOptions struct {
	BitbucketRepo *bbgit.BitbucketRepo
	IsFSRepo      bool
	Client        *client.Client
}
