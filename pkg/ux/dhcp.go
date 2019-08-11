package ux

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/plunder-app/pldrctl/pkg/oui"
	"github.com/plunder-app/plunder/pkg/services"
)

//LeasesGetFormat -
func LeasesGetFormat(leases []services.Lease) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Mac Address\tTime Seen\tTime since\tHardware Vendor")
	for i := range leases {

		// Build output template
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			leases[i].Nic,
			leases[i].Expiry.Format("Mon Jan _2 15:04:05 2006"),
			time.Since(leases[i].Expiry).Truncate(time.Second),
			oui.LookupMacVendor(leases[i].Nic))
	}
	w.Flush()
}
