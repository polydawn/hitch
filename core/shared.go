package core

import (
	"bytes"
	stdjson "encoding/json"
	"io"

	. "github.com/polydawn/go-errcat"
	"github.com/polydawn/refmt/json"

	"go.polydawn.net/hitch/api"
	"go.polydawn.net/hitch/core/db"
)

// Bundles the "UI" types -- stdin/out/err.
// If a function has this as a parameter, it's a top-level command function.
type UI struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// Load the requested catalog from the db.
// All db pkg errors are mapped into core pkg errors;
// "not found" is treated as an error.
func loadCatalog(dbctrl *db.Controller, name api.CatalogName) (api.Catalog, error) {
	catalog, err := dbctrl.LoadCatalog(name)
	switch Category(err) {
	case nil:
		return catalog, nil
	case db.ErrNotFound:
		return api.Catalog{}, Errorf(ErrDataNotFound, "no catalog found named %q", name)
	case db.ErrIO:
		return api.Catalog{}, Errorf(ErrFS, "error while reading db -- %s", err)
	case db.ErrStorageCorrupt:
		return api.Catalog{}, Errorf(ErrCorruptState, "error while reading db -- %s", err)
	default:
		panic(err)
	}
}

// Serialize the 'thing' with the api.Atlas and write
// as pretty-printed json.
//
// Errors from pushing into the writer are returned as-is:
// caller should filter and flag them based on what job the writer
// was doing (either ErrFS or ErrPiping may be approripate).
// Serialization fails will be panicked: they are bugs.
func emitPrettyJson(thing interface{}, w io.Writer) error {
	msg, err := json.MarshalAtlased(thing, api.Atlas)
	if err != nil {
		panic(err) // marshalling into a buffer shouldn't fail!
	}
	var buf bytes.Buffer
	stdjson.Indent(&buf, msg, "", "\t")
	buf.WriteString("\n")
	_, err = buf.WriteTo(w)
	return err
}
