package core

import (
	"io"

	"go.polydawn.net/hitch/api"
	"go.polydawn.net/hitch/core/db"
	. "go.polydawn.net/hitch/lib/errcat"
	"go.polydawn.net/hitch/lib/locator"
)

func Show(ui UI, nameStr string) error {
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

	// Args parsery.
	// This one is more interesting than usual because the behavior of this
	//  command changes radically based on how precise of a request was made.
	tuple, err := api.ParseReleaseItemID(nameStr)

	// Load the requested catalog from the db.
	catalog, err := dbctrl.LoadCatalog(tuple.CatalogName)
	switch Category(err) {
	case nil:
		// pass!
	case db.ErrNotFound:
		return Errorf(ErrDataNotFound, "no catalog found named %q", tuple.CatalogName)
	case db.ErrIO:
		return Errorf(ErrFS, "error while reading db -- %s", err)
	case db.ErrStorageCorrupt:
		return Errorf(ErrCorruptState, "error while reading db -- %s", err)
	default:
		panic(err)
	}

	// Switch behavior based on specificity of args.
	_ = catalog
	switch {
	case tuple.ItemName != "":
		return showItem(ui.Stdout, tuple)
	case tuple.ReleaseName != "":
		return showRelease(ui.Stdout, tuple)
	default:
		return showCatalog(ui.Stdout, tuple)
	}
}

func showItem(w io.Writer, tuple api.ReleaseItemID) error {
	return nil
}

func showRelease(w io.Writer, tuple api.ReleaseItemID) error {
	return nil
}

func showCatalog(w io.Writer, tuple api.ReleaseItemID) error {
	return nil
}
