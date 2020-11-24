package options

import (
	bbgit "github.com/craftamap/bb/git"
	"github.com/craftamap/bb/internal"
)

type GlobalOptions struct {
	BitbucketRepo *bbgit.BitbucketRepo
	Client        *internal.Client
}
