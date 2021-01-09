package cli

import (
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
	cmds := []prompt.Suggest{
		{Text: "get", Description: "get command"},
		{Text: "set", Description: "set command"},
		{Text: "del", Description: "del command"},
		{Text: "show", Description: "show info"},
		{Text: "keys", Description: "iterate keys"},
		{Text: "use", Description: "change db"},
		{Text: "exit", Description: "exit buntdb shell client"},
	}
	return prompt.FilterHasPrefix(cmds, cmd, true)
}

func optionCompleter(cmd string, args []string) []prompt.Suggest {
	switch cmd {
	case "get":
	case "set":
	case "del":
	case "show":
		return []prompt.Suggest{
			{Text: "index"},
			{Text: "db"},
		}
	case "keys":
	case "use":
	default:
		return []prompt.Suggest{}
	}
	return []prompt.Suggest{}
}
