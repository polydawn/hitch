package db

import (
	"bytes"
	stdjson "encoding/json"
	"os"
	"path/filepath"
	"strings"

	. "github.com/polydawn/go-errcat"
	"github.com/polydawn/refmt/json"

	"go.polydawn.net/hitch/api"
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
func (dbctrl *Controller) LoadCatalog(catalogName api.CatalogName) (api.Catalog, error) {
	catalogNameChunks := strings.Split(string(catalogName), "/")
	catalogPath := filepath.Join(append([]string{dbctrl.BasePath}, catalogNameChunks...)...)

	f, err := os.OpenFile(filepath.Join(catalogPath, "catalog.json"), os.O_RDONLY, 0)
	if err != nil {
		if os.IsNotExist(err) {
			return api.Catalog{}, Errorf(ErrNotFound, "no catalog named %q in db", catalogName)
		}
		return api.Catalog{}, Recategorize(ErrIO, err)
	}
	defer f.Close()

	var catalog api.Catalog
	err = json.NewUnmarshallerAtlased(f, api.Atlas).
		Unmarshal(&catalog)
	return catalog, Recategorize(ErrStorageCorrupt, err)
}

/*
	Save a catalog object into the db.

	If the extended fields (e.g. replay) are set, they will also be saved;
	if absent, those components on disk will not be modified.
*/
func (dbctrl *Controller) SaveCatalog(catalog api.Catalog) error {
	catalogNameChunks := strings.Split(string(catalog.Name), "/")
	catalogPath := filepath.Join(append([]string{dbctrl.BasePath}, catalogNameChunks...)...)

	if err := os.MkdirAll(catalogPath, 0755); err != nil {
		return Recategorize(ErrIO, err)
	}
	f, err := os.OpenFile(filepath.Join(catalogPath, "catalog.json"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return Recategorize(ErrIO, err)
	}
	defer f.Close()

	msg, err := json.MarshalAtlased(catalog, api.Atlas)
	if err != nil {
		panic(err) // marshalling into a buffer shouldn't fail!
	}
	var buf bytes.Buffer
	stdjson.Indent(&buf, msg, "", "\t")
	buf.WriteString("\n")
	_, err = buf.WriteTo(f)
	return Recategorize(ErrIO, err)
}
