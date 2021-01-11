package cli

import (
	"fmt"
	"github.com/Sora233/buntdb-cli/db"
	"github.com/alecthomas/kong"
	"strings"
)

type transactionRequireType int64

const (
	noNeed transactionRequireType = iota
	any
	nonNil
)

func commandTransactionRequireType(command string) transactionRequireType {
	switch command {
	case "rwbegin", "rbegin", "rollback", "commit", "shrink", "save":
		return noNeed
	case "use":
		return any
	default:
		return nonNil
	}
}

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
	switch commandTransactionRequireType(cmd) {
	case noNeed:
		err = ctx.Run()
		if err != nil {
			fmt.Printf("ERR: %v\n", err)
		}
		return
	case any:
		tx, _ := db.GetCurrentTransaction()
		if cmd == "use" {
			err = ctx.Run(tx)
			if err != nil {
				fmt.Printf("ERR: %v\n", err)
			}
			return
		}
	case nonNil:
		tx, rw, closeTx := db.GetCurrentOrNewTransaction()
		if Debug {
			fmt.Printf("GetCurrentOrNewTransaction(%v)\n", db.RWDescribe(rw))
		}
		defer func() {
			err = closeTx()
			if err != nil {
				fmt.Printf("ERR: %v\n", err)
			}
		}()
		err = ctx.Run(tx)
		if err != nil {
			fmt.Printf("ERR: %v\n", err)
			return
		}
	default:
		fmt.Printf("ERR: unknown transaction require\n")
		return
	}
}
