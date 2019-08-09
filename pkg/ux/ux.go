package ux

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ghodss/yaml"
)

//CheckOutFlag - parses the flag and will attempt to parse the output
func CheckOutFlag(output string, o interface{}) error {
	switch strings.ToLower(output) {
	case "":
		return nil
	case "yaml":
		return OutputYAML(o)
	case "json":
		return OutputJSON(o)
	default:
		return fmt.Errorf("Unknown output type [%s]", output)
	}
}

// OutputJSON - takes an interface and writes out the JSON to stdout
func OutputJSON(o interface{}) error {
	// Marshall output interface to yaml
	b, err := json.MarshalIndent(o, "", "\t")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", b)
	return nil
}

// OutputYAML - takes an interface and writes out the JSON to stdout
func OutputYAML(o interface{}) error {
	// Marshall output interface to yaml
	b, err := yaml.Marshal(o)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", b)
	return nil
}
