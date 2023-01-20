package main

import (
	"fmt"
	"log"
	"os"
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
		verifyVscmr(verifyVscmrFlags())
	default:
		fmt.Println("usage: urna vscmr verify <file_1> ... <file_n>")
	}
}

func verifyVscmr(files []string) {
	for _, f := range files {
		log.Printf("processing file %s", f)

		if strings.HasSuffix(f, ".zip") {
			urna.VerifyAssinaturaZip(f)
		}

		if strings.HasSuffix(f, ".vscmr") {
			urna.VerifyAssinaturaVscmr(f)
		}
	}
}

func verifyVscmrFlags() []string {
	if len(os.Args) > 3 {
		return os.Args[3:]
	}

	fmt.Println("usage: urna vscmr verify <file_1> ... <file_n>")
	os.Exit(1)

	return []string{}
}
