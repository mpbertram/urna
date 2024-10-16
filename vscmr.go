package main

import (
	"encoding/csv"
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
		verifyVscmr(GetFlags())
	case "csv":
		vscmrToCsv(GetFlags())
	case "cert":
		parseCerts(GetFlags())
	case "export":
		exportCerts(GetFlags())
	default:
		fmt.Println("usage: urna vscmr <verify|csv|cert|export> <file_1> ... <file_n>")
	}
}

func exportCerts(files []string) {
	for _, f := range files {
		log.Printf("processing file %s", f)

		if strings.HasSuffix(f, ".zip") {
			urna.ExportCertsZip(f)
		}

		if strings.HasSuffix(f, ".vscmr") {
			urna.ExportCertsVscmr(f)
		}
	}
}

func parseCerts(files []string) {
	for _, f := range files {
		log.Printf("processing file %s", f)

		if strings.HasSuffix(f, ".zip") {
			results := urna.VerifyCertsZip(f)
			if len(results) > 0 {
				for _, r := range results {
					print(r)
				}
			}
		}

		if strings.HasSuffix(f, ".vscmr") {
			results := urna.VerifyCertsVscmr(f)
			if len(results) > 0 {
				for _, r := range results {
					print(r)
				}
			}
		}
	}
}

func vscmrToCsv(files []string) {
	w := csv.NewWriter(os.Stdout)
	w.Write(
		[]string{
			"Municipio",
			"Zona",
			"Secao",
			"Arquivo",
			"Tipo",
			"Status",
			"Erro"})

	for _, f := range files {
		log.Printf("processing file %s", f)

		if strings.HasSuffix(f, ".zip") {
			result := urna.VerifyAssinaturaZip(f)
			for _, r := range result {
				writeToCsv(r, w)
			}
		}

		if strings.HasSuffix(f, ".vscmr") {
			result := urna.VerifyAssinaturaVscmr(f)
			for _, r := range result {
				writeToCsv(r, w)
			}
		}
	}

	w.Flush()
}

func verifyVscmr(files []string) {
	for _, f := range files {
		log.Printf("processing file %s", f)

		if strings.HasSuffix(f, ".zip") {
			result := urna.VerifyAssinaturaZip(f)
			for _, r := range result {
				print(r)
			}
		}

		if strings.HasSuffix(f, ".vscmr") {
			result := urna.VerifyAssinaturaVscmr(f)
			for _, r := range result {
				print(r)
			}
		}
	}
}

func print(r urna.VerificationResult) {
	log.Println(r.Msg())
	if r.Err != nil {
		log.Println(r.Err)
	}
}

func writeToCsv(r urna.VerificationResult, w *csv.Writer) {
	if r.Err != nil {
		w.Write([]string{
			r.Municipio,
			r.Zona,
			r.Secao,
			r.Filename,
			r.Type.String(),
			r.Ok.String(),
			r.Err.Error()})
	} else {
		w.Write([]string{
			r.Municipio,
			r.Zona,
			r.Secao,
			r.Filename,
			r.Type.String(),
			r.Ok.String(),
			""})
	}
}

func GetFlags() []string {
	if len(os.Args) > 3 {
		return os.Args[3:]
	}

	fmt.Println("usage: urna vscmr <verify|csv|cert> <file_1> ... <file_n>")
	os.Exit(1)

	return []string{}
}
