package api

import (
	"polydawn.net/hitch/api/rdef"
)

type (
	CatalogName string // oft like "project.org/thing".  The first part of an identifying triple.
	ReleaseName string // oft like "1.8".  The second part of an identifying triple.
	ItemLabel   string // oft like "linux-amd64" or "docs".  The third part of an identifying triple.
)

/*
	A Catalog is the accumulated releases for a particular piece of software.

	A Catalog indicates a single author.  When observing new releases and/or
	metadata updates in a Catalog over time, you should expect to see it signed
	by the same key.  (Signing is not currently a built-in operation of `hitch`,
	but may be added in future releases.)
*/
type Catalog struct {
	// Name of self.
	Name CatalogName

	// Ordered list of release entries.
	// Order not particularly important, though UIs generally display in this order.
	// Most recent entries are typically placed at the top (e.g. index zero).
	//
	// Each entry must have a unique ReleaseName in the scope of its Catalog.
	Releases []ReleaseEntry
}

type ReleaseEntry struct {
	Name     ReleaseName
	Items    map[ItemLabel]rdef.WareID
	Metadata map[string]string
	Hazards  map[string]string
	Replay   *Replay
}

type Replay struct {
	// TODO
}
