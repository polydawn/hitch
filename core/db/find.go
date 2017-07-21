package db

import (
	"go.polydawn.net/hitch/lib/locator"
)

const workspaceName = "hitch.db"

func LoadByPath(startPath string) (*Controller, error) {
	baseDir, err := locator.FindByDirname(startPath, workspaceName)
	if err != nil {
		return nil, err
	}
	return load(baseDir)
}

func LoadByCwd() (*Controller, error) {
	baseDir, err := locator.FindByDirnameFromCwd(workspaceName)
	if err != nil {
		return nil, err
	}
	return load(baseDir)
}

func load(basePath string) (*Controller, error) {
	ctrl := &Controller{BasePath: basePath}
	return ctrl, nil
}
