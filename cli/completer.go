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
		return cmdCompleter(args[0])
	} else {
		return optionCompleter(args[0], args[1:])
	}
}

func cmdCompleter(cmd string) []prompt.Suggest {
	if Debug {
		fmt.Printf("cmdCompleter %v\n", cmd)
	}
	cmds := []prompt.Suggest{
		{Text: "get", Description: "get command"},
		{Text: "set", Description: "set command"},
		{Text: "del", Description: "del command"},
		{Text: "ttl", Description: "ttl command"},
		{Text: "show", Description: "show info"},
		{Text: "keys", Description: "iterate keys"},
		{Text: "use", Description: "change db"},
		{Text: "exit", Description: "exit buntdb shell client"},
	}
	tx, _ := db.GetCurrentTransaction()
	if tx == nil {
		cmds = append(cmds,
			prompt.Suggest{Text: "rbegin", Description: "open a readonly transaction"},
			prompt.Suggest{Text: "rwbegin", Description: "open a read/write transaction"},
		)
	} else {
		cmds = append(cmds,
			prompt.Suggest{Text: "rollback", Description: "rollback a transaction"},
			prompt.Suggest{Text: "commit", Description: "commit a transaction"},
		)
	}
	return prompt.FilterHasPrefix(cmds, cmd, true)
}

func optionCompleter(cmd string, args []string) []prompt.Suggest {
	if Debug {
		fmt.Printf("optionCompleter %v %v\n", cmd, args)
	}
	switch cmd {
	case "get":
	case "set":
	case "del":
	case "ttl":
	case "show":
		return []prompt.Suggest{
			{Text: "index"},
			{Text: "db"},
		}
	case "keys":
	case "use":
	case "rbegin":
	case "rwbegin":
	case "rollback":
	case "commit":
	default:
		return []prompt.Suggest{}
	}
	return []prompt.Suggest{}
}