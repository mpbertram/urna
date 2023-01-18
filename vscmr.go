package main

import (
	"flag"
	"log"
	"os"

	urna "github.com/mpbertram/urna/ue"
)

func Vscmr() {
	var function = os.Args[2]
	switch function {
	case "verify":
		verifyVscmrFlags()
		verifyVscmr()
	default:
		log.Fatalf("function %s is invalid", function)
	}
}

func verifyVscmr() {
	urna.VerifyAssinaturas(folder)
}

func verifyVscmrFlags() {
	verifyFlags := flag.NewFlagSet("verify", flag.ContinueOnError)
	verifyFlags.StringVar(&folder, "folder", ".", "<folder>")
	verifyFlags.Parse(os.Args[3:])
}
