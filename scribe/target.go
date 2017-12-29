package scribe

import (
	"fmt"

	"gopkg.in/mgo.v2"
)

type target struct {
	isDB bool
	databaseName string
	session *mgo.Session

	isFile bool
	rootPath string
	fileSize int64
}

func NewTarget(isDB, isFile bool, dbName, dbStore, root string, fileSize int64) (*target, error){
	target := &target {
		databaseName:dbStore,
		isFile:isFile,
		isDB:isDB,
		rootPath:root,
		fileSize:fileSize,
	}

	if isDB {
		s, err := mgo.Dial(dbName)
		if err != nil {
			return nil, fmt.Errorf("failed to dial database: %v", err)
		}
		target.session = s
	}

	return target, nil
}