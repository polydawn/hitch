package db

import (
	"os"

	"go.polydawn.net/hitch/lib/locator"
)

const workspaceFilename = ".hitch"

func LoadByPath(startPath string) (*Controller, error) {
	rootDir, f, err := locator.FindByFilename(startPath, workspaceFilename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return load(rootDir, f)
}

func LoadByCwd() (*Controller, error) {
	rootDir, f, err := locator.FindByFilenameFromCwd(workspaceFilename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return load(rootDir, f)
}

func load(basePath string, f *os.File) (*Controller, error) {
	ctrl := &Controller{basePath: basePath}
	// When we actually have a main info file: parse it here.
	return ctrl, nil
}
