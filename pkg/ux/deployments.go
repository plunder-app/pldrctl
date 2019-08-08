package ux

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/plunder-app/plunder/pkg/services"
)

// DeploymentsGetFormat -
func DeploymentsGetFormat(plunderConfig services.DeploymentConfigurationFile) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Mac Address\tDeploymemt\tHostname\tIP Address")
	for i := range plunderConfig.Configs {
		d := plunderConfig.Configs[i]

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", d.MAC, d.ConfigName, d.ConfigHost.ServerName, d.ConfigHost.IPAddress)
	}
	w.Flush()
}

// DeploymentDescribeBootFormat -
func DeploymentDescribeBootFormat(h services.DeploymentConfig, plunderURL string) {
	fmt.Printf("Boot description for deployment [%s]\n", h.ConfigHost.ServerName)
	fmt.Printf("-----------------------------------------------------------\n\n")
	fmt.Printf("Phase one\n------------------\n")
	fmt.Printf("DHCP Request -> TFTP Boot -> iPXE boot with config file -> http://%s/%s.ipxe\n", plunderURL, h.ConfigName)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "%s:\t%s\n", "Deployment", h.ConfigName)
	fmt.Fprintf(w, "%s:\t%s\n", "Kernel", h.ConfigBoot.Kernel)
	fmt.Fprintf(w, "%s:\t%s\n", "Initrd", h.ConfigBoot.Initrd)
	fmt.Fprintf(w, "%s:\t%s\n", "cmdline", h.ConfigBoot.Cmdline)

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
