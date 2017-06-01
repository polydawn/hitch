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
	Computations map[CommissionName]Commission // review: this could still be `map[StepName]rdef.SetupHash`, because reminder, we're not a planner; this is replay not pipeline.
	Products     map[ItemLabel]struct {        // n.b. since this is final outputs only, implied here is commissions can also use an input "wire:<commissionName>:<outputPath>" for intermediates.
		CommissionName CommissionName
		OutputSlot     rdef.AbsPath
	}
	RunRecords map[CommissionName]map[rdef.RunRecordHash]*rdef.RunRecord // runRecords are stored indexed per step.  It is forbidden to have two runrecords for the same step with conflicting outputs IFF that output slot is referenced by either the final Products or any intermediate Wire.
}
