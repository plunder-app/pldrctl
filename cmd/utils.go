package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/ghodss/yaml"
)

type resourceContainer struct {
	Definition string          `json:"definition"`
	Resource   json.RawMessage `json:"resource"`
}

func NewResourceContainer(containerDefinition string, resource json.RawMessage) *resourceContainer {
	if containerDefinition == "" {
		return nil
	}

	return &resourceContainer{
		Definition: containerDefinition,
		Resource:   resource,
	}
}

func UnPackResourceContainer(b []byte) (containerDefinition string, resource json.RawMessage, err error) {

	var container resourceContainer
	jsonBytes, err := yaml.YAMLToJSON(b)
	if err == nil {
		// If there were no errors then the YAML => JSON was successful, no attempt to unmarshall
		err = json.Unmarshal(jsonBytes, &container)
		if err != nil {
			return "", nil, fmt.Errorf("Unable to parse configuration as either yaml or json")
		}

	} else {
		// Couldn't parse the yaml to JSON
		// Attempt to parse it as JSON
		err = json.Unmarshal(b, &container)
		if err != nil {
			return "", nil, fmt.Errorf("Unable to parse configuration as either yaml or json")
		}
	}
	return container.Definition, container.Resource, nil
}
