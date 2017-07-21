/*
	Package 'stage' handles assembling a release.

	Making a new release often takes a series of hitch commands --
	this matches how making a release often requires *several*
	large computations -- so all the intermediate staged states
	are serializable to disk.
*/
package stage

import (
	"bytes"
	stdjson "encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/polydawn/refmt/json"

	"go.polydawn.net/hitch/api"
	"go.polydawn.net/hitch/core/db"
	. "go.polydawn.net/hitch/lib/errcat"
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
		return nil, Errorw(ErrIO, err)
	}
	f, err := os.OpenFile(filepath.Join(dbctrl.BasePath, stagePath, "stage.json"), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return nil, Errorw(ErrIO, err)
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
		return Errorw(ErrIO, err)
	}
	defer f.Close()
	return stageCtrl.flush(f)
}

func (stageCtrl *Controller) flush(w io.Writer) error {
	msg, err := json.MarshalAtlased(stageCtrl.Catalog, api.Atlas)
	if err != nil {
		panic(err) // marshalling into a buffer shouldn't fail!
	}
	var buf bytes.Buffer
	stdjson.Indent(&buf, msg, "", "\t")
	_, err = buf.WriteTo(w)
	return Errorw(ErrIO, err)
}

func Load(dbctrl *db.Controller, stagePath string) (*Controller, error) {
	f, err := os.OpenFile(filepath.Join(dbctrl.BasePath, stagePath, "stage.json"), os.O_RDONLY, 0)
	if err != nil {
		return nil, Errorw(ErrIO, err)
	}
	defer f.Close()

	stageCtrl := &Controller{
		dbctrl:    dbctrl,
		stagePath: stagePath,
	}
	return stageCtrl, stageCtrl.load(f)
}

func (stageCtrl *Controller) load(r io.Reader) error {
	err := json.NewUnmarshallerAtlased(r, api.Atlas).
		Unmarshal(&stageCtrl.Catalog)
	return Errorw(ErrStorageCorrupt, err)
}
