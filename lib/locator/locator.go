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

type ErrNotFound struct {
	meep.TraitAutodescribing
}

type ErrCantRead struct {
	meep.TraitAutodescribing
	meep.TraitCausable
}

var fpSepStr string = string([]rune{filepath.Separator})
