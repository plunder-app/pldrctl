package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"text/tabwriter"

	"github.com/plunder-app/plunder/pkg/apiserver"
	"github.com/plunder-app/plunder/pkg/services"
	log "github.com/sirupsen/logrus"
)

func processURL(urlString, endpoint string) (*url.URL, error) {
	// Check that an address was actually entered
	if urlString == "" {
		return nil, fmt.Errorf("No Plunder server Address has been submitted")
	}

	// Check that the URL can be parsed
	u, err := url.Parse(urlString)
	if err != nil {
		return nil, fmt.Errorf("URL can't be parsed [%s]", err.Error())
	}

	u.Path = path.Join(u.Path, endpoint)

	return u, nil
}

func parseDeployments(u *url.URL) error {
	log.Debugf("Querying the Deployment Server [%s]", u.String())

	// Create a CA certificate pool and add cert.pem to it
	caCert, err := ioutil.ReadFile("plunder.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create a HTTPS client and supply the created CA pool
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}

	resp, err := client.Get(u.String())
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode > 200 {
		return fmt.Errorf(resp.Status)
	}

	var response apiserver.Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	if response.Error != "" {
		return fmt.Errorf("%s", response.Error)
	}

	var deployments services.DeploymentConfigurationFile

	err = json.Unmarshal(response.Payload, &deployments)
	if err != nil {
		return err
	}

	//test := response.Payload.()
	// if len(cfg.FriendlyError) == 0 {
	// 	log.Warnln("No deployment configurations found")
	// }

	// a, err := json.MarshalIndent(response.Payload, "", "\t")
	// if err != nil {
	// 	return err
	// }
	//fmt.Printf("%s\n%s\n", a, test)
	formatOutput(deployments)
	return nil
}

func formatOutput(plunderConfig services.DeploymentConfigurationFile) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Mac Address\tDeploymemt\tAllocated IP")
	for i := range plunderConfig.Configs {
		d := plunderConfig.Configs[i]

		fmt.Fprintf(w, "%s\t%s\t%s\n", d.MAC, d.ConfigName, d.ConfigHost.IPAddress)
	}
	w.Flush()

}
