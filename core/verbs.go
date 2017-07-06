package core

import (
	"polydawn.net/hitch/api"
	"polydawn.net/hitch/api/rdef"
)

type ReleaseEntryBuilder struct {
}

func (x *ReleaseEntryBuilder) AppendStep(
	name api.StepName,
	upstream map[*rdef.AbsPath]api.ReleaseItemID, // must onto (but not necessarily bijection, though lack of may emit warns) the formula inputs.
	formula *rdef.Formula, // yes, with hashes.  these HAD BETTER match the upstreams if you check it, but, if upstreams mutate, then, well, that's why we vendored it here.
	runRecord *rdef.RunRecord, // REVIEW maybe append these separately; verify checks if any step has zero runrecords at end.
) {

}

/*
	Call Verify or MustVerify to check the connectedness of all steps in the entry so far.

	Calling this after every append is possible if you know you're streaming in
	records in the same order they were built, but it's equally valid to append
	an unordered set of records and then call verify once at the end.
*/
func (x *ReleaseEntryBuilder) MustVerify() {

}
