package core

import (
	"fmt"
	"os"
	"path/filepath"

	"go.polydawn.net/hitch/core/db"
	"go.polydawn.net/hitch/core/stage"
	"go.polydawn.net/hitch/lib/locator"
)

func ReleaseStart(ui UI, catalogName, releaseName string) ExitCode {
	// Find hitch.db root.
	dbctrl, err := db.LoadByCwd()
	switch err.(type) {
	case nil:
		// pass!
	case *locator.ErrNotFound:
		fmt.Fprintf(ui.Stderr, "no hitch.db found -- try `hitch init` first?\n")
		return EXIT_DBNOTFOUND
	default:
		fmt.Fprintf(ui.Stderr, "error while searching for hitch.db -- %s\n", err)
		return EXIT_WEIRDFS
	}

	// Check for staging file.  Reject command if staging file already exists.
	// Default location is a dir at "$root/_stage/".
	_, err = os.Stat(filepath.Join(dbctrl.BasePath, stage.DefaultPath))
	switch {
	case err == nil:
		fmt.Fprintf(ui.Stderr, "a release is already in progress!\nif this doesn't sound right, use 'hitch release reset' to discard the information (or, 'rm -r _stage').\n")
		return EXIT_INPROGRESS
	case os.IsNotExist(err):
		// pass!
	default:
		fmt.Fprintf(ui.Stderr, "error while reading staged state -- %s\n", err)
		return EXIT_WEIRDFS
	}

	// Check for catalog already existing.  Reject if not (this is fat-finger avoidance).
	// TODO : come back here after writing more db inspection methods.

	// Check for catalog+release already existing.  Reject if released before.
	// TODO : come back here after writing more db inspection methods.

	// If catalog has signing keys set up, check that we have those keys.
	// FUTURE : signing keys are in the roadmap, but a fair ways off.

	// All checks passed!
	// Initialize stage state on disk.
	_, err = stage.Create(dbctrl, stage.DefaultPath)
	if err != nil {
		fmt.Fprintf(ui.Stderr, "error while initializing stage state -- %s\n", err)
		return EXIT_WEIRDFS
	}

	return EXIT_SUCCESS
}
