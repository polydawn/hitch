package core

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

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
				So(Init(ui), ShouldEqual, EXIT_SUCCESS)
			})
			Convey("init repeated should fail", func() {
				So(Init(ui), ShouldEqual, EXIT_SUCCESS)
				So(Init(ui), ShouldEqual, EXIT_INPROGRESS)
			})
			Convey("init inside another db should fail", func() {
				So(Init(ui), ShouldEqual, EXIT_SUCCESS)
				So(os.Mkdir("deeper", 0755), ShouldBeNil)
				So(os.Chdir("deeper"), ShouldBeNil)
				So(Init(ui), ShouldEqual, EXIT_INPROGRESS)
			})
		})
	})
}
