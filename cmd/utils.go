package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"text/tabwriter"

	"github.com/plunder-app/plunder/pkg/server"
	log "github.com/sirupsen/logrus"
)

func processURL(urlString, username, password, endpoint string) (*url.URL, error) {
	// Check that an address was actually entered
	if urlString == "" {
		return nil, fmt.Errorf("No Plunder server Address has been submitted")
	}

	// Check that the URL can be parsed
	u, err := url.Parse(urlString)
	if err != nil {
		return nil, fmt.Errorf("URL can't be parsed [%s]", err.Error())
	}

	//TODO - thebsdbox this will need removing

	if disableauth {
		// Check if a username was entered
		if username == "" {
			// if no username does one exist as part of the url
			if u.User.Username() == "" {
				return nil, fmt.Errorf("No Username has been submitted")
			}
		} else {
			// A username was submitted update the url
			u.User = url.User(username)
		}

		if password == "" {
			_, set := u.User.Password()
			if set == false {
				return nil, fmt.Errorf("No Password has been submitted")
			}
		} else {
			u.User = url.UserPassword(u.User.Username(), password)
		}
	}
	u.Path = path.Join(u.Path, endpoint)

	return u, nil
}

func parseDeployments(u *url.URL) error {
	log.Infof("Querying the Deployment Server [%s]", u.String())

	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode > 200 {
		return fmt.Errorf(resp.Status)
	}

	var cfg server.DeploymentConfigurationFile
	err = json.Unmarshal(body, &cfg)
	if err != nil {
		return err
	}
	if len(cfg.Deployments) == 0 {
		log.Warnln("No deployment configurations found")
	}
	formatOutput(cfg)
	return nil
}

func formatOutput(plunderConfig server.DeploymentConfigurationFile) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Mac Address\tDeploymemt\tAllocated IP")
	for i := range plunderConfig.Deployments {
		d := plunderConfig.Deployments[i]

		fmt.Fprintf(w, "%s\t%s\t%s\n", d.MAC, d.Deployment, d.Config.IPAddress)
	}
	w.Flush()

}
