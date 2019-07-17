package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	GetPlunderCmd.AddCommand(getAll)
}

var getAll = &cobra.Command{
	Use:   "all",
	Short: "Retrieve all data from a Plunder server",
	Run: func(cmd *cobra.Command, args []string) {
		// Parse through the flags and attempt to build a correct URL
		var err error

		PlunderServer.URL, err = processURL(urlString, username, password)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
	},
}
