package cli

import (
	"fmt"
	"github.com/Sora233/buntdb-cli/db"
	"github.com/c-bata/go-prompt"
	"strings"
)

func BuntdbCompleter(d prompt.Document) []prompt.Suggest {
	if d.TextBeforeCursor() == "" {
		return []prompt.Suggest{}
	}
	args := strings.Split(d.TextBeforeCursor(), " ")
	if len(args) == 1 {
		// input command
		return cmdCompleter(d, args[0])
	} else {
		return optionCompleter(d, args[0], args[1:])
	}
}

func cmdCompleter(d prompt.Document, cmd string) []prompt.Suggest {
	if Debug {
		fmt.Printf("|cmdCompleter %v|\n", cmd)
	}
	cmds := []prompt.Suggest{
		{Text: "get", Description: "get command"},
		{Text: "set", Description: "set command"},
		{Text: "del", Description: "del command"},
		{Text: "ttl", Description: "ttl command"},
		{Text: "show", Description: "show info"},
		{Text: "keys", Description: "iterate keys"},
		{Text: "search", Description: "search for string in all values"},
		{Text: "use", Description: "change db"},
		{Text: "exit", Description: "exit buntdb shell client"},
		{Text: "drop", Description: "drop the index"},
		{Text: "help", Description: "show available commands"},
	}
	tx, _ := db.GetCurrentTransaction()
	if tx == nil {
		cmds = append(cmds,
			prompt.Suggest{Text: "rbegin", Description: "open a readonly transaction"},
			prompt.Suggest{Text: "rwbegin", Description: "open a read/write transaction"},
			prompt.Suggest{Text: "shrink", Description: "shrink command"},
			prompt.Suggest{Text: "save", Description: "save db to file"},
		)
	} else {
		cmds = append(cmds,
			prompt.Suggest{Text: "rollback", Description: "rollback a transaction"},
			prompt.Suggest{Text: "commit", Description: "commit a transaction"},
		)
	}
	return prompt.FilterHasPrefix(cmds, cmd, true)
}

func optionCompleter(d prompt.Document, cmd string, args []string) []prompt.Suggest {
	if Debug {
		fmt.Printf("|optionCompleter %v [%v]|\n", cmd, strings.Join(args, ":"))
	}
	var result = make([]prompt.Suggest, 0)
	switch cmd {
	case "get":
	case "set":
	case "del":
	case "ttl":
	case "show":
		result = []prompt.Suggest{
			{Text: "index"},
			{Text: "db"},
		}
		if len(args) == 0 {
			break
		}
		arg := args[0]
		if Debug {
			fmt.Printf("|arg %v|", arg)
		}
		switch arg {
		case "index":
			result = []prompt.Suggest{}
		case "db":
			result = []prompt.Suggest{}
		default:
			result = prompt.FilterHasPrefix(result, arg, true)
		}
	case "keys":
	case "use":
	case "rbegin":
	case "rwbegin":
	case "rollback":
	case "commit":
	case "shrink":
	case "save":
	case "drop":
		result = []prompt.Suggest{
			{Text: "index"},
		}
		if len(args) == 0 {
			break
		}
		switch args[0] {
		case "index":
			result = []prompt.Suggest{}
			tx, _, closeTx := db.GetCurrentOrNewTransaction()
			defer closeTx()
			indexes, err := tx.Indexes()
			if err == nil {
				for _, index := range indexes {
					result = append(result, prompt.Suggest{Text: index})
				}
				if len(args) >= 2 {
					result = prompt.FilterHasPrefix(result, args[1], true)
				}
			}
		default:
			result = prompt.FilterHasPrefix(result, args[0], true)
		}
	default:
	}
	return result
}
