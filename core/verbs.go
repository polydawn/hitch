package hcore

import (
	"polydawn.net/hitch/api"
	"polydawn.net/hitch/api/rdef"
)

type StagedRelease struct {
	Steps map[api.StepName]StagedStep
}

type StagedStep struct {
	formula  *rdef.Formula
	upstream map[*rdef.AbsPath]api.ReleaseItemID
	records  []*rdef.RunRecord,
}

/*
	Returns keys in upstream which do not map to a formula input path
*/
func (ss *StagedStep) validateUpstreamInjective() ([]rdef.AbsPath, []rdef.AbsPath) {
	formulaOnly := []rdef.AbsPath{}
	for k, _ := range ss.formula.Inputs {
		_, ok := ss.upstream[key]
		if !ok {
			formulaOnly = append(formulaOnly, key)
		}
	}
	return upstreamOnly, formulaOnly
}

/*
	Returns input paths in formula which do not map to an upstream entry
*/
func (ss *StagedStep) validateUpstreamSurjective() []rdef.AbsPath {
	upstreamOnly := []rdef.AbsPath{}
	for key, _ := range ss.upstream {
		_, ok := ss.formula.Inputs[key]
		if !ok {
			upstreamOnly = append(upstreamOnly, key)
		}
	}
	return upstreamOnly
}

func (ss *StagedStep) validateRecords(outputs []rdef.AbsPath) []error {
	// TODO: Validate that all run records have the same output
	// only need to error on wired outputs
	return []error{}
}


/*
	Call to add a step, naming it and providing its formula, and
	optionally the ReleaseItemIDs naming where the inputs were selected from.
*/
func (x *StagedRelease) AppendStep(
	name api.StepName,
	formula *rdef.Formula,
	upstream map[rdef.AbsPath]api.ReleaseItemID, // must be onto (but not necessarily bijection, though lack of may emit warns) the formula inputs.
) error {
	step = StagedStep{
		formuala: formula,
		upstream: upstream,
	}
	results := step.validateUpstreamSurjective()
	if len(results) > 0 {
		return fmt.Errorf("Invalid upstreams: %v", results)
	}
	x.Steps[name] = step
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
) error {
	step, ok = x[name]
	if !ok {
		return fmt.Errorf("Step %s does not exist", name)
	}
	if runRecord == nil {
		return fmt.Errorf("Can't add nil run records")
	}
	step.records = append(step.records, runRecord)
}

/*
	Call Verify or MustVerify to check the connectedness of all steps in the entry so far.

	Calling this after every append is possible if you know you're streaming in
	records in the same order they were built, but it's equally valid to append
	an unordered set of records and then call verify once at the end.
*/
func (x *StagedRelease) MustVerify() {

}
