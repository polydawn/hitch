package db

import (
	"os"
)

/*
	Create a new hitch db.  Makes a dir, and creates the sigil file.
*/
func Create(basePath string) (*Controller, error) {
	err := os.MkdirAll(basePath, 0755)
	if err != nil {
		return nil, err
	}
	f, err := os.OpenFile(workspaceFilename, os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return load(basePath, f)
}
