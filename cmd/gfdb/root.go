package gfdb

import (
	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize()

	RootCmd.AddCommand(
		GenStructCmd(),
		GenResultCmd(),
		GenGoaTypeCmd(),
		GenGoaResultTypeCmd(),
	)
}

var RootCmd = &cobra.Command{
	Use:   "gfdb",
	Short: "command line generate from database",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
