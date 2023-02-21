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
		countBu(countBuFlags())
	case "verify":
		verifyBu(verifyBuFlags())
	case "csv":
		buToCsv(csvBuFlags())
	default:
		fmt.Println("usage: urna bu <count|verify|csv> <options>")
		fmt.Printf("provided function '%s' is none of (count, verify, csv)\n", function)
	}
}

func buToCsv(files []string) {
	cargo := urna.CargoConstitucionalFromString(cargo)
	candidatos := splitCandidatosIntoSlice()

	w := csv.NewWriter(os.Stdout)
	w.Write(
		append(
			[]string{
				"UF",
				"Municipio",
				"Tipo urna",
				"Tipo arquivo",
				"Apuracao (tipo)",
				"Apuracao (motivo)",
				"Zona",
				"Local",
				"Secao"}, candidatos...))

	for _, f := range files {
		log.Printf("processing file %s", f)

		if strings.HasSuffix(f, ".zip") {
			urna.ProcessZip(f, func(eeg urna.EntidadeEnvelopeGenerico) error {
				bu, err := eeg.ReadBu()
				if err != nil {
					log.Fatal(err)
				}

				w.Write(countVotos(bu, cargo, candidatos))
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

			w.Write(countVotos(bu, cargo, candidatos))
		}
	}

	w.Flush()
}

func countVotos(bu urna.EntidadeBoletimUrna, cargo urna.CargoConstitucional, candidatos []string) []string {
	votos := urna.CountVotosBu(bu, []urna.CargoConstitucional{cargo})
	var votosForCandidato []string
	for _, candidato := range candidatos {
		votosForCandidato = append(votosForCandidato, fmt.Sprint(votos[cargo][candidato]))
	}

	apuracao, err := bu.Urna.ReadMotivoUtilizacaoSA()
	if err != nil {
		log.Println(err)
	}

	return append(
		[]string{
			bu.IdentificacaoSecao.Municipio().Uf,
			bu.IdentificacaoSecao.Municipio().Nome,
			bu.Urna.Tipo().String(),
			bu.Urna.TipoDeArquivo().String(),
			apuracao.Tipo().String(),
			apuracao.Motivo(),
			fmt.Sprint(bu.IdentificacaoSecao.MunicipioZona.Zona),
			fmt.Sprint(bu.IdentificacaoSecao.Local),
			fmt.Sprint(bu.IdentificacaoSecao.Secao),
		}, votosForCandidato...)
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

				results := urna.ValidateVotosBu(ebu)
				for _, r := range results {
					log.Println(r.Msg)
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

			results := urna.ValidateVotosBu(bu)
			for _, r := range results {
				log.Println(r.Msg)
			}
		}
	}
}

func countBuFlags() []string {
	countFlags := flag.NewFlagSet("count", flag.ContinueOnError)
	countFlags.StringVar(&cargo, "cargo", "", "e.g. Presidente")

	err := countFlags.Parse(os.Args[3:])
	if err != nil {
		os.Exit(1)
	}

	if len(cargo) == 0 {
		fmt.Println("usage: urna bu count -cargo <cargo> <file_1> ... <file_n>")
		countFlags.PrintDefaults()
		os.Exit(1)
	}

	return countFlags.Args()
}

func csvBuFlags() []string {
	csvFlags := flag.NewFlagSet("csv", flag.ContinueOnError)
	csvFlags.StringVar(&cargo, "cargo", "", "e.g. Presidente")
	csvFlags.StringVar(&candidatos, "candidatos", "", "Comma-separated list; e.g. 'Branco,Nulo,99'")
	err := csvFlags.Parse(os.Args[3:])
	if err != nil {
		os.Exit(1)
	}

	if len(cargo) == 0 || len(candidatos) == 0 {
		fmt.Println("usage: urna bu csv -cargo <cargo> -candidatos <candidatos> <file_1> ... <file_n>")
		csvFlags.PrintDefaults()
		os.Exit(1)
	}

	return csvFlags.Args()
}

func verifyBuFlags() []string {
	if len(os.Args) > 3 {
		return os.Args[3:]
	}

	fmt.Println("usage: urna bu verify <file_1> ... <file_n>")
	os.Exit(1)

	return []string{}
}

func splitCandidatosIntoSlice() []string {
	candidatos := strings.Split(candidatos, ",")
	for i := range candidatos {
		candidatos[i] = strings.Trim(candidatos[i], " ")
	}

	return candidatos
}
