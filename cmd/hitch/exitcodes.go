package main

import (
	"go.polydawn.net/hitch/core"
	"go.polydawn.net/hitch/lib/errcat"
)

type ExitCode int

const (
	EXIT_SUCCESS       ExitCode = 0
	EXIT_BADARGS       ExitCode = 1 // Indicates usage errors.
	EXIT_PANIC         ExitCode = 2 // Placeholder.  We don't use this.  '2' happens when golang exits due to panic.
	EXIT_DATANOTFOUND  ExitCode = 4 // Indicates content not found -- the db was, but the e.g. catalog name (or release name, etc; whatever was requested) was not present.
	EXIT_CORRUPT       ExitCode = 6 // Indicates saved state is corrupt somehow (does not parse, or fails invariant checks).
	EXIT_WEIRDFS       ExitCode = 5 // Indicates some I/O error: permission denied, etc.
	EXIT_DBNOTFOUND    ExitCode = 6 // Returned when a hitch command is used outside of a hitch.db path.
	EXIT_INPROGRESS    ExitCode = 7 // Indicates desired operation is already begun -- e.g., `hitch init` is used and a hitch.db already exists; `hitch release start` when something is already staged, etc.
	EXIT_NAMECOLLISION ExitCode = 8 // Indicates that a name is already used, e.g. a catalog already has a release of the same name you just tried to create, etc.
)

var exitMapping = map[core.ErrorCategory]ExitCode{
	core.ErrBadArgs:       EXIT_BADARGS,
	core.ErrCorruptState:  EXIT_CORRUPT,
	core.ErrDataNotFound:  EXIT_DATANOTFOUND,
	core.ErrFS:            EXIT_WEIRDFS,
	core.ErrDBNotFound:    EXIT_DBNOTFOUND,
	core.ErrInProgress:    EXIT_INPROGRESS,
	core.ErrNameCollision: EXIT_NAMECOLLISION,
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
