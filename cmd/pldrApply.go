package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/plunder-app/plunder/pkg/apiserver"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var applyPathFlag string

func init() {
	pldrctlApply.Flags().StringVarP(&applyPathFlag, "file", "f", "", "Path of a Plunder configuration to be applied to a server")
}

//pldrctlApply - is used for it's subcommands for pulling data from a plunder server
var pldrctlApply = &cobra.Command{
	Use:   "apply",
	Short: "Apply a configuration to plunder",
	Long:  "This subcommand can also read from STDIN by passing \"-\" as a flag",
	Run: func(cmd *cobra.Command, args []string) {
		// Parse through the flags and attempt to build a correct URL
		log.SetLevel(log.Level(logLevel))
		var output []byte

		if len(args) == 1 && args[0] == "-" {
			reader := bufio.NewReader(os.Stdin)

			for {
				input, err := reader.ReadByte()
				if err != nil && err == io.EOF {
					break
				}
				output = append(output, input)
			}
		}
		if applyPathFlag != "" {
			// Check the actual path from the string
			if _, err := os.Stat(applyPathFlag); !os.IsNotExist(err) {
				output, err = ioutil.ReadFile(applyPathFlag)
				if err != nil {
					cmd.Help()

					log.Fatalf("%v", err)
				}
			} else {
				cmd.Help()

				log.Fatalf("Unable to open [%s]", applyPathFlag)
			}
		}
		if len(output) == 0 {
			log.Fatalf("No data could be read")
		}
		resource, err := UnPackResourceContainer(output)
		if err != nil {
			cmd.Help()
			log.Fatalln(err)
		}

		err = parseApply(resource.Definition, resource.Resource)
		if err != nil {
			cmd.Help()

			log.Fatalf("%v", err)
		}
	},
}

func parseApply(resourceDefinition string, resource json.RawMessage) error {
	u, c, err := apiserver.BuildEnvironmentFromConfig(pathFlag, urlFlag)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	//dashMac := strings.Replace(args[0], ":", "-", -1)

	switch strings.ToLower(resourceDefinition) {
	case "boot":
	case "config":
	case "deployment":

		ep, resp := apiserver.FindFunctionEndpoint(u, c, "deployment", "POST")
		if resp.Error != "" {
			log.Warnf(resp.Warning)
			log.Fatalf(resp.Error)

		}

		u.Path = path.Join(u.Path, ep.Path)

		response, err := apiserver.ParsePlunderPost(u, c, resource)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
		parseResponseError(response)

	case "deployments":
		// Set the url

		ep, resp := apiserver.FindFunctionEndpoint(u, c, "deployments", "POST")
		if resp.Error != "" {
			log.Warnf(resp.Warning)
			log.Fatalf(resp.Error)
		}

		u.Path = path.Join(u.Path, ep.Path)

		// Apply the POST
		response, err := apiserver.ParsePlunderPost(u, c, resource)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
		parseResponseError(response)

	case "globalconfig":
		// Set the url
		ep, resp := apiserver.FindFunctionEndpoint(u, c, "deployments", "POST")
		if resp.Error != "" {
			log.Warnf(resp.Warning)
			log.Fatalf(resp.Error)
		}

		u.Path = path.Join(u.Path, ep.Path+"/global")

		// Apply the POST
		response, err := apiserver.ParsePlunderPost(u, c, resource)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
		parseResponseError(response)

	case "parlay":
		ep, resp := apiserver.FindFunctionEndpoint(u, c, "parlay", "POST")
		if resp.Error != "" {
			log.Warnf(resp.Warning)
			log.Fatalf(resp.Error)
		}

		u.Path = path.Join(u.Path, ep.Path)

		// Apply the POST
		response, err := apiserver.ParsePlunderPost(u, c, resource)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
		parseResponseError(response)

	default:
		return fmt.Errorf("Unknown resource Definition [%s]", resourceDefinition)
	}
	return nil
}
