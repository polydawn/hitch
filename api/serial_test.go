package api

import (
	"encoding/json"
	"testing"

	"github.com/polydawn/refmt"
	"github.com/polydawn/refmt/obj/atlas"
	. "github.com/smartystreets/goconvey/convey"
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
	})
}
