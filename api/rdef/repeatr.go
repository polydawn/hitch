/*
	NOTE YE WELL: this is a placeholder package,
	wherein we're mirroring many types declared in repeatr.

	We're evolving them freely and independently for the moment, but
	the time will come when we have to make both projects line up again!
*/
package rdef

import (
	"strings"

	"github.com/polydawn/refmt/obj/atlas"
)

/*
	Ware IDs are content-addressable, cryptographic hashes which uniquely identify
	a "ware" -- a packed filesystem snapshot.

	Ware IDs are serialized as a string in two parts, separated by a colon --
	for example like "git:f23ae1829" or "tar:WJL8or32vD".
	The first part communicates which kind of packing system computed the hash,
	and the second part is the hash itself.
*/
type WareID struct {
	Type string
	Hash string
}

var WareID_AtlasEntry = atlas.BuildEntry(WareID{}).Transform().
	TransformMarshal(atlas.MakeMarshalTransformFunc(
		func(x WareID) (string, error) {
			return string(x.Type) + ":" + string(x.Hash), nil
		})).
	TransformUnmarshal(atlas.MakeUnmarshalTransformFunc(
		func(x string) (WareID, error) {
			ss := strings.Split(x, ":")
			return WareID{ss[0], ss[1]}, nil
		})).
	Complete()

type SetupHash string // HID of formula

type AbsPath string // Identifier for output slots.  Coincidentally, a path.

type RunRecord struct {
	UID       string             // random number, presumed globally unique.
	Time      int64              // time at start of build.
	FormulaID SetupHash          // HID of formula ran.
	Results   map[AbsPath]WareID // wares produced by the run!

	// --- below: addntl optional metadata ---

	Hostname string            // hostname.  not a trusted field, but useful for debugging.
	Metadata map[string]string // escape valve.  you can attach freetext here.
}

var RunRecord_AtlasEntry = atlas.BuildEntry(RunRecord{}).StructMap().Autogenerate().Complete()

type RunRecordHash string // HID of RunRecord.  Includes UID, etc, so quite unique.  Prefer this to UID for primary key in storage (it's collision resistant).

type Formula interface{} // TODO bother to finish fleshing this back out
