package cmd

import (
	"encoding/json"

	"github.com/plunder-app/plunder/pkg/apiserver"
	"github.com/plunder-app/plunder/pkg/services"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var createTypeFlag string
var deployment services.DeploymentConfig
var bootConfig services.BootConfig

func init() {
	pldrctlCreateBoot.Flags().StringVarP(&bootConfig.ConfigName, "name", "n", "demo", "The name of the new boot configuration")
	pldrctlCreateBoot.Flags().StringVarP(&bootConfig.ConfigType, "type", "t", "default", "The type of the new boot configuration")

	pldrctlCreateBoot.Flags().StringVarP(&bootConfig.Kernel, "kernel", "k", "", "The path of the kernel to be booted")
	pldrctlCreateBoot.Flags().StringVarP(&bootConfig.Initrd, "initrd", "i", "", "The path of init ramdisk to be booted")
	pldrctlCreateBoot.Flags().StringVarP(&bootConfig.Cmdline, "cmdline", "c", "", "Additional kernel commandline flags (optional)")
	pldrctlCreateBoot.Flags().StringVar(&bootConfig.ISOPath, "isoPath", "", "The path of an ISO to read from (optional)")
	pldrctlCreateBoot.Flags().StringVar(&bootConfig.ISOPrefix, "isoPrefix", "", "The prefix to bind the ISO too (optional)")

	pldrctlCreateDeployment.Flags().StringVarP(&createTypeFlag, "type", "t", "", "Type of resource to create")
	pldrctlCreateDeployment.Flags().StringVarP(&deployment.MAC, "mac", "m", "", "Mac Address of the resource to create")
	pldrctlCreateDeployment.Flags().StringVarP(&deployment.ConfigName, "config", "c", "", "The config to apply to the new resource")
	pldrctlCreateDeployment.Flags().StringVarP(&deployment.ConfigHost.IPAddress, "address", "a", "", "A Static address to apply to the new resource")
	pldrctlCreateDeployment.Flags().StringVarP(&deployment.ConfigHost.ServerName, "serverName", "n", "", "The hostname to apply to the new resource")

	pldrctlCreate.AddCommand(pldrctlCreateBoot)
	pldrctlCreate.AddCommand(pldrctlCreateDeployment)

}

//pldrctlCreate - is used for it's subcommands for pulling data from a plunder server
var pldrctlCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a new configuration for plunder",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		cmd.Help()
	},
}

//pldrctlCreateBoot - is used for it's subcommands for pulling data from a plunder server
var pldrctlCreateBoot = &cobra.Command{
	Use:   "boot",
	Short: "Create a new boot configuration for plunder",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		b, err := json.Marshal(bootConfig)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
		u, c, err := apiserver.BuildEnvironmentFromConfig(pathFlag, urlFlag)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
		ep, resp := apiserver.FindFunctionEndpoint(u, c, "configBoot", "POST")
		if resp.Error != "" {
			log.Debug(resp.Error)
			log.Fatalf(resp.Warning)
		}
		u.Path = ep.Path + "/" + bootConfig.ConfigName

		response, err := apiserver.ParsePlunderPost(u, c, b)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
		parseResponseError(response)

	},
}

//pldrctlCreate - is used for it's subcommands for pulling data from a plunder server
var pldrctlCreateDeployment = &cobra.Command{
	Use:   "deployment",
	Short: "Create a new deployment for plunder",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		if deployment.MAC == "" {
			cmd.Help()
			log.Fatalf("The MAC address is a required field")
		}

		if deployment.ConfigHost.IPAddress == "" {
			cmd.Help()
			log.Fatalf("The IP address is a required field")
		}

		b, err := json.Marshal(deployment)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
		// Build the environment
		u, c, err := apiserver.BuildEnvironmentFromConfig(pathFlag, urlFlag)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
		ep, resp := apiserver.FindFunctionEndpoint(u, c, "deployment", "POST")
		parseResponseError(resp)
		u.Path = ep.Path

		response, err := apiserver.ParsePlunderPost(u, c, b)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}

		parseResponseError(response)

	},
}
