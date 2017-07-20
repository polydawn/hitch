package core

import (
	"fmt"
	"os"

	"go.polydawn.net/hitch/core/db"
	"go.polydawn.net/hitch/lib/locator"
)

func Init(ui UI) ExitCode {
	// Check for existing hitch.db root anywhere above this.
	//  Reject if exists.  Nested repos would be silly.
	dbctrl, err := db.LoadByCwd()
	switch err.(type) {
	case nil:
		fmt.Fprintf(ui.Stderr, "cannot init new hitch.db -- one already exists, rooted at %q!\n", dbctrl.BasePath)
		return EXIT_INPROGRESS
	case *locator.ErrNotFound:
		// pass!
	default:
		fmt.Fprintf(ui.Stderr, "error while searching for hitch.db -- %s\n", err)
		return EXIT_WEIRDFS
	}

	// Make hitch.db sigil file in cwd.
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(ui.Stderr, "error while creating hitch.db -- %s\n", err)
		return EXIT_WEIRDFS
	}
	_, err = db.Create(cwd)
	if err != nil {
		fmt.Fprintf(ui.Stderr, "error while creating hitch.db -- %s\n", err)
		return EXIT_WEIRDFS
	}
	return EXIT_SUCCESS
}
