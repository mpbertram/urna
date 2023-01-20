package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

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
	files, err := filepath.Glob(glob)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if strings.HasSuffix(f, ".zip") {
			urna.VerifyAssinaturaZip(f)
		}

		if strings.HasSuffix(f, ".vscmr") {
			urna.VerifyAssinaturaVscmr(f)
		}
	}
}

func verifyVscmrFlags() {
	verifyFlags := flag.NewFlagSet("verify", flag.ContinueOnError)
	verifyFlags.StringVar(&glob, "glob", "", "glob of files to process")
	err := verifyFlags.Parse(os.Args[3:])
	if err != nil {
		os.Exit(1)
	}

	if len(glob) == 0 {
		fmt.Println("usage: urna vscmr verify -glob <glob>")
		verifyFlags.PrintDefaults()
		os.Exit(1)
	}
}
