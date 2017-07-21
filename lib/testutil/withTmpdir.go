package testutil

import (
	"io/ioutil"
	"os"
)

// Run a function, giving it a string naming a tempdir it may use.
// The tempdir will be removed when the function returns.
func WithTmpdir(fn func(tmpDir string)) {
	tmpBase := "/tmp/hitch-test/"
	err := os.MkdirAll(tmpBase, os.FileMode(0755)|os.ModeSticky)
	if err != nil {
		panic(err)
	}

	tmpDir, err := ioutil.TempDir(tmpBase, "")
	if err != nil {
		panic(err)
	}

	defer os.RemoveAll(tmpDir)
	fn(tmpDir)
}

// Run a function, after *changing process cwd* to a tempdir it may use.
// The cwd will be reset and the tempdir removed when the function returns.
func WithChdirTmpdir(fn func()) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	WithTmpdir(func(tmpDir string) {
		if err := os.Chdir(tmpDir); err != nil {
			panic(err)
		}
		defer os.Chdir(cwd)
		fn()
	})
}
