package main

import (
	"github.com/Sora233/buntdb-cli/cli"
	"github.com/Sora233/buntdb-cli/db"
	"github.com/alecthomas/kong"
	"github.com/c-bata/go-prompt"
	"os"
	"path"
)

var CLI struct {
	Path    string  `arg:"" optional:"" help:"buntudb file path, default a tempfile"`
	Debug   bool    `optional:"" help:"enable debug output"`
	Version version `optional:"" short:"v" help:"print version"`
}

func main() {
	kong.Parse(&CLI, kong.UsageOnError(), kong.Name("buntdb-cli"))
	defer os.Exit(0)

	if CLI.Debug {
		cli.Debug = true
	}

	if CLI.Path == "" {
		CLI.Path = db.GetTempDbPath("buntdb-cli")
	}

	db.InitBuntDB(CLI.Path)
	defer db.Close()

	p := prompt.New(
		cli.BuntdbExecutor,
		cli.BuntdbCompleter,
		prompt.OptionTitle("buntdb-cli: an interactive buntdb shell client"),
		prompt.OptionPrefix(path.Base(CLI.Path)+"> "),
	)
	p.Run()
}
