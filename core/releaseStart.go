package core

import (
	"os"
	"path/filepath"

	. "github.com/polydawn/go-errcat"

	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/hitch/core/db"
	"go.polydawn.net/hitch/core/stage"
	"go.polydawn.net/hitch/lib/locator"
)

func ReleaseStart(ui UI, catalogNameStr, releaseNameStr string) error {
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

	// Check for staging file.  Reject command if staging file already exists.
	// Default location is a dir at "$root/_stage/".
	_, err = os.Stat(filepath.Join(dbctrl.BasePath, stage.DefaultPath))
	switch {
	case err == nil:
		return Errorf(ErrInProgress, "a release is already in progress!\nif this doesn't sound right, use 'hitch release reset' to discard the information (or, 'rm -r _stage').")
	case os.IsNotExist(err):
		// pass!
	default:
		return Errorf(ErrFS, "error while reading staged state -- %s", err)
	}

	// Validate names fit within acceptable string ranges.
	// TODO : write some regexps for this and do real checks.
	catalogName := api.CatalogName(catalogNameStr)
	releaseName := api.ReleaseName(releaseNameStr)

	// Check for catalog already existing.  Reject if not (this is fat-finger avoidance).
	catalog, err := dbctrl.LoadCatalog(catalogName)
	switch Category(err) {
	case nil:
		// pass!
	case db.ErrNotFound:
		return Errorf(ErrDataNotFound, "a catalog must exist before you can commit a release to it! -- %s", err)
	case db.ErrIO:
		return Errorf(ErrFS, "error while reading db -- %s", err)
	case db.ErrStorageCorrupt:
		return Errorf(ErrCorruptState, "error while reading db -- %s", err)
	default:
		panic(err)
	}

	// Check for catalog+release already existing.  Reject if released before.
	for _, release := range catalog.Releases {
		if release.Name == releaseName {
			return Errorf(ErrNameCollision, "\"%s:%s\" is already a catalogued release!  releases must have a unique name.", catalog.Name, release.Name)
		}
	}

	// If catalog has signing keys set up, check that we have those keys.
	// FUTURE : signing keys are in the roadmap, but a fair ways off.

	// All checks passed!
	// Initialize stage state on disk.
	_, err = stage.Create(dbctrl, stage.DefaultPath, catalogName, releaseName)
	if err != nil {
		return Errorf(ErrFS, "error while initializing stage state -- %s", err)
	}

	return nil
}
