package main

import (
	"log"
	"os"
)

var cargo string
var candidatos string
var folder string

func main() {
	module := os.Args[1]

	switch module {
	case "bu":
		Bu()
	case "vscmr":
		Vscmr()
	default:
		log.Fatalf("module %s is invalid", module)
	}
}
