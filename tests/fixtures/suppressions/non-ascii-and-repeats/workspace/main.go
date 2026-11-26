package main

// é//nolint after a non-ASCII letter still has a boundary
var a = 1 // éeslint-disable
var b = 2 // nolint is not a directive: a no-break space is not \s
var c = 3 //nolint // lint:ignore both on one line
var d = 4 // # NOQA and # noqa count once, with //nolint on one line
var e = 5 // eslint-disable-line still ends at a boundary
