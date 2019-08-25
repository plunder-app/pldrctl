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
func LeasesGetFormat(leases []services.Lease, noColour bool) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Mac Address\tHardware Vendor\tTime Seen\tTime since")
	for i := range leases {
		format := "%s\t%s\t%s\t%s\n"
		if !noColour {
			switch duration := time.Since(leases[i].Expiry).Minutes(); {
			case duration < 2:
				// Short duration typically, quickly rebooted
				format = "%s\t%s\t%s\t\x1b[0032m%s\x1b[0000m\n"

			case duration <= 30:
				// If we've not seen it for nearly ten minutes .. probably iLO
				format = "%s\t%s\t%s\t\x1b[0033m%s\x1b[0000m\n"

			case duration > 30:
				// It's not coming back at this point.
				format = "%s\t%s\t%s\t\x1b[0031m%s\x1b[0000m\n"
			}
		}
		// Build output template
		fmt.Fprintf(w, format,
			leases[i].Nic,                                       // Mac Address
			oui.LookupMacVendor(leases[i].Nic),                  // Vendor
			leases[i].Expiry.Format("Mon Jan _2 15:04:05 2006"), // Time Added
			time.Since(leases[i].Expiry).Truncate(time.Second))  // Time last seen
	}
	w.Flush()
}
