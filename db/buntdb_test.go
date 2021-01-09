package db

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var testDb = ".test"

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
