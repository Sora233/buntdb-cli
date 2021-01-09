package main

import (
	"github.com/Sora233/buntdb-cli/db"
	"github.com/alecthomas/kong"
	"github.com/c-bata/go-prompt"
	"io/ioutil"
	"os"
	"path"
)

var CLI struct {
	Path    string  `arg:"" optional:"" help:"buntudb file path, default a tempfile"`
	Version version `optional:"" short:"v" help:"print version"`
}

func main() {
	kong.Parse(&CLI, kong.UsageOnError(), kong.Name("buntdb-cli"))
	defer os.Exit(0)

	if CLI.Path == "" {
		f, err := ioutil.TempFile("", "buntdb-cli")
		if err != nil {
		}
		f.Close()
		CLI.Path = f.Name()
		defer os.Remove(f.Name())
	}
	db.InitBuntDB(CLI.Path)
	defer db.Close()

	p := prompt.New(
		buntdbExecutor,
		buntdbCompleter,
		prompt.OptionTitle("buntdb-cli: an interactive buntdb shell client"),
		prompt.OptionPrefix(path.Base(CLI.Path)+">"),
	)
	p.Run()
}
