package oui

import (
	"bufio"
	"strings"
)

var hotLookup map[string]string

func init() {
	// Initilise the quick lookup map
	hotLookup = make(map[string]string)
}

func lookupMacVendor(mac string) (string, error) {
	if vendor, ok := hotLookup[mac]; !ok {
		// Couldn't find it in the cache, look in the "cold" pool
		scanner := bufio.NewScanner(strings.NewReader(ouiLookup))
		for scanner.Scan() {
			//lines := append(lines, scanner.Text())
			return vendor, nil
		}
		err := scanner.Err()
		return "", err
	}
	return "", nil

}
