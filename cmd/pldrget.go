package cmd

import (
	"encoding/json"
	"os"
	"path"

	"github.com/plunder-app/pldrctl/pkg/plunderapi"
	"github.com/plunder-app/pldrctl/pkg/ux"

	"github.com/plunder-app/plunder/pkg/apiserver"
	"github.com/plunder-app/plunder/pkg/services"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	pldrctlGet.PersistentFlags().StringVar(&urlFlag, "url", os.Getenv("pURL"), "The Url of a plunder server")

	pldrctlGet.AddCommand(getDeployments)
	pldrctlGet.AddCommand(getGlobal)
	pldrctlGet.AddCommand(getConfig)
	pldrctlGet.AddCommand(getUnLeased)
}

var pldrctlGet = &cobra.Command{
	Use:   "get",
	Short: "Retrieve data from a Plunder server",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
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
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
		// If an error has been returned then handle the error gracefully and terminate
		if response.FriendlyError != "" || response.Error != "" {
			log.Debugln(response.Error)
			log.Fatalln(response.FriendlyError)
		}

		var deployments services.DeploymentConfigurationFile

		err = json.Unmarshal(response.Payload, &deployments)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}

		if outputFlag != "" {
			err = ux.CheckOutFlag(outputFlag, NewResourceContainer("deployments", response.Payload))
		} else {
			ux.DeploymentsGetFormat(deployments)
		}
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
			log.Debugln(response.Error)
			log.Fatalln(response.FriendlyError)
		}
		var deployments services.DeploymentConfigurationFile

		err = json.Unmarshal(response.Payload, &deployments)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}

		globalConfigJSON, _ := json.Marshal(deployments.GlobalServerConfig)

		if outputFlag != "" {
			err = ux.CheckOutFlag(outputFlag, NewResourceContainer("globalConfig", globalConfigJSON))
		} else {
			ux.GlobalFormat(deployments.GlobalServerConfig)
		}
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
			log.Debugln(response.Error)
			log.Fatalln(response.FriendlyError)
		}
		var serverConfig services.BootController

		err = json.Unmarshal(response.Payload, &serverConfig)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
		if outputFlag != "" {
			err = ux.CheckOutFlag(outputFlag, NewResourceContainer("serverConfig", response.Payload))
		} else {
			ux.ServerFormat(serverConfig)
		}
	},
}

var getUnLeased = &cobra.Command{
	Use:   "unleased",
	Short: "Retrieve the addresses that Plunder hasn't allocated",
	Run: func(cmd *cobra.Command, args []string) {
		// Parse through the flags and attempt to build a correct URL
		log.SetLevel(log.Level(logLevel))

		u, c, err := plunderapi.BuildEnvironmentFromConfig(pathFlag, urlFlag)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}

		u.Path = path.Join(u.Path, apiserver.DHCPAPIPath()+"/unleased")

		response, err := plunderapi.ParsePlunderGet(u, c)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
		// If an error has been returned then handle the error gracefully and terminate
		if response.FriendlyError != "" || response.Error != "" {
			log.Fatalf("%s", err.Error())

		}
		var unleased []services.Lease

		err = json.Unmarshal(response.Payload, &unleased)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}

		ux.LeasesGetFormat(unleased)
	},
}
