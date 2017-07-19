/*
	Package 'stage' handles assembling a release.

	Making a new release often takes a series of hitch commands --
	this matches how making a release often requires *several*
	large computations -- so all the intermediate staged states
	are serializable to disk.
*/
package stage

import (
	"io"
	"os"
	"path/filepath"

	"go.polydawn.net/hitch/api"
	"go.polydawn.net/hitch/core/db"
)

const DefaultPath = "_stage"

type Controller struct {
	dbctrl    *db.Controller
	stagePath string

	Catalog api.Catalog // catalog struct, sync'd with file.  always must have exactly one release entry.
}

/*
	Create a new empty release staging state.  Makes a dir, and creates the sigil file.
*/
func Create(
	dbctrl *db.Controller, stagePath string,
	catalogName api.CatalogName, releaseName api.ReleaseName,
) (*Controller, error) {
	err := os.MkdirAll(filepath.Join(dbctrl.BasePath, stagePath), 0755)
	if err != nil {
		return nil, err
	}
	f, err := os.OpenFile(filepath.Join(dbctrl.BasePath, stagePath, "stage.json"), os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	stageCtrl := &Controller{
		dbctrl:    dbctrl,
		stagePath: stagePath,

		Catalog: api.Catalog{
			Name: catalogName,
			Releases: []api.ReleaseEntry{
				{Name: releaseName},
			},
		},
	}
	return stageCtrl, stageCtrl.flush(f)
}

func (stageCtrl *Controller) Save() error {
	f, err := os.OpenFile(filepath.Join(stageCtrl.dbctrl.BasePath, stageCtrl.stagePath, "stage.json"), os.O_WRONLY|os.O_TRUNC, 0)
	if err != nil {
		return err
	}
	defer f.Close()
	return stageCtrl.flush(f)
}

func (stageCtrl *Controller) flush(w io.Writer) error {
	return nil // TODO serialize
}

func Load(dbctrl *db.Controller, stagePath string) (*Controller, error) {
	f, err := os.OpenFile(filepath.Join(dbctrl.BasePath, stagePath, "stage.json"), os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	stageCtrl := &Controller{
		dbctrl:    dbctrl,
		stagePath: stagePath,
	}
	return stageCtrl, stageCtrl.load(f)
}

func (stageCtrl *Controller) load(r io.Reader) error {
	return nil // TODO deserialize
}
