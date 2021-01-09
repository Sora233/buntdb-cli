package main

import (
	"github.com/c-bata/go-prompt"
	"strings"
)

func buntdbCompleter(d prompt.Document) []prompt.Suggest {
	if d.TextBeforeCursor() == "" {
		return []prompt.Suggest{}
	}
	args := strings.Split(d.TextBeforeCursor(), " ")
	if len(args) == 1 {
		// input command
		return cmdCompleter(args[0], true)
	} else {
		return optionCompleter(args[0], args[1:])
	}
}

func cmdCompleter(cmd string, help bool) []prompt.Suggest {
	cmds := []prompt.Suggest{
		{Text: "get", Description: "get command"},
		{Text: "set", Description: "set command"},
		{Text: "del", Description: "del command"},
		{Text: "show", Description: "show info"},
		{Text: "scan", Description: "iterate keys"},
		{Text: "use", Description: "change db"},
	}
	if help {
		cmds = append(cmds, prompt.Suggest{Text: "help", Description: "show help for command"})
	}
	return prompt.FilterHasPrefix(cmds, cmd, true)
}

func optionCompleter(cmd string, args []string) []prompt.Suggest {
	last := args[len(args)-1]
	switch cmd {
	case "get":
	case "set":
	case "del":
	case "show":
		return []prompt.Suggest{
			{Text: "index"},
			{Text: "db"},
		}
	case "scan":
	case "help":
		return cmdCompleter(last, false)
	default:
		return []prompt.Suggest{}
	}
	return []prompt.Suggest{}
}
