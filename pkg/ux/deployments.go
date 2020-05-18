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
	fmt.Fprintln(w, "Mac Address\tBoot Config Name\tHostname\tIP Address")
	for i := range plunderConfig.Configs {
		d := plunderConfig.Configs[i]

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", d.MAC, d.ConfigName, d.ConfigHost.ServerName, d.ConfigHost.IPAddress)
	}
	w.Flush()
}

// DeploymentDescribeBootFormat -
func DeploymentDescribeBootFormat(h services.DeploymentConfig, plunderURL, dashmac string) {
	// Boot Configuration
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	// Config overview
	fmt.Fprintf(w, "%s:\t\n", "Config")
	fmt.Fprintf(w, "\t%s:\t%s\n", "Boot Name", h.ConfigBoot.ConfigName)
	fmt.Fprintf(w, "\t%s:\t%s\n", "Boot Type", h.ConfigBoot.ConfigType)
	fmt.Fprintf(w, "\t%s:\t%s\n", "Kernel", h.ConfigBoot.Kernel)
	fmt.Fprintf(w, "\t%s:\t%s\n", "Initrd", h.ConfigBoot.Initrd)
	fmt.Fprintf(w, "\t%s:\t%s\n", "cmdline", h.ConfigBoot.Cmdline)
	fmt.Fprintf(w, "\t%s:\t%s\n", "Adapter", h.ConfigHost.Adapter)
	fmt.Fprintf(w, "\t%s:\t%s\n", "Server Name", h.ConfigHost.ServerName)
	fmt.Fprintf(w, "\t%s:\t%s\n", "IP Address", h.ConfigHost.IPAddress)

	// Boot Phases
	fmt.Fprintf(w, "%s:\t\n", "Phase One")
	fmt.Fprintf(w, "\t%s:\t%s\n", "Action", "DHCP Request -> TFTP Boot -> iPXE boot with config file")
	fmt.Fprintf(w, "\t%s:\thttp://%s/%s.ipxe\n", "Config", plunderURL, dashmac)
	// Further Boot Phases
	if h.ConfigBoot.ConfigType == "preseed" || h.ConfigBoot.ConfigType == "vsphere" {
		// Phase two boot
		fmt.Fprintf(w, "%s:\t\n", "Phase Two")
		fmt.Fprintf(w, "\t%s:\t%s\n", "Action", "OS Bootstraps with config")
		fmt.Fprintf(w, "\t%s:\thttp://%s/%s.cfg\n", "Config", plunderURL, dashmac)
		if h.ConfigBoot.ConfigType == "vsphere" {
			fmt.Fprintf(w, "%s:\t\n", "Phase Three")
			fmt.Fprintf(w, "\t%s:\t%s\n", "Action", "vSphere installer with config")
			fmt.Fprintf(w, "\t%s:\thttp://%s/%s.ks\n", "Config", plunderURL, dashmac)
		}
	}
	if h.ConfigBoot.ConfigType == "booty" {
		fmt.Fprintf(w, "%s:\t\n", "Phase Two")
		fmt.Fprintf(w, "\t%s:\t%s\n", "Action", "BOOTy the OS installer will begin")
		fmt.Fprintf(w, "\t%s:\thttp://%s/%s.bty\n", "Config", plunderURL, dashmac)
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
