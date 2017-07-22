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
	if err != nil {
		return Errorf(ErrBadArgs, "malformed ID tuple: %s", err)
	}

	// Switch behavior based on specificity of args.
	switch {
	default:
		return showCatalog(dbctrl, tuple, ui.Stdout)
	case tuple.ReleaseName != "":
		return showRelease(dbctrl, tuple, ui.Stdout)
	case tuple.ItemName != "":
		return showItem(dbctrl, tuple, ui.Stdout)
	}
}

func showCatalog(dbctrl *db.Controller, tuple api.ReleaseItemID, w io.Writer) error {
	catalog, err := loadCatalog(dbctrl, tuple.CatalogName)
	if err != nil {
		return err
	}
	return emitPrettyJson(catalog, w)
}

func showRelease(dbctrl *db.Controller, tuple api.ReleaseItemID, w io.Writer) error {
	return nil
}

func showItem(dbctrl *db.Controller, tuple api.ReleaseItemID, w io.Writer) error {
	return nil
}
