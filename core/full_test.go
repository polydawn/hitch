package core

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	. "github.com/polydawn/go-errcat"
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
				So(Init(ui), ErrorShouldHaveCategory, nil)
			})
			Convey("init repeated should fail", func() {
				So(Init(ui), ErrorShouldHaveCategory, nil)
				So(Init(ui), ErrorShouldHaveCategory, ErrInProgress)
			})
			Convey("init inside another db should fail", func() {
				So(Init(ui), ErrorShouldHaveCategory, nil)
				So(os.Mkdir("deeper", 0755), ShouldBeNil)
				So(os.Chdir("deeper"), ShouldBeNil)
				So(Init(ui), ErrorShouldHaveCategory, ErrInProgress)
			})
			Convey("catalog-create finds db in cwd", func() {
				So(Init(ui), ErrorShouldHaveCategory, nil)
				So(CatalogCreate(ui, "cn"), ErrorShouldHaveCategory, nil)
			})
			Convey("catalog-create finds db when deeper", func() {
				So(Init(ui), ErrorShouldHaveCategory, nil)
				So(os.Mkdir("deeper", 0755), ShouldBeNil)
				So(os.Chdir("deeper"), ShouldBeNil)
				So(CatalogCreate(ui, "cn"), ErrorShouldHaveCategory, nil)
			})
		})
	})
	Convey("Release staging operations", t, func() {
		WithChdirTmpdir(func() {
			Must(Init(ui))

			Convey("starting a release to a nonexistent catalog should be rejected", func() {
				So(ReleaseStart(ui, "cn", "rn"), ErrorShouldHaveCategory, ErrDataNotFound)
			})
			Convey("setting up a simple release with several items, happy path", func() {
				So(CatalogCreate(ui, "cn"), ErrorShouldHaveCategory, nil)
				So(ReleaseStart(ui, "cn", "rn"), ErrorShouldHaveCategory, nil)
				So(ReleaseAddItem(ui, "label-foo", "tar:asdfasdf"), ErrorShouldHaveCategory, nil)
				So(ReleaseAddItem(ui, "label-bar", "tar:asdfasdf"), ErrorShouldHaveCategory, nil)
				Convey("adding a clearly invalid wareID should be rejected", func() {
					So(ReleaseAddItem(ui, "label-bar", "not a ware id!"), ErrorShouldHaveCategory, ErrBadArgs)
				})
				// TODO : overwrite detection not implemented yet!
				//	Convey("overwriting an item should be rejected", func() {
				//		So(ReleaseAddItem(ui, "label-bar", "tar:asdfasdf"), ErrorShouldHaveCategory, ErrOverwrite)
				//	})
				Convey("...should be able to commit!", func() {
					So(ReleaseCommit(ui), ErrorShouldHaveCategory, nil)
					Convey("starting a new release with a new name should fly", func() {
						So(ReleaseStart(ui, "cn", "rn-v200"), ErrorShouldHaveCategory, nil)
					})
					Convey("starting a new release with the same name should be rejected", func() {
						So(ReleaseStart(ui, "cn", "rn"), ErrorShouldHaveCategory, ErrNameCollision)
					})
				})
			})
			Convey("starting a release twice should be rejected", func() {
				So(CatalogCreate(ui, "cn"), ErrorShouldHaveCategory, nil)
				So(ReleaseStart(ui, "cn", "rn"), ErrorShouldHaveCategory, nil)
				So(ReleaseStart(ui, "cn", "zw"), ErrorShouldHaveCategory, ErrInProgress)
				// TODO : `hitch release reset` not implemented yet!
				//	Convey("reseting it should allow starting over", func() {
				//		So(ReleaseReset(ui), ErrorShouldHaveCategory, nil)
				//		So(ReleaseStart(ui, "cn", "zw"), ErrorShouldHaveCategory, nil)
				//	})
			})
		})
	})
	Convey("Show command operations", t, func() {
		WithChdirTmpdir(func() {
			Must(Init(ui))

			Convey("given a sizable catalog", func() {
				Must(CatalogCreate(ui, "cn"))
				Must(ReleaseStart(ui, "cn", "v0.1"))
				Must(ReleaseAddItem(ui, "label-foo", "tar:asdfasdf"))
				Must(ReleaseAddItem(ui, "label-bar", "tar:asdfqwer"))
				Must(ReleaseCommit(ui))
				Must(ReleaseStart(ui, "cn", "v0.2"))
				Must(ReleaseAddItem(ui, "label-foo", "tar:qwerasdf"))
				Must(ReleaseAddItem(ui, "label-bar", "tar:qwerqwer"))
				Must(ReleaseAddItem(ui, "label-qux", "tar:qwerzxcv"))
				Must(ReleaseCommit(ui))

				Convey("`hitch show <catalog>` should say a *lot*", func() {
					output, err := grabOutput(func(ui UI) error {
						return Show(ui, "cn")
					})
					So(err, ErrorShouldHaveCategory, nil)
					So(output, ShouldContainSubstring, `"cn"`)              // catalogs contain their own name
					So(output, ShouldContainSubstring, `"v0.1"`)            // both release names should appear
					So(output, ShouldContainSubstring, `"v0.2"`)            // both release names should appear
					So(output, ShouldContainSubstring, `"tar:asdfasdf"`)    // wareIDs in the first release should appear -- showing a catalog is *loud*!
					So(output, ShouldContainSubstring, `"tar:qwerqwer"`)    // wareIDs in the second release should appear -- showing a catalog is *loud*!
					So(strings.Count(output, `"metadata"`), ShouldEqual, 2) // keyword "metadata" should appear once per release entry
				})
				Convey("`hitch show <catalog>` requesting a non-existent catalog name should error", func() {
					So(Show(ui, "notgonnafindit"), ErrorShouldHaveCategory, ErrDataNotFound)
				})
				Convey("`hitch show <catalog>:<release>` should only show that release", func() {
					output, err := grabOutput(func(ui UI) error {
						return Show(ui, "cn:v0.1")
					})
					So(err, ErrorShouldHaveCategory, nil)
					So(output, ShouldContainSubstring, `"v0.1"`)            // releases contain their own name
					So(output, ShouldContainSubstring, `"tar:asdfasdf"`)    // wareIDs in the first release should appear
					So(strings.Count(output, `"metadata"`), ShouldEqual, 1) // keyword "metadata" should appear once, because it's just one release entry
					So(output, ShouldNotContainSubstring, `"cn"`)           // releases don't repeat the name of the catalog that contains them
					So(output, ShouldNotContainSubstring, `"v0.2"`)         // the other releases names certainly shouldn't appear
					So(output, ShouldNotContainSubstring, `"tar:qwerqwer"`) // wareIDs in the other releases certainly shouldn't should appear
				})
				Convey("`hitch show <catalog>:<release>` requesting a non-existent catalog name should error", func() {
					So(Show(ui, "notgonnafindit:wompwomp"), ErrorShouldHaveCategory, ErrDataNotFound)
				})
				Convey("`hitch show <catalog>:<release>` requesting a non-existent release name should error", func() {
					So(Show(ui, "cn:wompwomp"), ErrorShouldHaveCategory, ErrDataNotFound)
				})
				Convey("`hitch show <catalog>:<release>:<item>` should only show that WareID -- unquoted!", func() {
					output, err := grabOutput(func(ui UI) error {
						return Show(ui, "cn:v0.1:label-foo")
					})
					So(err, ErrorShouldHaveCategory, nil)
					So(output, ShouldEqual, "tar:asdfasdf\n")
				})
				Convey("`hitch show <catalog>:<release>:<item>` requesting a non-existent item name should error", func() {
					So(Show(ui, "cn:v0.1:notathing"), ErrorShouldHaveCategory, ErrDataNotFound)
				})
			})

		})
	})
}

func grabOutput(fn func(UI) error) (string, error) {
	var buf bytes.Buffer
	err := fn(UI{nil, &buf, &buf})
	return buf.String(), err
}
