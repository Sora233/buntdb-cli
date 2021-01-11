package db

import (
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/buntdb"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var testDb = ".testdb"
var testDb2 = ".testdb-2"

func TestInitBuntDB(t *testing.T) {
	assert.Nil(t, InitBuntDB(testDb))
	assert.Nil(t, Close())
	os.Remove(testDb)
}

func TestGetClient(t *testing.T) {
	db, err := GetClient()
	assert.NotNil(t, err)
	assert.Nil(t, db)

	InitBuntDB(testDb)
	defer os.Remove(testDb)
	defer Close()

	db, err = GetClient()
	assert.Nil(t, err)
	assert.NotNil(t, db)

	assert.Equal(t, GetDbPath(), testDb)
}

func TestSwitchBuntDB(t *testing.T) {
	assert.Nil(t, InitBuntDB(testDb))
	db, err := GetClient()
	assert.Nil(t, err)
	err = db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set("a", "b", nil)
		return err
	})
	assert.Nil(t, err)
	assert.Nil(t, InitBuntDB(testDb2))
	db, err = GetClient()
	assert.Nil(t, err)
	err = db.View(func(tx *buntdb.Tx) error {
		_, err := tx.Get("a")
		assert.Equal(t, buntdb.ErrNotFound, err)
		return nil
	})
	Close()
	os.Remove(testDb)
	os.Remove(testDb2)

}

func TestGetTempDbPath(t *testing.T) {
	path1 := GetTempDbPath("tmppath")
	assert.True(t, strings.HasPrefix(filepath.Base(path1), "tmppath"))
	path2 := GetTempDbPath("tmppath")
	assert.True(t, strings.HasPrefix(filepath.Base(path2), "tmppath"))
	assert.Equal(t, path1, path2)
	os.Remove(path1)
	os.Remove(path2)
}

func TestRWDescribe(t *testing.T) {
	assert.Equal(t, "rw", RWDescribe(true))
	assert.Equal(t, "r", RWDescribe(false))
}

func TestBegin(t *testing.T) {
	assert.Nil(t, InitBuntDB(testDb))
	tx, err := Begin(true)
	assert.Nil(t, err)
	_, err = tx.Get("a")
	assert.Equal(t, buntdb.ErrNotFound, err)
	_, _, err = tx.Set("a", "a", nil)
	assert.Nil(t, err)
	val, err := tx.Get("a")
	assert.Nil(t, err)
	assert.Equal(t, "a", val)
	assert.Nil(t, Rollback())

	db, _ := GetClient()
	db.View(func(tx *buntdb.Tx) error {
		_, err := tx.Get("a")
		assert.Equal(t, buntdb.ErrNotFound, err)
		return nil
	})

	tx, err = Begin(false)
	assert.Nil(t, err)
	_, _, err = tx.Set("a", "a", nil)
	assert.Equal(t, buntdb.ErrTxNotWritable, err)
	assert.NotNil(t, Commit())
	assert.Nil(t, Rollback())

	assert.Equal(t, ErrNoTransaction, Commit())
	assert.Equal(t, ErrNoTransaction, Rollback())
	tx, err = Begin(true)
	assert.Nil(t, err)
	tx, err = Begin(false)
	assert.Equal(t, ErrNestedTransaction, err)

	assert.Equal(t, ErrTransactionExist, InitBuntDB(testDb))
	assert.Nil(t, Rollback())
	Close()

	os.Remove(testDb)
	os.Remove(testDb2)
}

func TestGetCurrentTransaction(t *testing.T) {
	assert.Nil(t, InitBuntDB(testDb))
	Begin(true)
	tx, rw := GetCurrentTransaction()
	assert.True(t, rw)
	assert.NotNil(t, tx)
	Commit()

	Begin(false)
	tx, rw = GetCurrentTransaction()
	assert.False(t, rw)
	assert.NotNil(t, tx)
	Shrink()
	Rollback()
	tx, rw = GetCurrentTransaction()
	assert.Nil(t, tx)
	Shrink()
	os.Remove(testDb)
	os.Remove(testDb2)
}

func TestSave(t *testing.T) {
	assert.Nil(t, InitBuntDB(testDb))
	db, err := GetClient()
	assert.Nil(t, err)
	err = db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set("a", "testsave", nil)
		return err
	})
	assert.Nil(t, err)
	f, err := ioutil.TempFile("", "test_save")
	assert.Nil(t, err)
	Save(f)
	Close()
	f.Close()

	db, err = buntdb.Open(f.Name())
	assert.Nil(t, err)
	db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get("a")
		assert.Nil(t, err)
		assert.Equal(t, "testsave", val)
		return nil
	})
	db.Close()
	os.Remove(f.Name())
	os.Remove(testDb)
	os.Remove(testDb2)
}
