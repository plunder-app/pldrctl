package ux

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/plunder-app/plunder/pkg/services"
)

//BootFormat will display the global deployment configuration for a Plunder Server
func BootFormat(bootConfigs []services.BootConfig) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Config Name\tKernel Path\tInitrd Path\tCommand Line")
	for i := range bootConfigs {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			bootConfigs[i].ConfigName,
			bootConfigs[i].Kernel,
			bootConfigs[i].Initrd,
			bootConfigs[i].Cmdline)
	}
	w.Flush()

}

//ServerFormat will display the global deployment configuration for a Plunder Server
func ServerFormat(b services.BootController) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintf(w, "%s:\t%s\n", "Adapter", *b.AdapterName)
	fmt.Fprintf(w, "%s:\t%t\n", "Enable DHCP", *b.EnableDHCP)
	// DNS Configuration
	fmt.Fprintf(w, "\t%s:\t%s\n", "DHCP Start Address", *b.DHCPConfig.DHCPStartAddress)
	fmt.Fprintf(w, "\t%s:\t%s\n", "DHCP Server Address", *b.DHCPConfig.DHCPAddress)
	fmt.Fprintf(w, "\t%s:\t%s\n", "DHCP Gateway Address", *b.DHCPConfig.DHCPGateway)
	fmt.Fprintf(w, "\t%s:\t%s\n", "DHCP DNS Address", *b.DHCPConfig.DHCPDNS)
	fmt.Fprintf(w, "\t%s:\t%d\n", "DHCP Lease Pool Size", *b.DHCPConfig.DHCPLeasePool)
	fmt.Fprintf(w, "%s:\t%t\n", "Enable TFTP", *b.EnableTFTP)
	// TFTP Configuration
	fmt.Fprintf(w, "\t%s:\t%s\n", "TFTP Server Address", *b.TFTPAddress)

	fmt.Fprintf(w, "%s:\t%t\n", "Enable HTTP", *b.EnableHTTP)
	// TFTP Configuration
	fmt.Fprintf(w, "\t%s:\t%s\n", "DHCP Server Address", *b.HTTPAddress)
	fmt.Fprintf(w, "\t%s:\t%s\n", "PXE File Name", *b.PXEFileName)

	w.Flush()

}
