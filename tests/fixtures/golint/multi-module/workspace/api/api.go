package api

import "os"

func Remove(path string) {
	os.Remove(path)
}

func Unreachable() {
	return
	println("never")
}

func unusedHelper() int {
	return 2
}

func Pending() int {
	// TODO: compute this for real
	return 0
}
