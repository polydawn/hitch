package core

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	. "go.polydawn.net/hitch/lib/errcat"
	. "go.polydawn.net/hitch/lib/testutil"
)

func Test(t *testing.T) {
	ui := UI{
		nil,
		ioutil.Discard,
		ioutil.Discard,
	}
	Convey("Initialization / db-finding scenarios", t, func() {
		WithChdirTmpdir(func() {
			Convey("init happy path", func() {
				So(Init(ui), ShouldErrorWith, nil)
			})
			Convey("init repeated should fail", func() {
				So(Init(ui), ShouldErrorWith, nil)
				So(Init(ui), ShouldErrorWith, ErrInProgress)
			})
			Convey("init inside another db should fail", func() {
				So(Init(ui), ShouldErrorWith, nil)
				So(os.Mkdir("deeper", 0755), ShouldBeNil)
				So(os.Chdir("deeper"), ShouldBeNil)
				So(Init(ui), ShouldErrorWith, ErrInProgress)
			})
			Convey("release-start finds db in cwd", func() {
				So(Init(ui), ShouldErrorWith, nil)
				So(ReleaseStart(ui, "cn", "rn"), ShouldErrorWith, nil)
			})
			Convey("release-start finds db when deeper", func() {
				So(Init(ui), ShouldErrorWith, nil)
				So(os.Mkdir("deeper", 0755), ShouldBeNil)
				So(os.Chdir("deeper"), ShouldBeNil)
				So(ReleaseStart(ui, "cn", "rn"), ShouldErrorWith, nil)
			})
		})
	})
}
