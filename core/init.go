package core

import (
	"fmt"
	"os"

	"go.polydawn.net/hitch/core/db"
	"go.polydawn.net/hitch/lib/locator"
)

func Init() {
	// Check for existing hitch.db root anywhere above this.
	//  Reject if exists.  Nested repos would be silly.
	dbctrl, err := db.LoadByCwd()
	switch err.(type) {
	case nil:
		panic(fmt.Errorf("cannot init new hitch.db -- one already exists, rooted at %q!", dbctrl.BasePath))
	case *locator.ErrNotFound:
		// pass!
	default:
		panic(err)
	}

	// Make hitch.db sigil file in cwd.
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	_, err = db.Create(cwd)
	if err != nil {
		panic(err)
	}
}
