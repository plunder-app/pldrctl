package ux

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/plunder-app/plunder/pkg/services"
)

//LeasesGetFormat -
func LeasesGetFormat(leases []services.Lease) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Mac Address\tTime Seen\tTime since")
	for i := range leases {
		l := leases[i]

		fmt.Fprintf(w, "%s\t%s\t%s\n", l.Nic, l.Expiry.String(), time.Since(l.Expiry))
	}
	w.Flush()
}
