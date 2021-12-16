package cli

import (
	"fmt"
	"github.com/Sora233/buntdb-cli/db"
	"github.com/alecthomas/kong"
	"github.com/tidwall/buntdb"
	"os"
	"strings"
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

type SearchGrammar struct {
	Pattern string `arg:"" help:"the match pattern "`
	Delete  bool   `optional:"" short:"c" help:"delete all keys found (DANGEROUS)"`
}

func (s SearchGrammar) Run(ctx *kong.Context, tx *buntdb.Tx) error {
	var keys []string
	err := tx.AscendKeys("*", func(key, value string) bool {
		if strings.Contains(value, s.Pattern) {
			keys = append(keys, key)
			fmt.Fprintf(ctx.Stdout, "%v : %v\n", key, value)
		}
		return true
	})

	if err != nil {
		return err
	}

	if !s.Delete {
		return nil
	}

	fmt.Fprintf(ctx.Stdout, "\n\n-------------------\n\nFound %d keys, written above. \nSleeping 10 seconds before we delete.\n", len(keys))
	for n := 0; n != 10; n++ {
		time.Sleep(1 * time.Second)
		fmt.Fprintf(ctx.Stdout, ".")
	}

	for i, k := range keys {
		if i == 1 {
			fmt.Fprint(ctx.Stdout, "\n")
		}
		_, err := tx.Delete(k)
		if err == nil {
			fmt.Fprintf(ctx.Stdout, "Deleted: %v\n", k)
		}
	}
	return nil
}

type UseGrammar struct {
	Path   string `arg:"" help:"the new db path"`
	Create bool   `optional:"" short:"c" help:"create new db if path doesn't exist'"`
}

func (u *UseGrammar) Run(ctx *kong.Context, tx *buntdb.Tx) error {
	if tx != nil {
		return db.ErrTransactionExist
	}
	if u.Path == ":memory:" {
		return db.InitBuntDB(u.Path)
	}
	f, err := os.Lstat(u.Path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		if u.Create {
			return db.InitBuntDB(u.Path)
		} else {
			return fmt.Errorf("%v does not exist, set --create to create it", u.Path)
		}
	}
	if f.IsDir() {
		return fmt.Errorf("%v is a dir", u.Path)
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
	}
	fmt.Fprintln(ctx.Stdout, int64(ttl.Seconds()))
	return nil
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

type ShrinkGrammar struct{}

func (s *ShrinkGrammar) Run(ctx *kong.Context) error {
	return db.Shrink()
}

type SaveGrammar struct {
	Path  string `arg:"" help:"the path to save"`
	Force bool   `optional:"" help:"overwrite if the path exists"`
}

func (s *SaveGrammar) Run(ctx *kong.Context) error {
	f, err := os.Lstat(s.Path)
	if err == nil {
		if f.IsDir() {
			return fmt.Errorf("%v is a dir", s.Path)
		}
		if !s.Force {
			return fmt.Errorf("%v exist, use --force to overwrite it", s.Path)
		}
	}
	file, err := os.Create(s.Path)
	if err != nil {
		return err
	}
	return db.Save(file)
}

type DropGrammar struct {
	Index DropIndexGrammar `cmd:"" help:"drop the index with the given name"`
}

type DropIndexGrammar struct {
	Name string `arg:"" help:"the index name to drop"`
}

func (s *DropIndexGrammar) Run(ctx *kong.Context, tx *buntdb.Tx) error {
	return tx.DropIndex(s.Name)
}

type HelpGrammar struct {
	//
}

func (h *HelpGrammar) Run(ctx *kong.Context) (err error) {
	commands := "get\nset\ndel\nttl\nrbegin (begin a readonly transaction)\nrwbegin (begin a read/write transaction)\ncommit\nrollback\nshow\nkeys\nsearch\nuse\nshrink"
	_, err = fmt.Fprintln(ctx.Stdout,
		"Commands available:\n----------\n"+commands+"\n----------\nFor more help, try running a command with the -h switch.",
	)
	return
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
	Shrink   ShrinkGrammar   `cmd:"" help:"run database shrink command"`
	Save     SaveGrammar     `cmd:"" help:"save the db to file"`
	Drop     DropGrammar     `cmd:"" help:"drop command"`
	Search   SearchGrammar   `cmd:"" help:"Search for a string contained in any values"`
	Help     HelpGrammar     `cmd:""`
	Exit     bool            `kong:"-"`
}

func (g *Grammar) ExitWrapper(int) {
	g.Exit = true
}

func NewGrammar() *Grammar {
	return &Grammar{}
}
