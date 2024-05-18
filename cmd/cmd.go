package cmd

import (
	"github.com/chalfel/chi-auth-0/cmd/api"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "app",
		Short: "Chi auth 0 app",
	}

	root.AddCommand(api.NewApiCmd())

	return root
}
