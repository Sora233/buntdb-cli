package db

import (
	"errors"
	"github.com/tidwall/buntdb"
)

var db *buntdb.DB

var dbpath string

func InitBuntDB(filename string) error {
	buntDB, err := buntdb.Open(filename)
	if err != nil {
		return err
	}
	db = buntDB
	db.SetConfig(buntdb.Config{
		SyncPolicy: buntdb.Always,
	})
	dbpath = filename
	return nil
}

func GetClient() (*buntdb.DB, error) {
	if db == nil {
		return nil, errors.New("not initialized")
	}
	return db, nil
}

func GetDbPath() string {
	return dbpath
}

func Close() error {
	if db != nil {
		if err := db.Close(); err != nil {
			return err
		}
		db = nil
	}
	return nil
}
