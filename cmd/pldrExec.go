package cmd

import (
	"encoding/json"
	"path"

	"github.com/plunder-app/pldrctl/pkg/plunderapi"
	"github.com/plunder-app/plunder/pkg/apiserver"
	"github.com/plunder-app/plunder/pkg/parlay"
	"github.com/plunder-app/plunder/pkg/parlay/types"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var execHost, execCommand string

func init() {
	pldrctlExec.Flags().StringVarP(&execHost, "address", "a", "", "Host to submit the execution job too")
	pldrctlExec.Flags().StringVarP(&execCommand, "command", "c", "", "Command to submit to the remote host")

}

//pldrcltlDescribe - is used for it's subcommands for pulling data from a plunder server
var pldrctlExec = &cobra.Command{
	Use:   "exec",
	Short: "Execute a command against a host",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		if execHost == "" {
			cmd.Help()
			log.Fatalf("No Host was submitted")
		}
		if execCommand == "" {
			cmd.Help()
			log.Fatalf("No Command was submitted")
		}

		newDeployment := parlay.Deployment{
			Name:     "pldrctl exec",
			Parallel: false,
		}
		newDeployment.Hosts = append(newDeployment.Hosts, execHost)

		action := types.Action{
			ActionType: "command",
			Command:    execCommand,
			Name:       "pldrctl command",
		}
		newDeployment.Actions = append(newDeployment.Actions, action)

		var newMap parlay.TreasureMap
		newMap.Deployments = append(newMap.Deployments, newDeployment)

		// Pass the execution data to the API endpoint
		u, c, err := plunderapi.BuildEnvironmentFromConfig(pathFlag, urlFlag)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}

		u.Path = path.Join(u.Path, apiserver.ParlayAPIPath())
		b, err := json.Marshal(newMap)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
		response, err := plunderapi.ParsePlunderPost(u, c, b)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
		// If an error has been returned then handle the error gracefully and terminate
		if response.FriendlyError != "" || response.Error != "" {
			log.Debugln(response.Error)
			log.Fatalln(response.FriendlyError)
		}

	},
}
