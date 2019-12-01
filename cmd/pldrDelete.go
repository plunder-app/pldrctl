package cmd

import (
	"net/http"
	"path"
	"strings"

	"github.com/plunder-app/plunder/pkg/apiserver"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {

	pldrctlDelete.AddCommand(pldrctlDeleteBoot)
	pldrctlDelete.AddCommand(pldrctlDeleteDeployment)
	pldrctlDelete.AddCommand(pldrctlDeleteLogs)

}

// func deleteOperation(url string) (resp *apiserver.Response) {
// 	// Build the environment
// 	u, c, err := apiserver.BuildEnvironmentFromConfig(pathFlag, urlFlag)
// 	if err != nil {
// 		log.Fatalf("%s", err.Error())
// 	}

// 	// Build the URL
// 	u.Path = url

// 	// Run the delete
// 	resp, err = apiserver.ParsePlunderDelete(u, c)
// 	if err != nil {
// 		log.Fatalf("%s", err.Error())
// 	}

// 	return
// }

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

		u, c, err := apiserver.BuildEnvironmentFromConfig(pathFlag, urlFlag)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}

		ep, resp := apiserver.FindFunctionEndpoint(u, c, "deploymentID", http.MethodDelete)
		if resp.Error != "" {
			log.Debug(resp.Error)
			log.Fatalf(resp.FriendlyError)
		}

		u.Path = path.Join(u.Path, ep.Path+"/"+strings.Replace(args[0], ":", "-", -1))

		response, err := apiserver.ParsePlunderDelete(u, c)

		if response.FriendlyError != "" || response.Error != "" {
			log.Debugln(response.Error)
			log.Fatalln(response.FriendlyError)
		}
	},
}

//pldrctlDeleteLogs - is used for it's subcommands for pulling data from a plunder server
var pldrctlDeleteLogs = &cobra.Command{
	Use:   "logs",
	Short: "Delete Parlay logs from a deployment in Plunder",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		if len(args) != 1 {
			log.Fatalf("Only argument should be an IP address to have it's logs removed")
		}

		u, c, err := apiserver.BuildEnvironmentFromConfig(pathFlag, urlFlag)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}

		ep, resp := apiserver.FindFunctionEndpoint(u, c, "parlayLog", http.MethodDelete)
		if resp.Error != "" {
			log.Debug(resp.Error)
			log.Fatalf(resp.FriendlyError)
		}

		u.Path = path.Join(u.Path, ep.Path+"/"+strings.Replace(args[0], ":", "-", -1))

		response, err := apiserver.ParsePlunderDelete(u, c)

		if response.FriendlyError != "" || response.Error != "" {
			log.Debugln(response.Error)
			log.Fatalln(response.FriendlyError)
		}
	},
}

//pldrctlDeleteBoot - is used for it's subcommands for pulling data from a plunder server
var pldrctlDeleteBoot = &cobra.Command{
	Use:   "boot",
	Short: "Delete boot configuration from Plunder",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		if len(args) != 1 {
			log.Fatalf("Only argument should be an IP address to have it's logs removed")
		}
		u, c, err := apiserver.BuildEnvironmentFromConfig(pathFlag, urlFlag)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}

		ep, resp := apiserver.FindFunctionEndpoint(u, c, "configBoot", http.MethodDelete)
		if resp.Error != "" {
			log.Debug(resp.Error)
			log.Fatalf(resp.FriendlyError)
		}

		u.Path = path.Join(ep.Path + "/" + args[0])

		response, err := apiserver.ParsePlunderDelete(u, c)
		if response.FriendlyError != "" || response.Error != "" {
			log.Debugln(response.Error)
			log.Fatalln(response.FriendlyError)
		}
	},
}
