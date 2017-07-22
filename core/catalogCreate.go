package core

import (
	"go.polydawn.net/hitch/api"
	"go.polydawn.net/hitch/core/db"
	. "go.polydawn.net/hitch/lib/errcat"
	"go.polydawn.net/hitch/lib/locator"
)

func CatalogCreate(ui UI, catalogNameStr string) error {
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

	// Validate catalog name.
	catalogName := api.CatalogName(catalogNameStr) // TODO more :)

	// Check for an existing catalog of this name in the db.
	// Duplicates are *extremely* invalid!
	_, err = dbctrl.LoadCatalog(catalogName)
	switch Category(err) {
	case nil:
		return Errorf(ErrNameCollision, "a catalog of that name already exists!")
	case db.ErrNotFound:
		// pass!
	case db.ErrIO:
		return Errorf(ErrFS, "error while reading db -- %s", err)
	case db.ErrStorageCorrupt:
		return Errorf(ErrCorruptState, "error while reading db -- %s", err)
	default:
		panic(err)
	}

	// Initialize the clean new catalog.
	// Command the DB to save it.
	catalog := api.Catalog{
		Name: catalogName,
	}
	err = dbctrl.SaveCatalog(catalog)
	switch Category(err) {
	case nil:
		// pass!
	case db.ErrIO:
		return Errorf(ErrFS, "error while writing to db -- %s", err)
	default:
		panic(err)
	}

	return nil
}
