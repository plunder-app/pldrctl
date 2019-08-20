package cmd

import (
	"path"
	"strings"

	"github.com/plunder-app/plunder/pkg/apiserver"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var deleteTypeFlag string

func init() {
	pldrctlDelete.Flags().StringVarP(&deleteTypeFlag, "type", "t", "", "Type of resource to create")
}

//pldrctlCreate - is used for it's subcommands for pulling data from a plunder server
var pldrctlDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete a resource in Plunder",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		u, c, err := apiserver.BuildEnvironmentFromConfig(pathFlag, urlFlag)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
		switch strings.ToLower(deleteTypeFlag) {
		case "config":
		case "deployment":
			dashMac := strings.Replace(args[0], ":", "-", -1)

			u.Path = path.Join(u.Path, apiserver.DeploymentAPIPath()+"/"+dashMac)

			response, err := apiserver.ParsePlunderDelete(u, c)
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
		case "logs":
			dashAddress := strings.Replace(args[0], ":", "-", -1)

			u.Path = path.Join(u.Path, apiserver.ParlayAPIPath()+"/logs/"+dashAddress)

			response, err := apiserver.ParsePlunderDelete(u, c)
			if err != nil {
				log.Fatalf("%s", err.Error())
			}
			// If an error has been returned then handle the error gracefully and terminate
			if response.FriendlyError != "" || response.Error != "" {
				log.Debugln(response.Error)
				log.Fatalln(response.FriendlyError)
			}
		default:
			log.Fatalf("Unknown resource type [%s]", deleteTypeFlag)

		}
	},
}
