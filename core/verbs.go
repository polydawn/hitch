package hcore

import (
	"polydawn.net/hitch/api"
	"polydawn.net/hitch/api/rdef"
)

type StagedRelease struct {
}

/*
	Call to add a step, naming it and providing its formula, and
	optionally the ReleaseItemIDs naming where the inputs were selected from.
*/
func (x *StagedRelease) AppendStep(
	name api.StepName,
	formula *rdef.Formula,
	upstream map[*rdef.AbsPath]api.ReleaseItemID, // must onto (but not necessarily bijection, though lack of may emit warns) the formula inputs.
) {

}

/*
	Call to add a RunRecord to a step.
	The step name must already have been created by `AppendStep`.

	The RunRecord's recorded formula setuphash must match the setuphash
	of the formula for this step, or it is invalid.
	The RunRecord's outputs must match the output slots named in the formula
	for this step, or it is invalid (unless you're mucking with the document,
	output from `repeatr run` for the formula should always match).

	TODO: define at what point we check multiple RunRecords for coherency
	on wired outputs
*/
func (x *StagedRelease) AppendRunRecord(
	name api.StepName,
	runRecord *rdef.RunRecord,
) {

}

/*
	Call Verify or MustVerify to check the connectedness of all steps in the entry so far.

	Calling this after every append is possible if you know you're streaming in
	records in the same order they were built, but it's equally valid to append
	an unordered set of records and then call verify once at the end.
*/
func (x *StagedRelease) MustVerify() {

}
