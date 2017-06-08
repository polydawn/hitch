package api

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/polydawn/refmt"
	"github.com/polydawn/refmt/obj/atlas"
	. "github.com/smartystreets/goconvey/convey"

	"polydawn.net/hitch/api/rdef"
)

func TestSerial(t *testing.T) {
	atl := atlas.MustBuild(
		ReleaseItemID_AtlasEntry,
		Catalog_AtlasEntry,
		ReleaseEntry_AtlasEntry,
		Replay_AtlasEntry,
		Step_AtlasEntry,
		rdef.RunRecord_AtlasEntry,
	)
	Convey("ReleaseItemID serialization", t, func() {
		msg, err := refmt.JsonEncodeAtlased(atl,
			ReleaseItemID{"a", "b", "c"})
		So(err, ShouldBeNil)
		So(string(msg), ShouldResemble, `"a:b:c"`)
		var reheat string
		So(json.Unmarshal(msg, &reheat), ShouldBeNil)
		So(reheat, ShouldResemble, "a:b:c")
	})
	Convey("Catalog serialization", t, func() {
		Convey("empty catalog, no releases", func() {
			msg, err := refmt.JsonEncodeAtlased(atl,
				Catalog{
					"cname",
					[]ReleaseEntry{},
				})
			So(err, ShouldBeNil)
			So(string(msg), ShouldResemble, `{"name":"cname","releases":[]}`)
		})
		Convey("short catalog: one release, no replay", func() {
			msg, err := refmt.JsonEncodeAtlased(atl,
				Catalog{
					"cname",
					[]ReleaseEntry{
						{"1.0",
							map[ItemLabel]rdef.WareID{
								"item-a": "war:asdf",
								"item-b": "war:qwer",
							},
							map[string]string{
								"comment": "yes",
							},
							nil,
							nil,
						},
					},
				})
			So(err, ShouldBeNil)
			So(jsonPretty(msg), ShouldResemble,
				`{
	"name": "cname",
	"releases": [
		{
			"name": "1.0",
			"items": {
				"item-a": "war:asdf",
				"item-b": "war:qwer"
			},
			"metadata": {
				"comment": "yes"
			},
			"hazards": null,
			"replay": null
		}
	]
}`)
		})

		Convey("medium catalog: multiple releases, interesting replays", func() {
			msg, err := refmt.JsonEncodeAtlased(atl,
				Catalog{
					"cname",
					[]ReleaseEntry{
						{"1.1",
							map[ItemLabel]rdef.WareID{
								"src":         "tar:egruihieur",
								"docs":        "tar:387ty874yt",
								"linux-amd64": "tar:ooijpwoeij",
							},
							nil,
							nil,
							&Replay{
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
													"/task/output": "tar:ooijpwoeij",
													"/task/logs":   "tar:34t983hhei",
												},
											},
											"krljthklj": &rdef.RunRecord{
												UID:       "23456-2456792",
												Time:      23456,
												FormulaID: "h23hsfiuh48svi",
												Results: map[rdef.AbsPath]rdef.WareID{
													"/task/output": "tar:ooijpwoeij",
													"/task/logs":   "tar:poi2345926",
												},
											},
										},
									},
								},
								// The products map has the same keys as the Release's ItemLabels.
								// When these wires are resolved, the indicated RunRecords in this
								// Replay MUST express the exact same WareID hash
								// as the ReleaseRecord's top level ItemLabel->WareID map.
								Products: map[ItemLabel]ReleaseItemID{
									"src":         {"wire", "prepare-step", "/task/output/src"},
									"docs":        {"wire", "prepare-step", "/task/output/docs"},
									"linux-amd64": {"wire", "build-linux", "/task/output"},
								},
							},
						},
						{"1.0",
							map[ItemLabel]rdef.WareID{
								"src":         "war:asdf",
								"docs":        "war:ayhf",
								"linux-amd64": "war:qwer",
							},
							nil,
							map[string]string{
								"unreproducible": "missing replay! ;)",
							},
							nil,
						},
					},
				})
			So(err, ShouldBeNil)
			So(jsonPretty(msg), ShouldResemble,
				`{
	"name": "cname",
	"releases": [
		{
			"name": "1.1",
			"items": {
				"docs": "tar:387ty874yt",
				"linux-amd64": "tar:ooijpwoeij",
				"src": "tar:egruihieur"
			},
			"metadata": null,
			"hazards": null,
			"replay": {
				"steps": {
					"build-linux": {
						"imports": {
							"/": "hub.repeatr.io/base:2017-05-01:linux-amd64",
							"/app/compilr": "hub.repeatr.io/compilr:1.8:linux-amd64",
							"/task/src": "wire:prepare-step:/task/output/src"
						},
						"formula": {
							"action": null,
							"inputs": {
								"/": "tar:aLMH4qK1EdlPDavdhErOs0BPxqO0i6lUaeRE4DuUmnNMxhHtF56gkoeSulvwWNqT",
								"/app/compilr": "tar:jZ8NkMmCPUb5rTHtjBLZEe0usTSDjgGfD71hN07wuuPfkoqG6pLB0FR4GKmQRAva",
								"/task/output/src": "tar:egruihieur"
							},
							"outputs": {
								"/task/logs": "tar",
								"/task/output": "tar"
							}
						},
						"runRecords": {
							"krljthklj": {
								"uID": "23456-2456792",
								"time": 23456,
								"formulaID": "h23hsfiuh48svi",
								"results": {
									"/task/logs": "tar:poi2345926",
									"/task/output": "tar:ooijpwoeij"
								},
								"hostname": "",
								"metadata": null
							},
							"zjklalkjn": {
								"uID": "21552-2456792",
								"time": 23499,
								"formulaID": "h23hsfiuh48svi",
								"results": {
									"/task/logs": "tar:34t983hhei",
									"/task/output": "tar:ooijpwoeij"
								},
								"hostname": "",
								"metadata": null
							}
						}
					},
					"prepare-step": {
						"imports": {
							"/": "hub.repeatr.io/base:2017-05-01:linux-amd64",
							"/task/src": "team.net/theproj:2.1.1:src"
						},
						"formula": {
							"action": null,
							"inputs": {
								"/": "tar:aLMH4qK1EdlPDavdhErOs0BPxqO0i6lUaeRE4DuUmnNMxhHtF56gkoeSulvwWNqT",
								"/task/src": "git:e730adbee91e5584b12dd4cb438673785034ecbe"
							},
							"outputs": {
								"/task/output/docs": "tar",
								"/task/output/src": "tar"
							}
						},
						"runRecords": {
							"349h34tq34r9p8u": {
								"uID": "234852-23792",
								"time": 23495,
								"formulaID": "oeiru43t3ijjrieqo",
								"results": {
									"/task/output/docs": "tar:387ty874yt",
									"/task/output/src": "tar:egruihieur"
								},
								"hostname": "",
								"metadata": null
							}
						}
					}
				},
				"products": {
					"docs": "wire:prepare-step:/task/output/docs",
					"linux-amd64": "wire:build-linux:/task/output",
					"src": "wire:prepare-step:/task/output/src"
				}
			}
		},
		{
			"name": "1.0",
			"items": {
				"docs": "war:ayhf",
				"linux-amd64": "war:qwer",
				"src": "war:asdf"
			},
			"metadata": null,
			"hazards": {
				"unreproducible": "missing replay! ;)"
			},
			"replay": null
		}
	]
}`)
		})
	})
}

func jsonPretty(s []byte) string {
	var out bytes.Buffer
	json.Indent(&out, s, "", "\t")
	return out.String()
}
