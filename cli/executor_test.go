package cli

import (
	"github.com/Sora233/buntdb-cli/db"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/buntdb"
	"os"
	"testing"
)

func TestBuntdbExecutor(t *testing.T) {
	BuntdbExecutor("use -c testcli")
	assert.Equal(t, db.GetDbPath(), "testcli")
	bd, err := db.GetClient()
	assert.Nil(t, err)
	BuntdbExecutor("get a")
	BuntdbExecutor("get -h")
	BuntdbExecutor("get")
	BuntdbExecutor("set a b")
	bd.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get("a")
		assert.Nil(t, err)
		assert.Equal(t, "b", val)
		return nil
	})
	BuntdbExecutor("set a c")
	bd.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get("a")
		assert.Nil(t, err)
		assert.Equal(t, "c", val)
		return nil
	})
	BuntdbExecutor("get a")
	BuntdbExecutor("show db")
	BuntdbExecutor("show index")
	BuntdbExecutor("keys *")
	BuntdbExecutor("del a")
	bd.View(func(tx *buntdb.Tx) error {
		_, err := tx.Get("a")
		assert.Equal(t, buntdb.ErrNotFound, err)
		return nil
	})
	BuntdbExecutor("del b")
	BuntdbExecutor("use testcli-2")
	BuntdbExecutor("use -c testcli-2")
	assert.Equal(t, db.GetDbPath(), "testcli-2")
	BuntdbExecutor("use -c testcli")
	assert.Equal(t, db.GetDbPath(), "testcli")

	os.Remove("testcli")
	os.Remove("testcli-2")
}
