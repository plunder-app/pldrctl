package cmd

import (
	"encoding/json"
	"path"
	"strings"

	"github.com/plunder-app/pldrctl/pkg/ux"
	"github.com/plunder-app/plunder/pkg/apiserver"
	"github.com/plunder-app/plunder/pkg/services"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	pldrctlDescribe.AddCommand(describeDeploymentBootProcess)
	pldrctlDescribe.AddCommand(describeDeployment)

}

//pldrcltlDescribe - is used for it's subcommands for pulling data from a plunder server
var pldrctlDescribe = &cobra.Command{
	Use:   "describe",
	Short: "Describe provides the capability to inspect or look deeper into a plunder resource",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var describeDeploymentBootProcess = &cobra.Command{
	Use:   "boot",
	Short: "Describe the boot process for a deployment (through it's MAC address)",
	Run: func(cmd *cobra.Command, args []string) {
		// Parse through the flags and attempt to build a correct URL
		if len(args) != 1 {
			log.Fatalf("Only argument should be a MAC address to be described")
		}

		log.SetLevel(log.Level(logLevel))

		u, c, err := apiserver.BuildEnvironmentFromConfig(pathFlag, urlFlag)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
		dashMac := strings.Replace(args[0], ":", "-", -1)

		u.Path = path.Join(u.Path, apiserver.DeploymentAPIPath()+"/"+dashMac)

		response, err := apiserver.ParsePlunderGet(u, c)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
		// If an error has been returned then handle the error gracefully and terminate
		if response.FriendlyError != "" || response.Error != "" {
			log.Debugln(response.Error)
			log.Fatalln(response.FriendlyError)

		}
		var deployment services.DeploymentConfig

		err = json.Unmarshal(response.Payload, &deployment)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}

		ux.DeploymentDescribeBootFormat(deployment, u.Hostname(), dashMac)
	},
}

var describeDeployment = &cobra.Command{
	Use:   "deployment",
	Short: "Describe the details of a deployment (through it's MAC address)",
	Run: func(cmd *cobra.Command, args []string) {
		// Parse through the flags and attempt to build a correct URL
		if len(args) != 1 {
			log.Fatalf("Only argument should be a MAC address to be described")
		}

		log.SetLevel(log.Level(logLevel))

		u, c, err := apiserver.BuildEnvironmentFromConfig(pathFlag, urlFlag)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
		dashMac := strings.Replace(args[0], ":", "-", -1)

		u.Path = path.Join(u.Path, apiserver.DeploymentAPIPath()+"/"+dashMac)

		response, err := apiserver.ParsePlunderGet(u, c)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
		// If an error has been returned then handle the error gracefully and terminate
		if response.FriendlyError != "" || response.Error != "" {
			log.Debugln(response.Error)
			log.Fatalln(response.FriendlyError)

		}
		var deployment services.DeploymentConfig

		err = json.Unmarshal(response.Payload, &deployment)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
		if outputFlag != "" {
			err = ux.CheckOutFlag(outputFlag, NewResourceContainer("deployment", response.Payload))
		} else {
			ux.DeploymentDescribeBootFormat(deployment, u.Hostname(), dashMac)
		}
	},
}
