package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/plunder-app/plunder/pkg/apiserver"
	log "github.com/sirupsen/logrus"

	"github.com/ghodss/yaml"
)

// ResourceContainer is a container of a plunder resource, allowing it's definition to define it
type ResourceContainer struct {
	Definition string          `json:"definition"`
	Resource   json.RawMessage `json:"resource"`
}

// NewResourceContainer - will create a wrapper container for a particular pludner resource
func NewResourceContainer(containerDefinition string, resource json.RawMessage) *ResourceContainer {
	if containerDefinition == "" {
		return nil
	}

	return &ResourceContainer{
		Definition: containerDefinition,
		Resource:   resource,
	}
}

// UnPackResourceContainer will take raw data, determine if it's yaml or json and then return it as
func UnPackResourceContainer(b []byte) (container *ResourceContainer, err error) {

	jsonBytes, err := yaml.YAMLToJSON(b)
	if err == nil {
		// If there were no errors then the YAML => JSON was successful, no attempt to unmarshall
		err = json.Unmarshal(jsonBytes, &container)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse configuration as either yaml or json")
		}

	} else {
		// Couldn't parse the yaml to JSON
		// Attempt to parse it as JSON
		err = json.Unmarshal(b, &container)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse configuration as either yaml or json")
		}
	}
	return container, nil
}

func parseResponseError(r *apiserver.Response) {
	// If an error has been returned then handle the error gracefully and terminate
	if r.Warning != "" {
		log.Warnf(r.Warning)
		if r.Error != "" {
			log.Fatalf(r.Error)
		}
	}
}
