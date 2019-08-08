package ux

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/plunder-app/plunder/pkg/services"
)

// DeploymentFormat -
func DeploymentFormat(plunderConfig services.DeploymentConfigurationFile) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Mac Address\tDeploymemt\tHostname\tIP Address")
	for i := range plunderConfig.Configs {
		d := plunderConfig.Configs[i]

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", d.MAC, d.ConfigName, d.ConfigHost.ServerName, d.ConfigHost.IPAddress)
	}
	w.Flush()

}

//GlobalFormat will display the global deployment configuration for a Plunder Server
func GlobalFormat(g services.HostConfig) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintf(w, "%s:\t%s\n", "Adapter", g.Adapter)
	fmt.Fprintf(w, "%s:\t%s\n", "Gateway", g.Gateway)
	fmt.Fprintf(w, "%s:\t%s\n", "Subnet", g.Subnet)
	fmt.Fprintf(w, "%s:\t%s\n", "NameServer", g.NameServer)
	fmt.Fprintf(w, "%s:\t%s\n", "TimeServer", g.NTPServer)
	fmt.Fprintf(w, "%s:\t%s\n", "Username", g.Username)
	fmt.Fprintf(w, "%s:\t%s\n", "Password", g.Password)
	fmt.Fprintf(w, "%s:\t%s\n", "SSH Key Path", g.SSHKeyPath)
	fmt.Fprintf(w, "%s:\t%s\n", "Repository", g.RepositoryAddress)
	fmt.Fprintf(w, "%s:\t%s\n", "Ubuntu URL", g.MirrorDirectory)

	fmt.Fprintf(w, "%s:\t%s\n", "Packages", g.Packages)

	w.Flush()

}
