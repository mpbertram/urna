package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
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
		countBu(filesToProcess())
	case "verify":
		verifyBu(filesToProcess())
	case "csv":
		csvBuFlags()
		buToCsv(filesToProcess())
	default:
		fmt.Println("usage: urna bu <count|verify|csv> <options>")
		fmt.Printf("provided function '%s' is none of (count, verify, csv)\n", function)
	}
}

func buToCsv(files []string) {
	cargo := urna.CargoConstitucionalFromString(cargo)
	candidatos := splitCandidatosIntoSlice()

	w := csv.NewWriter(os.Stdout)
	w.Write(append([]string{"UF", "Municipio", "Local", "Secao"}, candidatos...))

	for _, f := range files {
		log.Printf("processing file %s", f)

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

func countBu(files []string) {
	cargos := []urna.CargoConstitucional{urna.CargoConstitucionalFromString(cargo)}
	votos := make(map[urna.CargoConstitucional]map[string]int)

	for _, f := range files {
		log.Printf("processing file %s", f)

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

func verifyBu(files []string) {
	for _, f := range files {
		log.Printf("processing file %s", f)

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

	err := countFlags.Parse(os.Args[3:])
	if err != nil {
		os.Exit(1)
	}

	if len(cargo) == 0 {
		fmt.Println("usage: urna bu count -glob <glob> -cargo <cargo>")
		countFlags.PrintDefaults()
		os.Exit(1)
	}
}

func csvBuFlags() {
	csvFlags := flag.NewFlagSet("csv", flag.ContinueOnError)
	csvFlags.StringVar(&cargo, "cargo", "", "e.g. Presidente")
	csvFlags.StringVar(&candidatos, "candidatos", "", "Comma-separated list; e.g. 'Branco,Nulo,99'")
	err := csvFlags.Parse(os.Args[3:])
	if err != nil {
		os.Exit(1)
	}

	if len(cargo) == 0 || len(candidatos) == 0 {
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
