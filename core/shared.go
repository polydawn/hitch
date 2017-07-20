package core

import "io"

// Bundles the "UI" types -- stdin/out/err.
// If a function has this as a parameter, it's a top-level command function.
type UI struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

type ExitCode int

const (
	EXIT_SUCCESS    ExitCode = 0
	EXIT_BADARGS             = 1 // Indicates usage errors.
	EXIT_PANIC               = 2 // Placeholder.  We don't use this.  '2' happens when golang exits due to panic.
	EXIT_CORRUPT             = 6 // Indicates saved state is corrupt somehow (does not parse, or fails invariant checks).
	EXIT_WEIRDFS             = 5 // Indicates some I/O error: permission denied, etc.
	EXIT_DBNOTFOUND          = 6 // Returned when a hitch command is used outside of a hitch.db path.
	EXIT_INPROGRESS          = 7 // Indicates desired operation is already begun -- e.g., `hitch init` is used and a hitch.db already exists; `hitch release start` when something is already staged, etc.
)
