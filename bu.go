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
	var function = os.Args[2]
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
		log.Fatalf("function %s is invalid", function)
	}
}

func buToCsv() {
	cargo := urna.CargoConstitucionalFromString(cargo)
	candidatos := splitCandidatosIntoSlice()

	w := csv.NewWriter(os.Stdout)
	w.Write(append([]string{"UF", "Municipio", "Local", "Secao"}, candidatos...))

	bus, err := urna.ReadAllBu(folder)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range bus {
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
	w.Flush()

	urna.ProcessAllZip(folder, func(eeg urna.EntidadeEnvelopeGenerico) error {
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

func countBu() {
	bus, err := urna.ReadAllBu(folder)
	if err != nil {
		log.Fatal(err)
	}

	cargos := []urna.CargoConstitucional{urna.CargoConstitucionalFromString(cargo)}
	votos := urna.CountVotos(bus, cargos)

	urna.ProcessAllZip(folder, func(eeg urna.EntidadeEnvelopeGenerico) error {
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

	log.Println(votos)
}

func verifyBu() {
	bus, err := urna.ReadAllBu(folder)
	if err != nil {
		log.Fatal(err)
	}

	err = urna.ValidateVotos(bus)
	if err != nil {
		log.Fatal(err)
	}

	urna.ProcessAllZip(folder, func(eeg urna.EntidadeEnvelopeGenerico) error {
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

func countBuFlags() {
	countFlags := flag.NewFlagSet("count", flag.ContinueOnError)
	countFlags.StringVar(&cargo, "cargo", "Presidente", "<Presidente|Governador>")
	countFlags.StringVar(&folder, "folder", ".", "<folder>")
	countFlags.Parse(os.Args[3:])
}

func verifyBuFlags() {
	verifyFlags := flag.NewFlagSet("verify", flag.ContinueOnError)
	verifyFlags.StringVar(&folder, "folder", ".", "<folder>")
	verifyFlags.Parse(os.Args[3:])
}

func csvBuFlags() {
	csvFlags := flag.NewFlagSet("csv", flag.ContinueOnError)
	csvFlags.StringVar(&cargo, "cargo", "Presidente", "<Presidente|Governador>")
	csvFlags.StringVar(&candidatos, "candidatos", "", "e.g. 'Branco,Nulo,99'")
	csvFlags.StringVar(&folder, "folder", ".", "<folder>")
	csvFlags.Parse(os.Args[3:])
}

func splitCandidatosIntoSlice() []string {
	candidatos := strings.Split(candidatos, ",")
	for i := range candidatos {
		candidatos[i] = strings.Trim(candidatos[i], " ")
	}

	return candidatos
}
