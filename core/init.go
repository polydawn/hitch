package core

import (
	"os"

	. "github.com/warpfork/go-errcat"

	"go.polydawn.net/hitch/core/db"
	"go.polydawn.net/hitch/lib/locator"
)

func Init(ui UI) error {
	// Check for existing hitch.db root anywhere above this.
	//  Reject if exists.  Nested repos would be silly.
	dbctrl, err := db.LoadByCwd()
	switch err.(type) {
	case nil:
		return Errorf(ErrInProgress, "cannot init new hitch.db -- one already exists, rooted at %q!", dbctrl.BasePath)
	case *locator.ErrNotFound:
		// pass!
	default:
		return Errorf(ErrFS, "error while searching for hitch.db -- %s", err)
	}

	// Make hitch.db sigil file in cwd.
	cwd, err := os.Getwd()
	if err != nil {
		return Errorf(ErrFS, "error while creating hitch.db -- %s", err)
	}
	_, err = db.Create(cwd)
	if err != nil {
		return Errorf(ErrFS, "error while creating hitch.db -- %s", err)
	}
	return nil
}
