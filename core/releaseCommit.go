package core

import (
	. "github.com/warpfork/go-errcat"

	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/hitch/core/db"
	"go.polydawn.net/hitch/core/stage"
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
	case db.ErrIO:
		return Errorf(ErrFS, "error while reading db -- %s", err)
	case db.ErrStorageCorrupt:
		return Errorf(ErrCorruptState, "error while reading db -- %s", err)
	default:
		panic(err)
	}

	// Do completion and sanity checks on the staged content.
	// You can exempt yourself from some of these by using additional command flags;
	//  the default is to halt and exit on any warnings at all.
	// TODO ...!

	// Merge the new release record into the existing catalog.
	// This process contains additional sanity checks; for example, colliding
	//  with an existing release name will be rejected.
	catalog, err = mergeCatalogs(catalog, stageCtrl.Catalog)
	switch Category(err) {
	case nil:
		// pass!
	case ErrNameCollision:
		return err
	default:
		panic(err)
	}

	// Save the updated catalog.  Delete the staged state.
	err = dbctrl.SaveCatalog(catalog)
	switch Category(err) {
	case nil:
		// pass!
	case db.ErrIO:
		return Errorf(ErrFS, "error while writing to db -- %s", err)
	default:
		panic(err)
	}
	stage.Clear(dbctrl, stage.DefaultPath)

	return nil
}

// this is only semantically correct for bringing in the half-catalog of a staged release.
// it started as intended to be a general merge, but it's not clear if it's actually
// gonna be reasonable to write things that way.
func mergeCatalogs(cat1, cat2 api.Catalog) (api.Catalog, error) {
	// Rack up all the new names, for efficient collision check.
	nNew := len(cat2.Releases)
	newNames := make(map[api.ReleaseName]struct{}, nNew)
	for _, release := range cat2.Releases {
		newNames[release.Name] = struct{}{}
	}
	// Allocate new array.  New stuff goes to top; old stuff to bottom.
	releases := make([]api.ReleaseEntry, len(cat1.Releases)+nNew)
	for i, release := range cat2.Releases {
		releases[i] = release
	}
	for i, release := range cat1.Releases {
		if _, exists := newNames[release.Name]; exists {
			return api.Catalog{}, Errorf(ErrNameCollision, "merge failed -- old and new catalogs both have a release named %q!", release.Name)
		}
		releases[nNew+i] = release
	}
	// Return cat1 with the updated releases slice stitched in.
	cat1.Releases = releases
	return cat1, nil
}
