package core

import "io"

// Bundles the "UI" types -- stdin/out/err.
// If a function has this as a parameter, it's a top-level command function.
type UI struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// Panic with this to cause an "exit" with the given code.
// It will be caught by the "main" method and returned as an int again,
// so tests can DTRT and get coverage of exit codes without forking.
type Exit struct {
	Code int
}
