package cmd

import (
	"github.com/plunder-app/plunder/pkg/apiserver"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	GetPlunderCmd.AddCommand(getAll)
	GetPlunderCmd.AddCommand(getGlobal)

}

var getAll = &cobra.Command{
	Use:   "all",
	Short: "Retrieve all data from a Plunder server",
	Run: func(cmd *cobra.Command, args []string) {
		// Parse through the flags and attempt to build a correct URL
		log.SetLevel(log.Level(logLevel))

		var err error
		PlunderServer.URL, err = processURL(urlString, apiserver.DeploymentsAPIPath())
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
		err = parseDeployments(PlunderServer.URL)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
	},
}

var getGlobal = &cobra.Command{
	Use:   "global",
	Short: "Retrieve the global configuration from a Plunder server",
	Run: func(cmd *cobra.Command, args []string) {
		// Parse through the flags and attempt to build a correct URL
		log.SetLevel(log.Level(logLevel))

		var err error
		PlunderServer.URL, err = processURL(urlString, apiserver.ConfigAPIPath())
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
		err = parseDeployments(PlunderServer.URL)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
	},
}

var getConfig = &cobra.Command{
	Use:   "config",
	Short: "Retrieve the Plunder server configuration from a Plunder instance",
	Run: func(cmd *cobra.Command, args []string) {
		// Parse through the flags and attempt to build a correct URL
		log.SetLevel(log.Level(logLevel))

		var err error
		PlunderServer.URL, err = processURL(urlString, apiserver.ConfigAPIPath())
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
		err = parseDeployments(PlunderServer.URL)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
	},
}
