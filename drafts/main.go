package main

func main() {}

type Catalog struct {
	Name     CatalogName
	Releases []Release
}

type CatalogName string // oft like "project.org/thing".  The first part of an identifying triple.
type ReleaseName string // oft like "1.8".  The second part of an identifying triple.
type ItemLabel string   // oft like "linux-amd64" or "docs".  The third part of an identifying triple.

type WareID string // oft like "git:f23ae1829" or "tar:WJL8or32vD".

type Release struct {
	Name     ReleaseName
	Items    map[ItemLabel]WareID
	Metadata map[string]string
	Hazards  map[string]string
	Replay   *Replay
}

type SetupHash string // HID of formula

type AbsPath string // Identifier for output slots.  Coincidentally, a path.

type RunRecord struct {
	Time     int64
	Hostname string
	Results  map[AbsPath]WareID
}

type Replay struct {
	// This is a bit odd looking, but here we are:
	// RunRecords include a timestamp and host ID,
	// so they unique themselves; each points back
	// to a formula by SetupHash.
	// Several RunRecords may point back to the same SetupHash.
	// There is no unique reason for them (though you
	// may intuit that multiple records pointing to the
	// same setup hash indicates the releaser was checking
	// for reproducibility!).
	//
	// TODO review that -- if the executor was checking
	// for reproducibility, shouldn't we be able to state that
	// clearly herein?
	// Maybe... maybe not: we wanted to keep planning and releasing
	// predicates *out* of our serial format specifically so
	// we can describe them in turing complete tools and
	// run them in other containers rather than needing to bless
	// a special few rules that we've preordained.
	RunRecords map[*RunRecord]SetupHash
}

type CellName string      // human name
type RunRecordHash string // HID

type Replay2 struct {
	// A named pointer to a ware ID, effectively.
	// Other formulas in the replay info bundle can refer to it as a tag for inputs;
	// the CatalogName is always "local".
	//
	// TODO / REVIEW -- this has ambiguity if two runrecords for the formula.
	// Is that just illegal in releases?
	// (Additional people can later append more runrecords of their attempts to
	// reproduce, of course, but they won't be linked explicitly from the release
	// record, nor of course signed by the releaser's private keys, so they won't
	// be a point of confusion herein.)
	Cells map[CellName]struct {
		FormulaID  SetupHash
		OutputSlot AbsPath
	}
}
