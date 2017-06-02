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
