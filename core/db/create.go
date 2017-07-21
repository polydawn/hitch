package db

import (
	"os"
	"path/filepath"
)

/*
	Create a new hitch db.  Makes a dir; the name itself is the sigil.
*/
func Create(inPath string) (*Controller, error) {
	basePath := filepath.Join(inPath, workspaceName)
	err := os.MkdirAll(basePath, 0755)
	if err != nil {
		return nil, err
	}
	return load(basePath)
}
