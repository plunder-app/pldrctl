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

func init() {

	pldrcltCmd.PersistentFlags().IntVar(&logLevel, "logLevel", 5, "Set the logging level [0=panic, 3=warning, 5=debug]")

	pldrcltCmd.AddCommand(pldrcltVersion)
	log.SetLevel(log.Level(logLevel))
}

var pldrcltCmd = &cobra.Command{
	Use:   "pldrctl",
	Short: "VMware vCenter Text User Interface",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		return
	},
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
	} else {
		// Default to logging anything Info and below
		logLevel = int(log.InfoLevel)
	}

	if err := pldrcltCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var pldrcltVersion = &cobra.Command{
	Use:   "version",
	Short: "Version and Release information about the plunder tool",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Plunder Release Information\n")
		fmt.Printf("Version:  %s\n", Release.Version)
		fmt.Printf("Build:    %s\n", Release.Build)
	},
}
