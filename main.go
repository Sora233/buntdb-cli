package main

import (
	"fmt"
	"github.com/Sora233/buntdb-cli/cli"
	"github.com/Sora233/buntdb-cli/db"
	"github.com/alecthomas/kong"
	"github.com/c-bata/go-prompt"
	"os"
	"path"
	"strings"
)

var CLI struct {
	Path    string  `arg:"" optional:"" help:"buntudb file path, default a tempfile"`
	Debug   bool    `optional:"" help:"enable debug output"`
	Version version `optional:"" short:"v" help:"print version"`
}

func main() {
	kong.Parse(&CLI, kong.UsageOnError(), kong.Name("buntdb-cli"))

	if CLI.Debug {
		cli.Debug = true
	}

	if CLI.Path == "" {
		CLI.Path = db.GetTempDbPath("buntdb-cli")
	}

	err := db.InitBuntDB(CLI.Path)
	if err != nil {
		fmt.Printf("ERR: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	p := prompt.New(
		cli.BuntdbExecutor,
		cli.BuntdbCompleter,
		prompt.OptionTitle("buntdb-cli: an interactive buntdb shell client"),
		prompt.OptionLivePrefix(func() (prefix string, useLivePrefix bool) {
			tx, rw := db.GetCurrentTransaction()
			if tx != nil {
				return path.Base(db.GetDbPath()) + fmt.Sprintf("(%v)", db.RWDescribe(rw)) + "> ", true
			} else {
				return path.Base(db.GetDbPath()) + "> ", true
			}
		}),
		prompt.OptionSetExitCheckerOnInput(func(in string, breakline bool) bool {
			return strings.TrimSpace(in) == "exit" && breakline
		}),
	)
	p.Run()
	tx, _ := db.GetCurrentTransaction()
	if tx != nil {
		fmt.Println("WARN: current transaction will rollback")
		db.Rollback()
	}
	fmt.Println("bye")
}
