package cmd

import (
	"fmt"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Release - this struct contains the release information populated when building katbox
var Release struct {
	Version string
	Build   string
}

var logLevel int
var urlFlag, pathFlag string

// Used to determine JSON/YAML output
var outputFlag string

func init() {
	pldrctlCmd.PersistentFlags().StringVarP(&outputFlag, "output", "o", "", "The Url of a plunder server")

	pldrctlCmd.PersistentFlags().IntVar(&logLevel, "logLevel", int(log.InfoLevel), "Set the logging level [0=panic, 3=warning, 5=debug]")
	pldrctlCmd.PersistentFlags().StringVarP(&pathFlag, "path", "p", "plunderclient.yaml", "Path to a custom Plunder Server configuation")

	// Add all of the subcommands to the pldrctl UX
	pldrctlCmd.AddCommand(pldrctlApply)
	pldrctlCmd.AddCommand(pldrctlCreate)
	pldrctlCmd.AddCommand(pldrctlDescribe)
	pldrctlCmd.AddCommand(pldrctlGet)
	pldrctlCmd.AddCommand(pldrctlVersion)
}

// Execute - starts the command parsing process
func Execute() {
	if os.Getenv("VCLOG") != "" {
		i, err := strconv.ParseInt(os.Getenv("VCLOG"), 10, 8)
		if err != nil {
			log.Fatalf("Error parsing environment variable [VCLOG")
		}
		// We've only parsed to an 8bit integer, however i is still a int64 so needs casting
		logLevel = int(i)
	}

	if err := pldrctlCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var pldrctlCmd = &cobra.Command{
	Use:   "pldrctl",
	Short: "Plunder CLI (command line interface)",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		cmd.Help()
		return
	},
}

var pldrctlVersion = &cobra.Command{
	Use:   "version",
	Short: "Version and Release information about the plunder tool",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Plunder Release Information\n")
		fmt.Printf("Version:  %s\n", Release.Version)
		fmt.Printf("Build:    %s\n", Release.Build)
	},
}
