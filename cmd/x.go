package main

import (
	"fmt"
	"github.com/laracro/x/internal/engine"
	"os"
)

var (
	command string
	options []string
)

func main() {

	x := engine.LoadEngine()

	if len(os.Args) < 2 {
		x.FetchAll()
		return
	}

	command = os.Args[1]
	options = os.Args[2:]

	fmt.Printf("command{%s} | options{%s}\n", command, options)
}
