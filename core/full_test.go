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
			Convey("catalog-create finds db in cwd", func() {
				So(Init(ui), ShouldErrorWith, nil)
				So(CatalogCreate(ui, "cn"), ShouldErrorWith, nil)
			})
			Convey("catalog-create finds db when deeper", func() {
				So(Init(ui), ShouldErrorWith, nil)
				So(os.Mkdir("deeper", 0755), ShouldBeNil)
				So(os.Chdir("deeper"), ShouldBeNil)
				So(CatalogCreate(ui, "cn"), ShouldErrorWith, nil)
			})
		})
	})
	Convey("Release staging operations", t, func() {
		WithChdirTmpdir(func() {
			So(Init(ui), ShouldErrorWith, nil)

			Convey("starting a release to a nonexistent catalog should be rejected", func() {
				So(ReleaseStart(ui, "cn", "rn"), ShouldErrorWith, ErrDataNotFound)
			})
			Convey("setting up a simple release with several items, happy path", func() {
				So(CatalogCreate(ui, "cn"), ShouldErrorWith, nil)
				So(ReleaseStart(ui, "cn", "rn"), ShouldErrorWith, nil)
				So(ReleaseAddItem(ui, "label-foo", "tar:asdfasdf"), ShouldErrorWith, nil)
				So(ReleaseAddItem(ui, "label-bar", "tar:asdfasdf"), ShouldErrorWith, nil)
				Convey("adding a clearly invalid wareID should be rejected", func() {
					So(ReleaseAddItem(ui, "label-bar", "not a ware id!"), ShouldErrorWith, ErrBadArgs)
				})
				// TODO : overwrite detection not implemented yet!
				//	Convey("overwriting an item should be rejected", func() {
				//		So(ReleaseAddItem(ui, "label-bar", "tar:asdfasdf"), ShouldErrorWith, ErrOverwrite)
				//	})
				Convey("...should be able to commit!", func() {
					So(ReleaseCommit(ui), ShouldErrorWith, nil)
					Convey("starting a new release with a new name should fly", func() {
						So(ReleaseStart(ui, "cn", "rn-v200"), ShouldErrorWith, nil)
					})
					// TODO : finish giving hitch-release-commit actually impacts!!!
					//	Convey("starting a new release with the same name should be rejected", func() {
					//		So(ReleaseStart(ui, "cn", "rn"), ShouldErrorWith, ErrInProgress)
					//	})
				})
			})
			Convey("starting a release twice should be rejected", func() {
				So(CatalogCreate(ui, "cn"), ShouldErrorWith, nil)
				So(ReleaseStart(ui, "cn", "rn"), ShouldErrorWith, nil)
				So(ReleaseStart(ui, "cn", "zw"), ShouldErrorWith, ErrInProgress)
				// TODO : `hitch release reset` not implemented yet!
				//	Convey("reseting it should allow starting over", func() {
				//		So(ReleaseReset(ui), ShouldErrorWith, nil)
				//		So(ReleaseStart(ui, "cn", "zw"), ShouldErrorWith, nil)
				//	})
			})
		})
	})
}
