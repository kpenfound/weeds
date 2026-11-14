package silenced

import "os"

func Remove(path string) {
	os.Remove(path) //nolint:errcheck
}

func Other(path string) {
	os.Remove(path) // nolint
}
