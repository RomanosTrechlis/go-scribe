package streamer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

// CheckPath checks the validity of a given path
func CheckPath(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("couldn't find path: %v", err)
	}
	info, err := f.Stat()
	if err != nil {
		return fmt.Errorf("couldn't get stat: %v", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("path is not a dir: %v", err)
	}
	return nil
}

func fileExceedsMaxSize(info os.FileInfo, maxSize int64, rootPath, path, filename string) (bool, error) {
	if info == nil {
		return false, nil
	}
	// never change filename due to size constraints
	if maxSize < 0 {
		return false, nil
	}

	if info.Size() < maxSize {
		return false, nil
	}

	oldPath := fmt.Sprintf("%s/%s/%s.log", rootPath, path, filename)
	newPath := fmt.Sprintf("%s/%s/%s_%v.log", rootPath, path, filename, PrintTime())

	err := os.Rename(oldPath, newPath)
	if err != nil {
		return false, fmt.Errorf("failed to rename file exceeding %dbytes: %v", maxSize, err)
	}
	return true, nil
}

func writeLine(rootPath, path, filename, line string, maxSize int64) error {
	logPath := fmt.Sprintf("%s/%s/%s.log", rootPath, path, filename)
	info, err := os.Stat(logPath)
	if os.IsNotExist(err) {
		// path doesn't exist and we need to create it.
		err = os.MkdirAll(filepath.Join(rootPath, path), os.ModePerm)
		if err != nil {
			return fmt.Errorf("couldn't create path '%s': %v",
				filepath.Join(rootPath, path), err)
		}
	}

	createNewFile, err := fileExceedsMaxSize(info, maxSize, rootPath, path, filename)
	if err != nil {
		return fmt.Errorf("failed to rename file: %v", err)
	}
	if createNewFile {
		// re create file if the old has exceeded max size
		err = os.MkdirAll(filepath.Join(rootPath, path), os.ModePerm)
		if err != nil {
			return fmt.Errorf("couldn't create path '%s': %v",
				filepath.Join(rootPath, path), err)
		}
	}

	f, err := os.OpenFile(logPath,
		syscall.O_CREAT|syscall.O_APPEND|syscall.O_WRONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("couldn't create to path '%s': %v", logPath, err)
	}
	defer f.Close()

	if !strings.HasSuffix(line, "\n") {
		line += "\n"
	}
	_, err = f.WriteString(line)
	if err != nil {
		return fmt.Errorf("couldn't write line: %v", err)
	}
	return nil
}
