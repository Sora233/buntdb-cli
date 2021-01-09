package db

import (
	"errors"
	"github.com/tidwall/buntdb"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var db *buntdb.DB

var dbpath string

func InitBuntDB(filename string) error {
	buntDB, err := buntdb.Open(filename)
	if err != nil {
		return err
	}
	if db != nil {
		db.Close()
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

func GetTempDbPath(prefix string) string {
	fileInfo, err := ioutil.ReadDir(os.TempDir())
	if err == nil {
		for _, fi := range fileInfo {
			if fi.IsDir() {
				continue
			}
			if strings.HasPrefix(filepath.Base(fi.Name()), prefix) {
				return path.Join(os.TempDir(), fi.Name())
			}
		}
	}
	f, err := ioutil.TempFile("", prefix)
	if err != nil {
		return ""
	}
	f.Close()
	return f.Name()
}
