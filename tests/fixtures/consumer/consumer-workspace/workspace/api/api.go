package api

import "os"

// Remove deletes the file.
func Remove(path string) {
	// the error does not matter here
	os.Remove(path)
}
