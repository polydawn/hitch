package api

import (
	"polydawn.net/hitch/api/rdef"
)

/*
	Release names come in three parts:
	the Catalog Name, the Release Name, and the Item Label.

	The catalog is a whole group of releases for a single project.
	For example, there's a "repeatr" releases catalog.

	The release name is a string attached to the release;
	it's specified by the releaser when creating a new release.

	The item label is a selector used to select a specific ware when there is
	more than one ware published in a single atomic release.
	The set of items in a release is considered immutable once the release is published.
	Generally, there's an expectation the ecosystem that the set of item labels available
	from each release will be the same: e.g., when upgrading from an older version
	of repeatr, one might expect to jump from "repeatr.io/repeatr:1.0:linux-amd64"
	to "repeatr.io/repeatr:1.1:linux-amd64".
*/
type (
	CatalogName string // oft like "project.org/thing".  The first part of an identifying triple.
	ReleaseName string // oft like "1.8".  The second part of an identifying triple.
	ItemLabel   string // oft like "linux-amd64" or "docs".  The third part of an identifying triple.
)

type ReleaseItemID struct {
	CatalogName
	ReleaseName
	ItemLabel
}

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
	// The set of steps recorded in this replay.
	// Each step contains a formula with precise instructions on how to run the step again,
	// and additional data on where the inputs were selected from, so that
	// recursive audits can work automatically.
	Steps map[StepName]Step

	// Map wiring the ItemLabels in the release outputs to a step and output slot
	// within that step's formula.
	//
	// As with "wire" mode in a Step's Imports, if the referenced step has more than
	// one RunRecord, then the wired output slot MUST have resulted in the same WareID
	// hash all the RunRecords.
	// Furthermore, the WareID in those RunRecords must match the WareID which
	// is directly listed for the ItemLabel in the the release entry; otherwise,
	// the replay isn't describing the same thing released!
	Products map[ItemLabel]struct {
		StepName   StepName
		OutputSlot rdef.AbsPath
	}
}

type StepName string
type Step struct {
	// Record upstream names for formula inputs.
	//
	// Each key must match an input key in the formula or it is invalid.
	// The formula may have inputs that are not explained here (though tools
	// should usually emit a warning about such unexplained blobs).
	//
	// Imports may either be the full `{CatalogName,ReleaseName,ItemLabel}` tuple
	// referring to another catalog, or, `{"wire",StepName,OutputSlot}`.
	//
	// In the "wire" mode, the reference is interpreted as another step in this replay.
	// Hashes coming from a "wire" may be purely internal to the replay
	// (meaning, practically speaking, that ware may be an intermediate which
	// is not actually be stored anywhere).
	// If step referred to by a "wire" has more than one RunRecord, the wired
	// output slot MUST have resulted in the same WareID hash
	// all the RunRecords, or the replay is invalid.
	Imports map[rdef.AbsPath]ReleaseItemID

	// The formula for this step, exactly as executed by the releaser.
	//
	// This includes inputs (with full hashes), the script run,
	// and the output slots saved.
	// Names of inputs are separately stored; they're in the `Import` field.
	// Results are separately stored: they're in the `RunRecords` field.
	//
	// Note: it's entirely possible for two steps with different names in a Replay
	// to have identical formulas (and thus identical setupHashes).
	// In this case, both steps may also share identical RunRecords(!), if the
	// original releaser used a formula runner smart enough to notice this
	// and dedup the computation; the steps are still stored separately, because
	// it is correct to render them separately in order to represent the
	// releaser's original intentions clearly.
	Formula rdef.Formula

	// RunRecords from executions of this formula.
	// May be one or multiple.
	//
	// These are only the records included by the releaser at the time of release.
	// Other rebuilders may have more RunRecords to share, but these are stored
	// elsewhere (you may look them up by using the Formula.SetupHash).
	//
	// It is forbidden to have two RunRecords in the same step to declare
	// different resulting WareIDs for an output slot
	// if that output slot is referenced by either the final Products
	// or any intermediate "wire"-mode Import.
	// (It's perfectly fine for RunRecords to have differing results for outputs
	// that *aren't* so referenced; an output slot which captures logs, for example,
	// may often differ between runs, but since it's not passed forward, so be it.)
	RunRecords map[rdef.RunRecordHash]*rdef.RunRecord
}

