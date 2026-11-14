// Copyright 2024 Example Corp.
// Licensed under the MIT licence.

package main

import _ "embed"

//go:embed data.txt
var data string

func main() {
	//go:noinline
	run()
}

func run() { _ = data }
