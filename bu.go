package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	urna "github.com/mpbertram/urna/ue"
)

func Bu() {
	var function string
	if len(os.Args) > 2 {
		function = os.Args[2]
	}

	switch function {
	case "count":
		countBuFlags()
		countBu()
	case "verify":
		verifyBuFlags()
		verifyBu()
	case "csv":
		csvBuFlags()
		buToCsv()
	default:
		fmt.Println("usage: urna bu <count|verify|csv> <options>")
		fmt.Printf("provided function '%s' is none of (count, verify, csv)\n", function)
	}
}

func buToCsv() {
	cargo := urna.CargoConstitucionalFromString(cargo)
	candidatos := splitCandidatosIntoSlice()

	w := csv.NewWriter(os.Stdout)
	w.Write(append([]string{"UF", "Municipio", "Local", "Secao"}, candidatos...))

	files, err := filepath.Glob(glob)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if strings.HasSuffix(f, ".zip") {
			urna.ProcessZip(f, func(eeg urna.EntidadeEnvelopeGenerico) error {
				bu, err := eeg.ReadBu()
				if err != nil {
					log.Fatal(err)
				}

				votos := urna.CountVotosBu(bu, []urna.CargoConstitucional{cargo})
				var votosForCandidato []string
				for _, candidato := range candidatos {
					votosForCandidato = append(votosForCandidato, fmt.Sprint(votos[cargo][candidato]))
				}
				w.Write(
					append(
						[]string{
							bu.IdentificacaoSecao.Municipio().Uf,
							bu.IdentificacaoSecao.Municipio().Nome,
							fmt.Sprint(bu.IdentificacaoSecao.Local),
							fmt.Sprint(bu.IdentificacaoSecao.Secao),
						},
						votosForCandidato...,
					),
				)

				w.Flush()
				return nil
			})
		}

		if strings.HasSuffix(f, ".bu") {
			entry := urna.BuEntry{Path: f}
			bu, err := entry.ReadBu()
			if err != nil {
				log.Fatal(err)
			}

			votos := urna.CountVotosBu(bu, []urna.CargoConstitucional{cargo})
			var votosForCandidato []string
			for _, candidato := range candidatos {
				votosForCandidato = append(votosForCandidato, fmt.Sprint(votos[cargo][candidato]))
			}
			w.Write(
				append(
					[]string{
						bu.IdentificacaoSecao.Municipio().Uf,
						bu.IdentificacaoSecao.Municipio().Nome,
						fmt.Sprint(bu.IdentificacaoSecao.Local),
						fmt.Sprint(bu.IdentificacaoSecao.Secao),
					},
					votosForCandidato...,
				),
			)
		}
	}

	w.Flush()
}

func countBu() {
	files, err := filepath.Glob(glob)
	if err != nil {
		log.Fatal(err)
	}

	cargos := []urna.CargoConstitucional{urna.CargoConstitucionalFromString(cargo)}
	votos := make(map[urna.CargoConstitucional]map[string]int)

	for _, f := range files {
		if strings.HasSuffix(f, ".zip") {
			urna.ProcessZip(f, func(eeg urna.EntidadeEnvelopeGenerico) error {
				ebu, err := eeg.ReadBu()
				if err != nil {
					log.Fatal(err)
				}
				for cargo, candidato := range urna.CountVotosBu(ebu, cargos) {
					if votos[cargo] == nil {
						votos[cargo] = candidato
					} else {
						for candidato, numVotos := range candidato {
							votos[cargo][candidato] = votos[cargo][candidato] + numVotos
						}
					}
				}
				return nil
			})
		}

		if strings.HasSuffix(f, ".bu") {
			entry := urna.BuEntry{Path: f}
			bu, err := entry.ReadBu()
			if err != nil {
				log.Fatal(err)
			}

			for cargo, candidato := range urna.CountVotosBu(bu, cargos) {
				if votos[cargo] == nil {
					votos[cargo] = candidato
				} else {
					for candidato, numVotos := range candidato {
						votos[cargo][candidato] = votos[cargo][candidato] + numVotos
					}
				}
			}
		}
	}

	log.Println(votos)
}

func verifyBu() {
	files, err := filepath.Glob(glob)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if strings.HasSuffix(f, ".zip") {
			urna.ProcessZip(f, func(eeg urna.EntidadeEnvelopeGenerico) error {
				ebu, err := eeg.ReadBu()
				if err != nil {
					log.Fatal(err)
				}

				err = urna.ValidateVotosBu(ebu)
				if err != nil {
					log.Fatal(err)
				}

				return nil
			})
		}
		if strings.HasSuffix(f, ".bu") {
			entry := urna.BuEntry{Path: f}
			bu, err := entry.ReadBu()
			if err != nil {
				log.Fatal(err)
			}

			err = urna.ValidateVotosBu(bu)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func countBuFlags() {
	countFlags := flag.NewFlagSet("count", flag.ContinueOnError)
	countFlags.StringVar(&cargo, "cargo", "", "e.g. Presidente")
	countFlags.StringVar(&glob, "glob", "", "glob of files to process")

	err := countFlags.Parse(os.Args[3:])
	if err != nil {
		os.Exit(1)
	}

	if len(cargo) == 0 || len(glob) == 0 {
		fmt.Println("usage: urna bu count -glob <glob> -cargo <cargo>")
		countFlags.PrintDefaults()
		os.Exit(1)
	}
}

func verifyBuFlags() {
	verifyFlags := flag.NewFlagSet("verify", flag.ContinueOnError)
	verifyFlags.StringVar(&glob, "glob", "", "glob of files to process")
	err := verifyFlags.Parse(os.Args[3:])
	if err != nil {
		os.Exit(1)
	}

	if len(glob) == 0 {
		fmt.Println("usage: urna bu verify -glob <glob>")
		verifyFlags.PrintDefaults()
		os.Exit(1)
	}
}

func csvBuFlags() {
	csvFlags := flag.NewFlagSet("csv", flag.ContinueOnError)
	csvFlags.StringVar(&cargo, "cargo", "", "e.g. Presidente")
	csvFlags.StringVar(&candidatos, "candidatos", "", "Comma-separated list; e.g. 'Branco,Nulo,99'")
	csvFlags.StringVar(&glob, "glob", "", "glob of files to process")
	err := csvFlags.Parse(os.Args[3:])
	if err != nil {
		os.Exit(1)
	}

	if len(cargo) == 0 || len(glob) == 0 || len(candidatos) == 0 {
		fmt.Println("usage: urna bu csv -cargo <cargo> -candidatos <candidatos> -glob <glob>")
		csvFlags.PrintDefaults()
		os.Exit(1)
	}
}

func splitCandidatosIntoSlice() []string {
	candidatos := strings.Split(candidatos, ",")
	for i := range candidatos {
		candidatos[i] = strings.Trim(candidatos[i], " ")
	}

	return candidatos
}
