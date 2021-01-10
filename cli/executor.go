package cli

import (
	"fmt"
	"github.com/Sora233/buntdb-cli/db"
	"github.com/alecthomas/kong"
	"strings"
)

func BuntdbExecutor(s string) {
	s = strings.TrimSpace(s)
	if s == "" || s == "exit" {
		return
	}
	args := ArgSplit(s)

	if Debug {
		fmt.Printf("args: %v\n", strings.Join(args, "/"))
	}

	grammar := NewGrammar()
	k := kong.Must(
		grammar,
		kong.Exit(grammar.ExitWrapper),
	)
	ctx, err := k.Parse(args)
	if grammar.Exit {
		return
	}
	if err != nil {
		fmt.Printf("ERR: %v\n", err)
		return
	}
	cmd := ctx.Selected().Name
	if cmd == "rwbegin" || cmd == "rbegin" || cmd == "rollback" || cmd == "commit" {
		err = ctx.Run()
		if err != nil {
			fmt.Printf("ERR: %v\n", err)
		}
		return
	}
	tx, rw := db.GetCurrentTransaction()
	if ctx.Selected().Name == "use" {
		err = ctx.Run(tx)
		if err != nil {
			fmt.Printf("ERR: %v\n", err)
		}
		return
	}
	if tx != nil {
		if Debug {
			fmt.Printf("got current %v transaction\n", db.RWDescribe(rw))
		}
		err = ctx.Run(tx)
		if err != nil {
			fmt.Printf("ERR: %v\n", err)
			return
		}
	} else {
		if Debug {
			fmt.Printf("no transaction, create a rw transaction\n")
		}
		tx, err := db.Begin(true)
		if err != nil {
			fmt.Printf("ERR: %v\n", err)
			return
		}
		defer func() {
			if Debug {
				fmt.Printf("transaction commit\n")
			}
			err := db.Commit()
			if err != nil {
				fmt.Printf("ERR: commit error %v\n", err)
			}
		}()
		err = ctx.Run(tx)
		if err != nil {
			fmt.Printf("ERR: %v\n", err)
			return
		}
	}
}
