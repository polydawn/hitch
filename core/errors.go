package core

type ErrorCategory string

const (
	ErrBadArgs       ErrorCategory = "badargs"       // Indicates usage errors.
	ErrCorruptState  ErrorCategory = "corruptstate"  // Indicates saved state is corrupt somehow (does not parse, or fails invariant checks).
	ErrDataNotFound  ErrorCategory = "datanotfound"  // Indicates content not found -- the db was, but the e.g. catalog name (or release name, etc; whatever was requested) was not present.
	ErrFS            ErrorCategory = "fs"            // Indicates some local I/O error: permission denied, etc.  Retries unwise.
	ErrDBNotFound    ErrorCategory = "dbnotfound"    // Returned when a hitch command is used outside of a hitch.db path.
	ErrInProgress    ErrorCategory = "inprogress"    // Indicates desired operation is already begun -- e.g., `hitch init` is used and a hitch.db already exists; `hitch release start` when something is already staged, etc.
	ErrNameCollision ErrorCategory = "namecollision" // Indicates that a name is already used, e.g. a catalog already has a release of the same name you just tried to create, etc.
	ErrPiping        ErrorCategory = "pipecollapse"  // Returned when IO to the *user* (e.g. stdout) fails.  (*Not* when disk or db IO fails.)
)
