package core

import (
	. "github.com/warpfork/go-errcat"

	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/hitch/core/db"
	"go.polydawn.net/hitch/core/stage"
	"go.polydawn.net/hitch/lib/locator"
)

func ReleaseAddItem(ui UI, itemNameStr, wareStr string) error {
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

	// Additional args parsery.
	// The "wareRef" param may be a "wire" type reference to a step (which is much more complicated),
	// or, just a regular "{type}:{hash}" WareID.
	// TODO : deal with this, or perhaps split the wire mode into a different subcommand for clarity.
	itemName := api.ItemName(itemNameStr)
	wareID, err := api.ParseWareID(wareStr)
	if err != nil {
		return Errorf(ErrBadArgs, "invalid ware reference -- %s", err)
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

	// Insert the item.  Then tell the stage state to save itself.
	items := stageCtrl.Catalog.Releases[0].Items
	if items == nil {
		items = make(map[api.ItemName]api.WareID)
	}
	items[itemName] = wareID
	stageCtrl.Catalog.Releases[0].Items = items
	if err := stageCtrl.Save(); err != nil {
		return Errorf(ErrFS, "error while saving staged state -- %s", err)
	}

	return nil
}
