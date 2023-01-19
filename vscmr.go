package main

import (
	"flag"
	"fmt"
	"os"

	urna "github.com/mpbertram/urna/ue"
)

func Vscmr() {
	var function string
	if len(os.Args) > 2 {
		function = os.Args[2]
	}

	switch function {
	case "verify":
		verifyVscmrFlags()
		verifyVscmr()
	default:
		fmt.Println("usage: urna vscmr verify <options>")
	}
}

func verifyVscmr() {
	urna.VerifyAssinaturas(folder)
}

func verifyVscmrFlags() {
	verifyFlags := flag.NewFlagSet("verify", flag.ContinueOnError)
	verifyFlags.StringVar(&folder, "folder", ".", "folder to search for *.vscmr files")
	err := verifyFlags.Parse(os.Args[3:])
	if err != nil {
		fmt.Println("usage: urna vscmr verify -folder <folder>")
		verifyFlags.PrintDefaults()
		os.Exit(1)
	}

	if len(folder) == 0 {
		fmt.Println("usage: urna vscmr verify -folder <folder>")
		verifyFlags.PrintDefaults()
		os.Exit(1)
	}
}
