package cli

import (
	"fmt"
	"github.com/Sora233/buntdb-cli/db"
	"github.com/alecthomas/kong"
	"github.com/tidwall/buntdb"
	"os"
	"time"
)

var null = "<nil>"

type GetGrammar struct {
	Key          string `arg:"" help:"the key to get"`
	IgnoreExpire bool   `help:"ignore expire"`
}

func (g *GetGrammar) Run(ctx *kong.Context, tx *buntdb.Tx) error {
	val, err := tx.Get(g.Key, g.IgnoreExpire)
	if err != nil {
		if err == buntdb.ErrNotFound {
			fmt.Fprintln(ctx.Stdout, null)
			return nil
		}
		return err
	}
	fmt.Fprintln(ctx.Stdout, val)
	return nil
}

type SetGrammar struct {
	Key   string `arg:"" help:"the key to set"`
	Value string `arg:"" help:"the value assign to the key"`
	TTL   int64  `arg:"" optional:"" help:"expire time in second"`
}

func (s *SetGrammar) Run(ctx *kong.Context, tx *buntdb.Tx) error {
	setOpt := &buntdb.SetOptions{}
	if s.TTL != 0 {
		setOpt.Expires = true
		setOpt.TTL = time.Duration(s.TTL) * time.Second
	}
	oldval, _, err := tx.Set(s.Key, s.Value, setOpt)
	if err != nil {
		return err
	} else {
		if oldval == "" {
			fmt.Fprintln(ctx.Stdout, null)
		} else {
			fmt.Fprintln(ctx.Stdout, oldval)
		}
	}
	return nil
}

type DelGrammar struct {
	Key string `arg:"" help:"the key to delete"`
}

func (d *DelGrammar) Run(ctx *kong.Context, tx *buntdb.Tx) error {
	_, err := tx.Delete(d.Key)
	if err != nil {
		if err == buntdb.ErrNotFound {
			fmt.Fprintln(ctx.Stdout, "0")
			return nil
		}
		return err
	} else {
		fmt.Fprintln(ctx.Stdout, "1")
	}
	return nil
}

type ShowGrammar struct {
	Cmd string `arg:"" enum:"db,index"`
}

func (s *ShowGrammar) Run(ctx *kong.Context, tx *buntdb.Tx) error {
	switch s.Cmd {
	case "db":
		fmt.Fprintln(ctx.Stdout, db.GetDbPath())
	case "index":
		indexes, err := tx.Indexes()
		if err != nil {
			return err
		}
		for _, index := range indexes {
			fmt.Fprintf(ctx.Stdout, "%16v\t", index)
		}
	}
	return nil
}

type KeysGrammar struct {
	Pattern string `arg:"" help:"the match pattern "`
}

func (k KeysGrammar) Run(ctx *kong.Context, tx *buntdb.Tx) error {
	return tx.AscendKeys(k.Pattern, func(key, value string) bool {
		fmt.Fprintf(ctx.Stdout, "%v : %v\n", key, value)
		return true
	})
}

type UseGrammar struct {
	Path   string `arg:"" help:"the new db path"`
	Create bool   `optional:"" short:"c" help:"create new db if path doesn't exist'"`
}

func (u *UseGrammar) Run(ctx *kong.Context, tx *buntdb.Tx) error {
	if tx != nil {
		return db.ErrTransactionExist
	}
	f, err := os.Lstat(u.Path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		if u.Create {
			return db.InitBuntDB(u.Path)
		} else {
			fmt.Fprintf(ctx.Stdout, "%v does not exist, set --create to create it.\n", u.Path)
			return nil
		}
	}
	if f.IsDir() {
		fmt.Fprintf(ctx.Stdout, "%v is a dir.\n", u.Path)
		return nil
	}
	return db.InitBuntDB(u.Path)
}

type TTLGrammar struct {
	Key string `arg:"" help:"the key to show ttl"`
}

func (t *TTLGrammar) Run(ctx *kong.Context, tx *buntdb.Tx) error {
	ttl, err := tx.TTL(t.Key)
	if err != nil {
		if err == buntdb.ErrNotFound {
			fmt.Fprintln(ctx.Stdout, null)
			return nil
		}
		return err
	} else {
		fmt.Fprintln(ctx.Stdout, int64(ttl.Seconds()))
		return nil
	}
}

type RWBeginGrammar struct{}

func (r *RWBeginGrammar) Run(ctx *kong.Context) error {
	_, err := db.Begin(true)
	return err
}

type RBeginGrammar struct{}

func (r *RBeginGrammar) Run(ctx *kong.Context) error {
	_, err := db.Begin(false)
	return err
}

type CommitGrammar struct{}

func (c *CommitGrammar) Run(ctx *kong.Context) error {
	return db.Commit()
}

type RollbackGrammar struct{}

func (r *RollbackGrammar) Run(ctx *kong.Context) error {
	return db.Rollback()
}

type Grammar struct {
	Get      GetGrammar      `cmd:"" help:"get a value from key, return the value if key exists, or <nil> if non-exists."`
	Set      SetGrammar      `cmd:"" help:"set a key-value [ttl], return the old value, or <nil> if old value doesn't exist."`
	Del      DelGrammar      `cmd:"" help:"delete a key, return 1 if success, or 0 if key doesn't exist."`
	Show     ShowGrammar     `cmd:"" help:"show index or db."`
	Keys     KeysGrammar     `cmd:"" help:"iterate over the key match the pattern, support '?' and '*'."`
	Use      UseGrammar      `cmd:"" help:"switch to other db."`
	TTL      TTLGrammar      `cmd:"" help:"get key ttl (seconds), 0 if no ttl, <nil> if key doesn't exist'"`
	RWBegin  RWBeginGrammar  `cmd:"" name:"rwbegin" help:"begin a read/write transaction"`
	RBegin   RBeginGrammar   `cmd:"" name:"rbegin" help:"begin a readonly transaction"`
	Commit   CommitGrammar   `cmd:"" help:"commit a transaction"`
	Rollback RollbackGrammar `cmd:"" help:"rollback a transaction"`
	Exit     bool            `kong:"-"`
}

func (g *Grammar) ExitWrapper(int) {
	g.Exit = true
}

func NewGrammar() *Grammar {
	return &Grammar{}
}
