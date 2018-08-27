package fs

import (
	"fmt"
	"os"
)

func CreateFolderIfNotExist(path string) (err error) {
	_, err = os.Stat(path)
	if err == nil {
		return nil
	}

	if !os.IsNotExist(err) {
		return fmt.Errorf("error accessing directory %s: %v", path, err)
	}

	err = os.Mkdir(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating directory %s: %v", path, err)
	}
	return nil
}
