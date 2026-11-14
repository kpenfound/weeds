// Copyright 2024 Example Corp. All rights reserved.
// Use of this source code is governed by the MIT licence.

// Package lib holds the allowed comment forms.
package lib

import "fmt"

// Greeter greets people.
type Greeter struct {
	// Name is who to greet.
	Name   string
	secret string
}

type (
	// Mode selects how to greet.
	Mode   int
	hidden int
)

const (
	// Loud shouts.
	Loud Mode = iota
	quiet
)

// Version is the library version.
var Version = "1"

// Speaker says things.
type Speaker interface {
	// Speak says it.
	Speak() string
}

// Greet greets the person.
func (g Greeter) Greet() string {
	return fmt.Sprint("hello ", g.Name)
}

// New returns a Greeter.
//go:noinline
func New(name string) Greeter {
	return Greeter{Name: name}
}

// Open opens the thing.
//
//go:noinline
//line open.y:12
func Open() {}
