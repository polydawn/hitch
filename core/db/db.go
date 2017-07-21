package db

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/polydawn/refmt/json"

	"go.polydawn.net/hitch/api"
	. "go.polydawn.net/hitch/lib/errcat"
)

type Controller struct {
	BasePath string
}

/*
	Load a catalog object from the db.

	This only loads the core catalog data -- names, releases, items, metadata.
	The replay info is *null* -- the replay info is not loaded by default,
	because that stuff may be much larger than many callers need;
	use the other db methods to ask for it explicitly.
*/
func (dbctrl *Controller) LoadCatalog(catalogName api.CatalogName) (*api.Catalog, error) {
	catalogNameChunks := strings.Split(string(catalogName), "/")
	catalogPath := filepath.Join(append([]string{dbctrl.BasePath}, catalogNameChunks...)...)

	f, err := os.OpenFile(filepath.Join(catalogPath, "catalog.json"), os.O_RDONLY, 0)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, Errorf(ErrNotFound, "no catalog named %q in db", catalogName)
		}
		return nil, Errorw(ErrIO, err)
	}
	defer f.Close()

	var catalog api.Catalog
	err = json.NewUnmarshallerAtlased(f, api.Atlas).
		Unmarshal(&catalog)
	return &catalog, Errorw(ErrStorageCorrupt, err)
}
