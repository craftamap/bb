package auth

import (
	"github.com/craftamap/bb/cmd/commands/auth/login"
	"github.com/craftamap/bb/cmd/options"
	"github.com/spf13/cobra"
)

func Add(rootCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage bb authentification state",
	}
	login.Add(authCmd, globalOpts)
	rootCmd.AddCommand(authCmd)
}
