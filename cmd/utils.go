package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/plunder-app/plunder/pkg/services"
)

func formatOutput(plunderConfig services.DeploymentConfigurationFile) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Mac Address\tDeploymemt\tHostname\tIP Address")
	for i := range plunderConfig.Configs {
		d := plunderConfig.Configs[i]

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", d.MAC, d.ConfigName, d.ConfigHost.ServerName, d.ConfigHost.IPAddress)
	}
	w.Flush()

}
