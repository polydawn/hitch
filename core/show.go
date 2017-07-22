package core

import (
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
	thing, err := getThing(dbctrl, tuple)
	if err != nil {
		return err
	}
	return Errorw(ErrPiping, emitPrettyJson(thing, ui.Stdout))
}

// Yields one of: a Catalog, ReleaseEntry, or WareID -- switching behavior
// based on how much info is provided in the tuple argument.
func getThing(dbctrl *db.Controller, tuple api.ReleaseItemID) (interface{}, error) {
	catalog, err := loadCatalog(dbctrl, tuple.CatalogName)
	if err != nil {
		return nil, err
	}
	if tuple.ReleaseName == "" {
		return catalog, nil
	}
	release, exists := selectRelease(catalog.Releases, tuple.ReleaseName)
	if !exists {
		return nil, Errorf(ErrDataNotFound, "no release named %q in catalog %q", tuple.ReleaseName, tuple.CatalogName)
	}
	if tuple.ItemName == "" {
		return release, nil
	}
	wareID, exists := release.Items[tuple.ItemName]
	if !exists {
		return nil, Errorf(ErrDataNotFound, "no item labeled %q in release %q", tuple.ItemName, tuple.ReleaseName)
	}
	return wareID, nil
}

func selectRelease(releases []api.ReleaseEntry, name api.ReleaseName) (api.ReleaseEntry, bool) {
	// O(n) search :/
	// But since we store them linearly, this is pretty much the way of it.
	// Bagging them up into a map would pay the same cost, and we tend not to
	// need to run this select more than once in the whole life of a task, so it's moot.
	for _, release := range releases {
		if release.Name == name {
			return release, true
		}
	}
	return api.ReleaseEntry{}, false
}
