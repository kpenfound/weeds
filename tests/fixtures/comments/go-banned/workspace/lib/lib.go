package lib

// TODO: remove the global.
var Global = 1

// Run starts work.
func Run() {
	// body comment
	x := 1
	_ = x
}

// helper is unexported.
func helper() {}

// func old() {
// return nil
// }
var Gone = 2

var Trailing = 3 // trailing note

// Detached is separated by a blank line.

func Detached() {}

// FIXME: broken
func Broken() {}

type Config struct {
	// Name is exported.
	Name string
	// secret is not.
	secret string
}
