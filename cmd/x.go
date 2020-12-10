package main

import (
	"fmt"
	"github.com/laracro/x/internal/client"
	"github.com/laracro/x/internal/engine"
	"log"
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
	}
	command = os.Args[1]
	switch command {
	case "web":
		fmt.Printf("command -> %s\n",command)
		break
	case "conn":
		options = os.Args[2:]
		x.GetInstance(options[0])
		c := client.NewClient(x.Instance)
		c.Login()
		break
	default:
		log.Fatal("unknown command")
	}
}