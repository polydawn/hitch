package hcore

import (
	"polydawn.net/hitch/api"
	"polydawn.net/hitch/api/rdef"
)

type ReleaseEntryBuilder struct {
}

func (x *ReleaseEntryBuilder) AppendStep(
	name api.CommissionName,
	upstream map[*rdef.AbsPath]struct { // must onto (but not necessarily bijection, though lack of may emit warns) the formula inputs.
		api.CatalogName
		api.ReleaseName
		api.ItemLabel
	},
	formula interface{}, // yes, with hashes.  these HAD BETTER match the upstreams if you check it, but, if upstreams mutate, then, well, that's why we vendored it here.
	runRecord *rdef.RunRecord,
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
