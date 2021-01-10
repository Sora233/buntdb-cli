package cli

import (
	"github.com/Sora233/buntdb-cli/db"
	"github.com/c-bata/go-prompt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuntdbCompleter(t *testing.T) {
	Debug = true
	buf := prompt.NewBuffer()
	d := buf.Document()
	assert.NotNil(t, d)
	sug := BuntdbCompleter(*d)
	assert.Empty(t, sug)
	buf.InsertText("g", false, true)
	d = buf.Document()
	assert.NotNil(t, d)
	sug = BuntdbCompleter(*d)
	assert.Len(t, sug, 1)
	assert.Equal(t, sug[0].Text, "get")

	buf.DeleteBeforeCursor(999)
	buf.InsertText("s", false, true)
	d = buf.Document()
	assert.NotNil(t, d)
	sug = BuntdbCompleter(*d)
	assert.Len(t, sug, 2)
	text := []string{sug[0].Text, sug[1].Text}
	assert.Contains(t, text, "set")
	assert.Contains(t, text, "show")

	buf.DeleteBeforeCursor(999)
	buf.InsertText("get ", false, true)
	d = buf.Document()
	assert.NotNil(t, d)
	sug = BuntdbCompleter(*d)
	assert.Empty(t, sug)

	buf.DeleteBeforeCursor(999)
	buf.InsertText("show ", false, true)
	d = buf.Document()
	assert.NotNil(t, d)
	sug = BuntdbCompleter(*d)
	assert.Len(t, sug, 2)
	text = []string{sug[0].Text, sug[1].Text}
	assert.Contains(t, text, "db")
	assert.Contains(t, text, "index")

	buf.DeleteBeforeCursor(999)
	buf.InsertText("fake ", false, true)
	d = buf.Document()
	assert.NotNil(t, d)
	sug = BuntdbCompleter(*d)
	assert.Empty(t, sug)

	buf.DeleteBeforeCursor(999)

	db.InitBuntDB(":memory:")
	db.Begin(true)
	buf.InsertText("r", false, true)
	d = buf.Document()
	assert.NotNil(t, d)
	sug = BuntdbCompleter(*d)
	assert.Len(t, sug, 1)
	assert.Equal(t, "rollback", sug[0].Text)
	db.Rollback()
	db.Close()

}
