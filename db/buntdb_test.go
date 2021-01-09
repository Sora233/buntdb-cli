package db

import (
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/buntdb"
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

func TestGetDbPath(t *testing.T) {
	InitBuntDB(testDb)
	defer os.Remove(testDb)
	defer Close()
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
	path1 := GetTempDbPath("testcli")
	assert.True(t, strings.HasPrefix(filepath.Base(path1), "testcli"))
	path2 := GetTempDbPath("testcli")
	assert.True(t, strings.HasPrefix(filepath.Base(path2), "testcli"))
	assert.Equal(t, path1, path2)
	Close()
	os.Remove(testDb)
	os.Remove(testDb2)
}
