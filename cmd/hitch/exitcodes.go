package main

import (
	"go.polydawn.net/hitch/core"
	"go.polydawn.net/hitch/lib/errcat"
)

type ExitCode int

const (
	EXIT_SUCCESS    ExitCode = 0
	EXIT_BADARGS    ExitCode = 1 // Indicates usage errors.
	EXIT_PANIC      ExitCode = 2 // Placeholder.  We don't use this.  '2' happens when golang exits due to panic.
	EXIT_CORRUPT    ExitCode = 6 // Indicates saved state is corrupt somehow (does not parse, or fails invariant checks).
	EXIT_WEIRDFS    ExitCode = 5 // Indicates some I/O error: permission denied, etc.
	EXIT_DBNOTFOUND ExitCode = 6 // Returned when a hitch command is used outside of a hitch.db path.
	EXIT_INPROGRESS ExitCode = 7 // Indicates desired operation is already begun -- e.g., `hitch init` is used and a hitch.db already exists; `hitch release start` when something is already staged, etc.
)

var exitMapping = map[core.ErrorCategory]ExitCode{
	core.ErrBadArgs:      EXIT_BADARGS,
	core.ErrCorruptState: EXIT_CORRUPT,
	core.ErrFS:           EXIT_WEIRDFS,
	core.ErrDBNotFound:   EXIT_DBNOTFOUND,
	core.ErrInProgress:   EXIT_INPROGRESS,
}

func mapToExitCode(err error) ExitCode {
	cat, ok := errcat.Category(err).(core.ErrorCategory)
	if !ok {
		return EXIT_PANIC
	}
	if v, ok := exitMapping[cat]; ok {
		return v
	}
	return EXIT_PANIC
}
