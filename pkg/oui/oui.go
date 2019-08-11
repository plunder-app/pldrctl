package oui

import (
	"bufio"
	"strings"

	log "github.com/sirupsen/logrus"
)

var hotLookup map[string]string

func init() {
	// Initilise the quick lookup map
	hotLookup = make(map[string]string)
}

// LookupMacVendor - will attempt to find a vendor in the hot cache or will search the large string
func LookupMacVendor(mac string) string {

	// Modify Mac address to use dashes (as per the oui.txt)
	dashMac := strings.Replace(mac, ":", "-", -1)

	// chop full MAC address to identifier (also upper case)
	macIdentifier := strings.ToUpper(string(dashMac[0:8]))

	// Search
	if vendor, ok := hotLookup[macIdentifier]; !ok {
		// Couldn't find it in the cache, look in the "cold" pool
		log.Debugf("Mac address [%s] not found in cache", macIdentifier)
		scanner := bufio.NewScanner(strings.NewReader(ouiLookup))
		for scanner.Scan() {
			// Look line-by-line
			FoundMAC := strings.Contains(scanner.Text(), macIdentifier)
			if FoundMAC {
				// Remove the prefix and the whitespaces
				vendor := strings.TrimSpace(strings.TrimPrefix(scanner.Text(), macIdentifier))
				// Add to cache
				hotLookup[macIdentifier] = vendor
				// Return the cold version
				return vendor
			}
		}
	} else {
		// Return the cached version
		return vendor
	}
	return "" // Return a blank

}

// splint and add to cache

//Return vendor
