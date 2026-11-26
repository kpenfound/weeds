package dotignored

import "os"

// removeIt is not exported.
func removeIt(path string) {
	os.Remove(path) //nolint:errcheck
}
