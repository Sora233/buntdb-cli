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
var tx *buntdb.Tx
var writable bool

var dbpath string

func InitBuntDB(filename string) error {
	if tx != nil {
		return ErrTransactionExist
	}
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

func Begin(_writable bool) (*buntdb.Tx, error) {
	if tx != nil {
		return nil, ErrNestedTransaction
	}
	bd, err := GetClient()
	if err != nil {
		return nil, err
	}
	_tx, err := bd.Begin(_writable)
	if err != nil {
		return nil, err
	}
	tx = _tx
	writable = _writable
	return tx, err
}

func Commit() error {
	if tx == nil {
		return ErrNoTransaction
	}

	err := tx.Commit()
	if err != nil {
		if err == buntdb.ErrTxNotWritable {
			return errors.New("readonly transaction can only rollback")
		}
		return err
	} else {
		tx = nil
		return nil
	}
}

func Rollback() error {
	if tx == nil {
		return ErrNoTransaction
	}
	err := tx.Rollback()
	if err != nil {
		return err
	} else {
		tx = nil
		return nil
	}
}

func GetCurrentTransaction() (*buntdb.Tx, bool) {
	return tx, writable
}

func RWDescribe(writable bool) string {
	if writable {
		return "rw"
	} else {
		return "r"
	}
}
