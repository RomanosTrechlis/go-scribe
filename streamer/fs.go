package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	pb "github.com/RomanosTrechlis/logStream/api"
)

func checkPath(path string) error {
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

func fileExceedsMaxSize(info os.FileInfo, req *pb.LogRequest, maxSize int, path string) (bool, error) {
	if info == nil {
		return false, nil
	}
	// never change filename due to size constraints
	if maxSize < 0 {
		return false, nil
	}

	if info.Size() < int64(maxSize) {
		return false, nil
	}

	oldPath := fmt.Sprintf("%s/%s/%s.log", path, req.GetPath(), req.GetFilename())
	newPath := fmt.Sprintf("%s/%s/%s_%v.log", path, req.GetPath(), req.GetFilename(), time.Now())

	err := os.Rename(oldPath, newPath)
	if err != nil {
		return false, fmt.Errorf("failed to rename file exceeding %dbytes: %v", maxSize, err)
	}
	return true, nil
}

func writeLine(req *pb.LogRequest, path string, maxSize int) error {
	logPath := fmt.Sprintf("%s/%s/%s.log", path, req.GetPath(), req.GetFilename())
	info, err := os.Stat(logPath)
	if os.IsNotExist(err) {
		// path doesn't exist and we need to create it.
		err = os.MkdirAll(filepath.Join(path, req.GetPath()), os.ModePerm)
		if err != nil {
			return fmt.Errorf("couldn't create path '%s': %v",
				filepath.Join(path, req.GetPath()), err)
		}
	}

	createNewFile, err := fileExceedsMaxSize(info, req, maxSize, path)
	if err != nil {
		return fmt.Errorf("failed to rename file: %v", err)
	}
	if createNewFile {
		// re create file if the old has exceeded max size
		err = os.MkdirAll(filepath.Join(path, req.GetPath()), os.ModePerm)
		if err != nil {
			return fmt.Errorf("couldn't create path '%s': %v",
				filepath.Join(path, req.GetPath()), err)
		}
	}

	f, err := os.OpenFile(logPath,
		syscall.O_CREAT|syscall.O_APPEND|syscall.O_WRONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("couldn't create to path '%s': %v", logPath, err)
	}
	defer f.Close()

	line := req.GetLine()
	if !strings.HasSuffix(line, "\n") {
		line += "\n"
	}
	_, err = f.WriteString(line)
	if err != nil {
		return fmt.Errorf("couldn't write line: %v", err)
	}
	return nil
}
