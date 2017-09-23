package main

import (
	"fmt"
	"os"
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

func fileExceedsMaxSize(info os.FileInfo, req *pb.LogRequest) (bool, error) {
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
