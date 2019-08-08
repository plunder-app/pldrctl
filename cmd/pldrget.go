package cmd

import (
	"encoding/json"
	"path"

	"github.com/plunder-app/pldrctl/pkg/plunderapi"
	"github.com/plunder-app/pldrctl/pkg/ux"

	"github.com/plunder-app/plunder/pkg/apiserver"
	"github.com/plunder-app/plunder/pkg/services"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {

	GetPlunderCmd.AddCommand(getDeployments)
	GetPlunderCmd.AddCommand(getGlobal)
	GetPlunderCmd.AddCommand(getConfig)

}

var getDeployments = &cobra.Command{
	Use:   "deployments",
	Short: "Retrieve all deployments from a Plunder server",
	Run: func(cmd *cobra.Command, args []string) {
		// Parse through the flags and attempt to build a correct URL
		log.SetLevel(log.Level(logLevel))

		u, c, err := plunderapi.BuildEnvironmentFromConfig(pathFlag, urlFlag)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}

		u.Path = path.Join(u.Path, apiserver.DeploymentsAPIPath())

		response, err := plunderapi.ParsePlunderGet(u, c)

		// If an error has been returned then handle the error gracefully and terminate
		if response.FriendlyError != "" || response.Error != "" {

		}

		var deployments services.DeploymentConfigurationFile

		err = json.Unmarshal(response.Payload, &deployments)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}

		ux.DeploymentFormat(deployments)
	},
}

var getGlobal = &cobra.Command{
	Use:   "global",
	Short: "Retrieve the global configuration from a Plunder server",
	Run: func(cmd *cobra.Command, args []string) {
		// Parse through the flags and attempt to build a correct URL
		log.SetLevel(log.Level(logLevel))

		u, c, err := plunderapi.BuildEnvironmentFromConfig(pathFlag, urlFlag)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}

		u.Path = path.Join(u.Path, apiserver.DeploymentsAPIPath())

		response, err := plunderapi.ParsePlunderGet(u, c)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
		// If an error has been returned then handle the error gracefully and terminate
		if response.FriendlyError != "" || response.Error != "" {
			log.Fatalf("%s", err.Error())

		}
		var deployments services.DeploymentConfigurationFile

		err = json.Unmarshal(response.Payload, &deployments)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}

		ux.GlobalFormat(deployments.GlobalServerConfig)
	},
}

var getConfig = &cobra.Command{
	Use:   "config",
	Short: "Retrieve the Plunder server configuration from a Plunder instance",
	Run: func(cmd *cobra.Command, args []string) {
		// Parse through the flags and attempt to build a correct URL
		log.SetLevel(log.Level(logLevel))

		u, c, err := plunderapi.BuildEnvironmentFromConfig(pathFlag, urlFlag)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}

		u.Path = path.Join(u.Path, apiserver.ConfigAPIPath())

		response, err := plunderapi.ParsePlunderGet(u, c)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
		// If an error has been returned then handle the error gracefully and terminate
		if response.FriendlyError != "" || response.Error != "" {
			log.Fatalf("%s", err.Error())

		}
		var serverConfig services.BootController

		err = json.Unmarshal(response.Payload, &serverConfig)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}

		ux.ServerFormat(serverConfig)
		//ux.GlobalFormat(deployments.GlobalServerConfig)

	},
}
