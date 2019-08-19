package cmd

import (
	"encoding/json"
	"os"
	"path"
	"strings"
	"time"

	"github.com/plunder-app/pldrctl/pkg/plunderapi"
	"github.com/plunder-app/pldrctl/pkg/ux"

	"github.com/plunder-app/plunder/pkg/apiserver"
	"github.com/plunder-app/plunder/pkg/plunderlogging"
	"github.com/plunder-app/plunder/pkg/services"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var watch int

func init() {
	pldrctlGet.PersistentFlags().StringVar(&urlFlag, "url", os.Getenv("pURL"), "The Url of a plunder server")

	getLogs.Flags().IntVarP(&watch, "watch", "w", 0, "Setting a watch timeout until \"completion\" or \"fail\"")

	pldrctlGet.AddCommand(getBoot)
	pldrctlGet.AddCommand(getDeployments)
	pldrctlGet.AddCommand(getGlobal)
	pldrctlGet.AddCommand(getConfig)
	pldrctlGet.AddCommand(getLogs)

	pldrctlGet.AddCommand(getUnLeased)
}

var pldrctlGet = &cobra.Command{
	Use:   "get",
	Short: "Retrieve data from a Plunder server",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var getBoot = &cobra.Command{
	Use:   "boot",
	Short: "Retrieve the Plunder server boot configuration from a Plunder instance",
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
			err = ux.CheckOutFlag(outputFlag, NewResourceContainer("bootConfig", response.Payload))
		} else {
			ux.BootFormat(serverConfig.BootConfigs)
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

var getLogs = &cobra.Command{
	Use:   "logs",
	Short: "Retrieve the logs from a Parlay deployment",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		u, c, err := plunderapi.BuildEnvironmentFromConfig(pathFlag, urlFlag)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
		if len(args) != 1 {
			log.Fatalf("An address of a remote server to retrieve the logs for is required")
		}

		dashAddress := strings.Replace(args[0], ".", "-", -1)

		u.Path = path.Join(u.Path, apiserver.ParlayAPIPath()+"/logs/"+dashAddress)
		for {
			response, err := plunderapi.ParsePlunderGet(u, c)
			if err != nil {
				log.Fatalf("%s", err.Error())
			}
			// If an error has been returned then handle the error gracefully and terminate
			if response.FriendlyError != "" || response.Error != "" {
				log.Debugln(response.Error)
				log.Fatalln(response.FriendlyError)
			}

			var logs plunderlogging.JSONLog

			err = json.Unmarshal(response.Payload, &logs)
			if err != nil {
				log.Fatalf("%s", err.Error())
			}

			if outputFlag != "" {
				err = ux.CheckOutFlag(outputFlag, NewResourceContainer("logs", response.Payload))
				break
			} else {
				if watch == 0 {
					ux.LogsGetFormat(logs)
					break
				} else {
					ux.LogsGetFormat(logs)
					if logs.State != "Running" {
						break
					}
					// Now repeat again
					time.Sleep(time.Duration(watch) * time.Second)
				}
			}
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
