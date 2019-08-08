package plunderapi

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/plunder-app/plunder/pkg/apiserver"
	log "github.com/sirupsen/logrus"
)

//BuildEnvironmentFromConfig will use the apiserver pkg to parse a configuration file and create a http client with the correct authentication and URL
func BuildEnvironmentFromConfig(path, urlFlag string) (*url.URL, *http.Client, error) {
	log.Debugf("Parsing Configuration file [%s]", path)

	// Open the configuration
	c, err := apiserver.OpenClientConfig(path)
	if err != nil {
		return nil, nil, err
	}
	// Retrieve the certificate
	cert, err := c.RetrieveClientCert()
	if err != nil {
		return nil, nil, err
	}

	// Build the certificate pool from the unencrypted cert
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cert)

	// Create a HTTPS client and supply the created CA pool
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}

	// Build the URL from the configuration
	serverURL := c.GetServerAddressURL()

	// Overwrite the configuration url if
	if urlFlag != "" {
		serverURL, err = url.Parse(urlFlag)
		if err != nil {
			return nil, nil, err
		}
	}

	return serverURL, client, nil
}

//ParsePlunderGet will attempt to retrieve data from the plunder API server
func ParsePlunderGet(u *url.URL, c *http.Client) (*apiserver.Response, error) {
	var response apiserver.Response

	log.Debugf("Querying the Plunder Server [%s]", u.String())

	resp, err := c.Get(u.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode > 200 {
		return nil, fmt.Errorf(resp.Status)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil

}
