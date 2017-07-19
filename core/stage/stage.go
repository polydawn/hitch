/*
	Package 'stage' handles assembling a release.

	Making a new release often takes a series of hitch commands --
	this matches how making a release often requires *several*
	large computations -- so all the intermediate staged states
	are serializable to disk.
*/
package stage

import (
	"os"
	"path/filepath"

	"go.polydawn.net/hitch/core/db"
)

const DefaultPath = "_stage"

type Controller struct {
	dbctrl    *db.Controller
	stagePath string
}

/*
	Create a new empty release staging state.  Makes a dir, and creates the sigil file.
*/
func Create(dbctrl *db.Controller, stagePath string) (*Controller, error) {
	err := os.MkdirAll(filepath.Join(dbctrl.BasePath, stagePath), 0755)
	if err != nil {
		return nil, err
	}
	f, err := os.OpenFile(filepath.Join(dbctrl.BasePath, stagePath, "stage.json"), os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return &Controller{
		dbctrl:    dbctrl,
		stagePath: stagePath,
	}, nil
}
