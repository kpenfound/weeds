package nesteddot

import "os"

// removeIt is not exported.
func removeIt(path string) {
	os.Remove(path) //nolint:errcheck
}
