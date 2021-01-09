package cli

import (
	"fmt"
	"github.com/alecthomas/kong"
	"runtime"
	"strings"
)

func BuntdbExecutor(s string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return
	} else if s == "exit" {
		fmt.Println("bye")
		runtime.Goexit()
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
	err = ctx.Run()
	if err != nil {
		fmt.Printf("ERR: %v\n", err)
	}
}
