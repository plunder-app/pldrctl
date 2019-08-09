package cmd

import "github.com/spf13/cobra"

//pldrctlCreate - is used for it's subcommands for pulling data from a plunder server
var pldrctlCreate = &cobra.Command{
	Use:   "create",
	Short: "Apply a configuration to plunder",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
