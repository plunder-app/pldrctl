package cmd

import (
	"strings"

	"github.com/plunder-app/plunder/pkg/apiserver"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var deleteTypeFlag string

func init() {
	pldrctlDelete.AddCommand(pldrctlDeleteDeployment)
	pldrctlDelete.AddCommand(pldrctlDeleteLogs)

	pldrctlDelete.Flags().StringVarP(&deleteTypeFlag, "type", "t", "", "Type of resource to create")
}

func deleteOperation(url string) (resp *apiserver.Response) {
	// Build the environment
	u, c, err := apiserver.BuildEnvironmentFromConfig(pathFlag, urlFlag)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	// Build the URL
	u.Path = url

	// Run the delete
	resp, err = apiserver.ParsePlunderDelete(u, c)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	return
}

//pldrctlDelete - is used for it's subcommands for pulling data from a plunder server
var pldrctlDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete a resource in Plunder",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		cmd.Help()
	},
}

//pldrctlDeleteDeployment - is used for it's subcommands for pulling data from a plunder server
var pldrctlDeleteDeployment = &cobra.Command{
	Use:   "deployment",
	Short: "Delete a a deployment in Plunder",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		if len(args) != 1 {
			log.Fatalf("Only argument should be a MAC address to be removed")
		}
		resp := deleteOperation(apiserver.DeploymentAPIPath() + "/" + strings.Replace(args[0], ":", "-", -1))
		if resp.FriendlyError != "" || resp.Error != "" {
			log.Debugln(resp.Error)
			log.Fatalln(resp.FriendlyError)
		}
	},
}

//pldrctlDeleteLogs - is used for it's subcommands for pulling data from a plunder server
var pldrctlDeleteLogs = &cobra.Command{
	Use:   "logs",
	Short: "Delete logs from a deployment in Plunder",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		if len(args) != 1 {
			log.Fatalf("Only argument should be an IP address to have it's logs removed")
		}
		resp := deleteOperation(apiserver.ParlayAPIPath() + "/logs/" + strings.Replace(args[0], ":", "-", -1))
		if resp.FriendlyError != "" || resp.Error != "" {
			log.Debugln(resp.Error)
			log.Fatalln(resp.FriendlyError)
		}
	},
}
