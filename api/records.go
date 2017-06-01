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

type CommissionName string
type Commission interface { /* it's a formula without pinned inputs */
}

type Replay struct {
	Computations map[CommissionName]Commission // review: this could still be `map[StepName]rdef.SetupHash`, because reminder, we're not a planner; this is replay not pipeline.
	Products     map[ItemLabel]struct {        // n.b. since this is final outputs only, implied here is commissions can also use an input "wire:<commissionName>:<outputPath>" for intermediates.
		CommissionName CommissionName
		OutputSlot     rdef.AbsPath
	}
	RunRecords map[CommissionName]map[rdef.RunRecordHash]*rdef.RunRecord // runRecords are stored indexed per step.  It is forbidden to have two runrecords for the same step with conflicting outputs IFF that output slot is referenced by either the final Products or any intermediate Wire.
}

var example = Replay{
	Computations: map[CommissionName]Commission{
		"prepare-step": map[string]interface{}{
			"inputs": map[rdef.AbsPath]interface{}{
				// remember, this is already thoroughly resolved: these infos are for your recursive lookup clarity.
				"/":         "hub.repeatr.io/base:2017-05-01:linux-amd64",
				"/task/src": "team.net/theproj:2.1.1:src",
			},
			"action": nil, // ... some preprocessor step, whatever ...
			"outputs": map[rdef.AbsPath]interface{}{
				"/task/output/docs": "tar",
				"/task/output/src":  "tar",
			},
		},
		"build-linux": map[string]interface{}{
			"inputs": map[rdef.AbsPath]interface{}{
				"/":            "hub.repeatr.io/base:2017-05-01:linux-amd64",
				"/app/compilr": "hub.repeatr.io/compilr:1.8:linux-amd64",
				"/task/src":    "wire:prepare-step:/task/output/src",
			},
			"action": nil, // ... some compiler is invoked ...
			"outputs": map[rdef.AbsPath]interface{}{
				"/task/output": "tar",
				"/task/logs":   "tar", // this is a byproduct (implicit: no products point at it).
			},
		},
	},
	Products: map[ItemLabel]struct {
		CommissionName CommissionName
		OutputSlot     rdef.AbsPath
	}{
		"src":         {"prepare-step", "/task/output/src"},
		"docs":        {"prepare-step", "/task/output/docs"},
		"linux-amd64": {"build-linux", "/task/output"},
	},
	RunRecords: map[CommissionName]map[rdef.RunRecordHash]*rdef.RunRecord{
		"prepare-step": map[rdef.RunRecordHash]*rdef.RunRecord{
			"aasdfasdf": &rdef.RunRecord{
				UID:       "234852-23792",
				Time:      23495,
				FormulaID: "oeiru43t", // Matches what happens when you put the inputs named in the commission together... but not currently otherwise explicitly mentioned.
				Results: map[rdef.AbsPath]rdef.WareID{
					"/task/output/docs": "tar:387ty874yt",
					"/task/output/src":  "tar:egruihieur",
				},
			},
		},
		// ... same for the other commissions.
	},
	// Implicitly now -- For this release record:
	//   - Items["src"] = "tar:egruihieur" -- this much match; correct doc format verifies this: the item label matches the products map, points to the commision name, has a runrecord, which has this wareID.
	//   - ... and the other items similarly.
}
