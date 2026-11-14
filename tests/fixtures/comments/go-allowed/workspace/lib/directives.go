//go:build linux
// +build linux

package lib

//go:generate stringer -type=Mode

//line generated.y:12
var Marker = 1

func helper() {
	//go:noinline
	_ = Marker
}
