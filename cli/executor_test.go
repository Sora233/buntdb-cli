package cli

import (
	"github.com/Sora233/buntdb-cli/db"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/buntdb"
	"os"
	"testing"
)

func TestBuntdbExecutor(t *testing.T) {
	Debug = true
	BuntdbExecutor("use -c testcli")
	assert.Equal(t, db.GetDbPath(), "testcli")
	bd, err := db.GetClient()
	assert.Nil(t, err)
	BuntdbExecutor("fake")
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
	BuntdbExecutor("set a d 999")
	bd.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get("a")
		assert.Nil(t, err)
		assert.Equal(t, "d", val)
		ttl, err := tx.TTL("a")
		assert.Nil(t, err)
		assert.Greater(t, ttl.Seconds(), 0.0)
		return nil
	})
	BuntdbExecutor("ttl a")
	BuntdbExecutor("ttl b")
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
	BuntdbExecutor(`set "a b" c`)
	bd.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get("a b")
		assert.Nil(t, err)
		assert.Equal(t, "c", val)
		return nil
	})
	BuntdbExecutor("del b")
	BuntdbExecutor("del a")

	BuntdbExecutor("rwbegin")
	BuntdbExecutor("set x y")
	BuntdbExecutor("set y x")
	BuntdbExecutor("commit")

	bd.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get("x")
		assert.Nil(t, err)
		assert.Equal(t, "y", val)
		val, err = tx.Get("y")
		assert.Nil(t, err)
		assert.Equal(t, "x", val)
		return nil
	})

	BuntdbExecutor("rwbegin")
	BuntdbExecutor("del x")
	BuntdbExecutor("del y")
	BuntdbExecutor("set x xxx")
	BuntdbExecutor("rollback")
	bd.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get("x")
		assert.Nil(t, err)
		assert.Equal(t, "y", val)
		val, err = tx.Get("y")
		assert.Nil(t, err)
		assert.Equal(t, "x", val)
		return nil
	})

	BuntdbExecutor("rbegin")
	BuntdbExecutor("del x")
	BuntdbExecutor("del y")
	BuntdbExecutor("shrink")
	BuntdbExecutor("save testcli_save")
	_, err = os.Lstat("testcli_save")
	assert.Nil(t, err)
	BuntdbExecutor("commit")
	BuntdbExecutor("rollback")
	bd.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get("x")
		assert.Nil(t, err)
		assert.Equal(t, "y", val)
		val, err = tx.Get("y")
		assert.Nil(t, err)
		assert.Equal(t, "x", val)
		return nil
	})
	BuntdbExecutor("shrink")
	BuntdbExecutor("set a xy")
	BuntdbExecutor("save testcli_save")
	_, err = os.Lstat("testcli_save")
	assert.Nil(t, err)
	BuntdbExecutor("save testcli_save")
	BuntdbExecutor("save --force testcli_save")

	BuntdbExecutor("use testcli-2")
	BuntdbExecutor("use -c testcli-2")
	assert.Equal(t, db.GetDbPath(), "testcli-2")
	BuntdbExecutor("use -c testcli")
	assert.Equal(t, db.GetDbPath(), "testcli")
	BuntdbExecutor("use :memory:")
	BuntdbExecutor("exit")
	BuntdbExecutor("")

	os.Remove("testcli_save")
	os.Remove("testcli")
	os.Remove("testcli-2")
}
