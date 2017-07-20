package core

import (
	"fmt"

	"go.polydawn.net/hitch/api"
	"go.polydawn.net/hitch/api/rdef"
	"go.polydawn.net/hitch/core/db"
	"go.polydawn.net/hitch/core/stage"
	"go.polydawn.net/hitch/lib/locator"
)

func ReleaseAddLabel(ui UI, labelNameStr, wareStr string) ExitCode {
	// Find hitch.db root.
	dbctrl, err := db.LoadByCwd()
	switch err.(type) {
	case nil:
		// pass!
	case *locator.ErrNotFound:
		fmt.Fprintf(ui.Stderr, "no hitch.db found -- try `hitch init` first?\n")
		return EXIT_DBNOTFOUND
	default:
		fmt.Fprintf(ui.Stderr, "error while searching for hitch.db -- %s\n", err)
		return EXIT_WEIRDFS
	}

	// Additional args parsery.
	// The "wareRef" param may be a "wire" type reference to a step (which is much more complicated),
	// or, just a regular "{type}:{hash}" WareID.
	// TODO : deal with this, or perhaps split the wire mode into a different subcommand for clarity.
	labelName := api.ItemLabel(labelNameStr)
	wareID, err := rdef.ParseWareID(wareStr)
	if err != nil {
		fmt.Fprintf(ui.Stderr, "invalid ware reference -- %s\n", err)
		return EXIT_BADARGS
	}

	// Load stage state.  Staging must have already been started by `hitch release start`.
	stageCtrl, err2 := stage.Load(dbctrl, stage.DefaultPath)
	switch {
	case err2 == nil:
		// pass!
	case err2.Category == stage.ErrIO:
		fmt.Fprintf(ui.Stderr, "error while reading staged state -- %s\n", err)
		return EXIT_WEIRDFS
	case err2.Category == stage.ErrStorageCorrupt:
		fmt.Fprintf(ui.Stderr, "error while reading staged state -- %s\n", err)
		return EXIT_CORRUPT
	default:
		panic(err2)
	}

	// Insert the label.  Then tell the stage state to save itself.
	items := stageCtrl.Catalog.Releases[0].Items
	if items == nil {
		items = make(map[api.ItemLabel]rdef.WareID)
	}
	items[labelName] = wareID
	stageCtrl.Catalog.Releases[0].Items = items
	if err := stageCtrl.Save(); err != nil {
		fmt.Fprintf(ui.Stderr, "error while saving staged state -- %s\n", err)
		return EXIT_WEIRDFS
	}

	return EXIT_SUCCESS
}
