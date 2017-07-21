package locator

import (
	"os"
	"path/filepath"

	"go.polydawn.net/meep"
)

/*
	Look for a workspace config file:
	starting from `startPath`,
	iterating up through parent directories,
	looking for files named `fileName`,
	and returning either the dir found and open file, or error.

	Errors:

		- *ErrNotFound -- if we never find a file named 'fileName'
		- *ErrCantRead -- if any other IO errors opening the file for reading
*/
func FindByFilename(startPath string, fileName string) (dirFound string, f *os.File, erm error) {
	dirFound = filepath.Clean(startPath)
	for {
		f, err := os.Open(filepath.Join(dirFound, fileName))
		if err == nil {
			return dirFound, f, nil
		}
		defer f.Close()
		if os.IsNotExist(err) {
			dirFound = filepath.Dir(dirFound)
			if dirFound == "." || dirFound == fpSepStr {
				return "", nil, meep.Meep(&ErrNotFound{})
			}
			continue
		}
		return dirFound, f, meep.Meep(&ErrCantRead{}, meep.Cause(err))
	}
}

/*
	As per FindByFilename, but the start dir is the process current working directory.

	Errors:

		- the errors from `FindByFilename()`
		- *ErrCantRead -- if the current process CWD cannot be detected.
*/
func FindByFilenameFromCwd(fileName string) (dirFound string, f *os.File, erm error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", nil, meep.Meep(&ErrCantRead{}, meep.Cause(err))
	}
	return FindByFilename(cwd, fileName)
}

/*
	Look for a workspace indicated by a directory with special name:
	starting from `startPath`,
	iterating up through parent directories,
	looking for dirs in each parent dir named `dirName`,
	and returning either the path found with that name, or error.

	Symlinks to directories are also acceptable.

	Errors:

		- *ErrNotFound -- if we never find a dir named 'dirName'
		- *ErrCantRead -- if any other IO errors iterating up the parent dirs for reading
*/
func FindByDirname(startPath string, dirName string) (string, error) {
	dirFound := filepath.Clean(startPath)
	if filepath.Base(dirFound) == dirName {
		return dirFound, nil
	}
	for {
		// `ls`.  Any errors: return.
		f, err := os.Open(dirFound)
		if err != nil {
			if os.IsNotExist(err) {
				return "", meep.Meep(&ErrNotFound{})
			}
			return "", meep.Meep(&ErrCantRead{}, meep.Cause(err))
		}
		defer f.Close()
		// Scan through all entries in the dir, looking for our fav.
		names, err := f.Readdirnames(-1)
		if err != nil {
			return "", meep.Meep(&ErrCantRead{}, meep.Cause(err))
		}
		for _, name := range names {
			if name != dirName {
				continue
			}
			dirFound := filepath.Join(dirFound, name)
			fi, err := os.Stat(dirFound)
			if err != nil {
				return "", meep.Meep(&ErrCantRead{}, meep.Cause(err))
			}
			if !fi.IsDir() {
				break
			}
			return dirFound, nil
		}
		// If basename'ing got us "/" this time, and we still didn't find it, terminate.
		if dirFound == "/" || dirFound == "." {
			return "", meep.Meep(&ErrNotFound{})
		}
		// Step up one dir.
		dirFound = filepath.Dir(dirFound)
	}
}

/*
	As per FindByDirname, but the start dir is the process current working directory.

	Errors:

		- the errors from `FindByDirname()`
		- *ErrCantRead -- if the current process CWD cannot be detected.
*/
func FindByDirnameFromCwd(dirName string) (dirFound string, err error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", meep.Meep(&ErrCantRead{}, meep.Cause(err))
	}
	return FindByDirname(cwd, dirName)
}

type ErrNotFound struct {
	meep.TraitAutodescribing
}

type ErrCantRead struct {
	meep.TraitAutodescribing
	meep.TraitCausable
}

var fpSepStr string = string([]rune{filepath.Separator})
