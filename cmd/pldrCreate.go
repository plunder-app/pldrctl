package cmd

import (
	"encoding/json"
	"path"
	"strings"

	"github.com/plunder-app/plunder/pkg/apiserver"
	"github.com/plunder-app/plunder/pkg/services"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var createTypeFlag string
var deployment services.DeploymentConfig

func init() {
	pldrctlCreate.Flags().StringVarP(&createTypeFlag, "type", "t", "", "Type of resource to create")
	pldrctlCreate.Flags().StringVarP(&deployment.MAC, "mac", "m", "", "Mac Address of the resource to create")
	pldrctlCreate.Flags().StringVarP(&deployment.ConfigName, "config", "c", "", "The config to apply to the new resource")
	pldrctlCreate.Flags().StringVarP(&deployment.ConfigHost.IPAddress, "address", "a", "", "A Static address to apply to the new resource")
	pldrctlCreate.Flags().StringVarP(&deployment.ConfigHost.ServerName, "serverName", "n", "", "The hostname to apply to the new resource")

}

//pldrctlCreate - is used for it's subcommands for pulling data from a plunder server
var pldrctlCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a new configuration for plunder",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		u, c, err := apiserver.BuildEnvironmentFromConfig(pathFlag, urlFlag)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}

		switch strings.ToLower(createTypeFlag) {
		case "boot":
		case "config":
		case "deployment":
			// Call external function (TODO)
			u.Path = path.Join(u.Path, apiserver.DeploymentAPIPath())
			b, err := json.Marshal(deployment)
			if err != nil {
				log.Fatalf("%s", err.Error())
			}
			response, err := apiserver.ParsePlunderPost(u, c, b)
			if err != nil {
				log.Fatalf("%s", err.Error())
			}
			// If an error has been returned then handle the error gracefully and terminate
			if response.FriendlyError != "" || response.Error != "" {
				log.Debugln(response.Error)
				log.Fatalln(response.FriendlyError)
			}

		case "deployments":
		case "globalConfig":
		default:
			cmd.Help()
			log.Fatalf("Unknown resource Definition [%s]", createTypeFlag)

		}
	},
}