var example = Replay{
	Steps: map[StepName]Step{
		"prepare-step": Step{
			Imports: map[rdef.AbsPath]ReleaseItemID{
				// Given the same snapshot of all relevant Catalogs, this section
				// should be sufficient to reproduce the Formula.Inputs.
				//
				// Remember, this is already thoroughly resolved: these infos are
				// for recursive lookup clarity.
				// The concept of "updating" lies somewhere outside of our demesne.
				"/":         {"hub.repeatr.io/base", "2017-05-01", "linux-amd64"},
				"/task/src": {"team.net/theproj", "2.1.1", "src"},
			},
			Formula: map[string]interface{}{
				"inputs": map[rdef.AbsPath]string{
					"/":         "tar:aLMH4qK1EdlPDavdhErOs0BPxqO0i6lUaeRE4DuUmnNMxhHtF56gkoeSulvwWNqT",
					"/task/src": "git:e730adbee91e5584b12dd4cb438673785034ecbe",
				},
				"action": nil, // ... some preprocessor step, whatever ...
				"outputs": map[rdef.AbsPath]interface{}{
					"/task/output/docs": "tar",
					"/task/output/src":  "tar",
				},
			},
			RunRecords: map[rdef.RunRecordHash]*rdef.RunRecord{
				"349h34tq34r9p8u": &rdef.RunRecord{
					UID:       "234852-23792",
					Time:      23495,
					FormulaID: "oeiru43t3ijjrieqo", // somewhat redundantly, the SetupHash of the above formula.
					Results: map[rdef.AbsPath]rdef.WareID{
						"/task/output/docs": "tar:387ty874yt",
						"/task/output/src":  "tar:egruihieur",
					},
				},
				// REVIEW.  A valid alternative way to do this would be putting the
				// blessed set of results here (e.g. only what's picked up by wires),
				// and listing a set of RunRecordHashes that back it up, but storing
				// them very far away (not specific to this step at all).
				// Doing so would mean runrecords aren't stored in "why" order on disk,
				// but would also more closely match how any memoizer tends to see things.
			},
		},
		"build-linux": Step{
			Imports: map[rdef.AbsPath]ReleaseItemID{
				"/":            {"hub.repeatr.io/base", "2017-05-01", "linux-amd64"},
				"/app/compilr": {"hub.repeatr.io/compilr", "1.8", "linux-amd64"},
				"/task/src":    {"wire", "prepare-step", "/task/output/src"},
			},
			Formula: map[string]interface{}{
				"inputs": map[rdef.AbsPath]string{
					"/":                "tar:aLMH4qK1EdlPDavdhErOs0BPxqO0i6lUaeRE4DuUmnNMxhHtF56gkoeSulvwWNqT",
					"/app/compilr":     "tar:jZ8NkMmCPUb5rTHtjBLZEe0usTSDjgGfD71hN07wuuPfkoqG6pLB0FR4GKmQRAva",
					"/task/output/src": "tar:egruihieur",
				},
				"action": nil, // ... some compiler is invoked ...
				"outputs": map[rdef.AbsPath]interface{}{
					"/task/output": "tar",
					"/task/logs":   "tar", // this is a byproduct (implicit: no products point at it).
				},
			},
			RunRecords: map[rdef.RunRecordHash]*rdef.RunRecord{
				"zjklalkjn": &rdef.RunRecord{
					UID:       "21552-2456792",
					Time:      23499,
					FormulaID: "h23hsfiuh48svi",
					Results: map[rdef.AbsPath]rdef.WareID{
						"/task/output": "tar:ooijpwoeijgwer",
						"/task/logs":   "tar:34t983hheiufrt",
					},
				},
				"krljthklj": &rdef.RunRecord{
					UID:       "23456-2456792",
					Time:      23456,
					FormulaID: "h23hsfiuh48svi",
					Results: map[rdef.AbsPath]rdef.WareID{
						"/task/output": "tar:ooijpwoeijgwer",
						"/task/logs":   "tar:poi23459268034",
					},
				},
			},
		},
		// and then you might imagine a "build-mac" step here as well, etc...
	},
	Products: map[ItemLabel]struct {
		StepName   StepName
		OutputSlot rdef.AbsPath
	}{
		"src":         {"prepare-step", "/task/output/src"},
		"docs":        {"prepare-step", "/task/output/docs"},
		"linux-amd64": {"build-linux", "/task/output"},
	},
	// Implicitly now -- For this release record:
	//   - Items["src"] = "tar:egruihieur" -- this much match; correct doc format verifies this: the item label matches the products map, points to the step name, has a runrecord, which has this wareID.
	//   - ... and the other items similarly.
}

// checkpoints: they're named in the same space as commissions, can be wired the same way, they just happen to not need execution of a formula and only have a single output ware (which is coincidentally also the input ware).

// split the computations upstream parts, per matching the verbs file.

// define all these types with full structy goodness; write test structs with primitives and test json roundtrip against them to keep our schema adhered to the ground.
