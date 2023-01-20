package main

import (
	"fmt"
	"os"
)

var cargo string
var candidatos string

func main() {
	var module string
	if len(os.Args) > 1 {
		module = os.Args[1]
	}

	switch module {
	case "bu":
		Bu()
	case "vscmr":
		Vscmr()
	default:
		fmt.Println("usage: urna <bu|vscmr> <function> <options>")
		fmt.Printf("provided module '%s' is none of (bu, vscmr)\n", module)
	}
}
