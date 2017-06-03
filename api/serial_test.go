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
		Convey("short catalog, one release", func() {
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
	})
}

func jsonPretty(s []byte) string {
	var out bytes.Buffer
	json.Indent(&out, s, "", "\t")
	return out.String()
}
