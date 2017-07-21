package core

import "io"

// Bundles the "UI" types -- stdin/out/err.
// If a function has this as a parameter, it's a top-level command function.
type UI struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}
