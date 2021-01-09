package main

import (
	"fmt"
	"runtime"
	"strings"
)

func buntdbExecutor(s string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return
	} else if s == "quit" || s == "exit" {
		fmt.Println("bye")
		runtime.Goexit()
	}
	fmt.Println("got", s)
}
