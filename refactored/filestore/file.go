package filestore

import (
	"fmt"
	types2 "github.com/brad-hawkins/clean-code/existing/types"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type TempFile struct {
	file os.File
}

// NewTempFile will create a new tempfile and return it as a ReadWriteCloser
func NewTempFile(directory string, filename string, fileType string) (io.ReadWriteCloser, error) {
	err := ensureDirectoryExists(directory)

	fileName := fmt.Sprintf("%s.*.%s", filename, fileType)

	tempFile, err := ioutil.TempFile(directory, fileName)
	if err != nil {
		return nil, err
	}

	return tempFile, nil
}

func ensureDirectoryExists(directory string) error {
	aPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		return types2.WrapError(err, "cannot get absolute path for the sync-app", types2.WithFileSystemError())
	}

	parentDir := filepath.Dir(aPath)

	destDir := filepath.Join(parentDir, directory)

	err = os.MkdirAll(destDir, 0600)
	if err != nil {
		return types2.WrapError(err, "unable to create document directory", types2.WithFileSystemError())
	}
	return nil
}
