package core

import (
	"go.polydawn.net/hitch/core/db"
	"go.polydawn.net/hitch/core/stage"
	. "go.polydawn.net/hitch/lib/errcat"
	"go.polydawn.net/hitch/lib/locator"
)

func ReleaseCommit(ui UI) error {
	// Find hitch.db root.
	dbctrl, err := db.LoadByCwd()
	switch err.(type) {
	case nil:
		// pass!
	case *locator.ErrNotFound:
		return Errorf(ErrDBNotFound, "no hitch.db found -- try `hitch init` first?")
	default:
		return Errorf(ErrFS, "error while searching for hitch.db -- %s", err)
	}

	// Load stage state.  Staging must have already been started by `hitch release start`.
	stageCtrl, err := stage.Load(dbctrl, stage.DefaultPath)
	switch Category(err) {
	case nil:
		// pass!
	case stage.ErrIO:
		return Errorf(ErrFS, "error while reading staged state -- %s", err)
	case stage.ErrStorageCorrupt:
		return Errorf(ErrCorruptState, "error while reading staged state -- %s", err)
	default:
		panic(err)
	}

	// Load the existing catalog from the db.
	//  (We do this before any more sanity checks on the staged stuff,
	//  because some of the sanity checks include looking at the existing content!)
	catalog, err := dbctrl.LoadCatalog(stageCtrl.Catalog.Name)
	switch Category(err) {
	case nil:
		// pass!
	case db.ErrNotFound:
		return Errorf(ErrCorruptState, "a catalog must exist before you can commit a release to it! -- %s", err)
		// FIXME : we're not enforcing that invariant at start time yet...
	case db.ErrIO:
		return Errorf(ErrFS, "error while reading db -- %s", err)
	case db.ErrStorageCorrupt:
		return Errorf(ErrCorruptState, "error while reading db -- %s", err)
	default:
		panic(err)
	}
	_ = catalog

	// Do completion and sanity checks.
	// You can exempt yourself from some of these by using additional command flags;
	//  the default is to halt and exit on any warnings at all.
	// TODO ...!

	// Merge the new release record into the existing catalog.
	// TODO

	// Save the updated catalog.  Delete the staged state.
	// TODO

	return nil
}
