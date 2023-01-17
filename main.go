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

var function = flag.String("function", "", "<compute|validate|all>")
var cargo = flag.String("cargo", "", "<Presidente|Governador>")
var candidatos = flag.String("candidatos", "", "e.g. 'Branco,Nulo,99'")
var folder = flag.String("folder", "", "<folder>")

func main() {
	flag.Parse()

	switch *function {
	case "compute":
		compute()
	case "validate":
		validate()
	case "all":
		all()
	}
}

func all() {
	cargo := urna.CargoConstitucionalFromString(*cargo)
	cargos := []urna.CargoConstitucional{cargo}
	candidatos := strings.Split(*candidatos, ",")
	for i := range candidatos {
		candidatos[i] = strings.Trim(candidatos[i], " ")
	}

	w := csv.NewWriter(os.Stdout)
	w.Write(append([]string{"UF", "Municipio", "Local", "Secao"}, candidatos...))

	bus, err := urna.ReadAllBu(*folder)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range bus {
		bu, err := entry.ReadBu()
		if err != nil {
			log.Fatal(err)
		}

		votos := urna.ComputeVotosBu(bu, cargos)
		mun, _ := urna.MunicipioFromId(int(bu.IdentificacaoSecao.MunicipioZona.Municipio))
		var votosForCandidato []string
		for _, candidato := range candidatos {
			votosForCandidato = append(votosForCandidato, fmt.Sprint(votos[cargo.String()][candidato]))
		}
		w.Write(
			append(
				[]string{
					mun.Uf,
					mun.Nome,
					fmt.Sprint(bu.IdentificacaoSecao.Local),
					fmt.Sprint(bu.IdentificacaoSecao.Secao)},
				votosForCandidato...))
	}
	w.Flush()

	urna.ProcessAllZip(*folder, func(bu urna.EntidadeBoletimUrna) error {
		votos := urna.ComputeVotosBu(bu, cargos)
		mun, _ := urna.MunicipioFromId(int(bu.IdentificacaoSecao.MunicipioZona.Municipio))
		var votosForCandidato []string
		for _, candidato := range candidatos {
			votosForCandidato = append(votosForCandidato, fmt.Sprint(votos[cargo.String()][candidato]))
		}
		w.Write(
			append(
				[]string{
					mun.Uf,
					mun.Nome,
					fmt.Sprint(bu.IdentificacaoSecao.Local),
					fmt.Sprint(bu.IdentificacaoSecao.Secao)},
				votosForCandidato...))

		w.Flush()
		return nil
	})
}

func compute() {
	bus, err := urna.ReadAllBu(*folder)
	if err != nil {
		log.Fatal(err)
	}

	cargos := []urna.CargoConstitucional{urna.CargoConstitucionalFromString(*cargo)}
	votos := urna.ComputeVotos(bus, cargos)

	urna.ProcessAllZip(*folder, func(ebu urna.EntidadeBoletimUrna) error {
		for cargo, candidato := range urna.ComputeVotosBu(ebu, cargos) {
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

func validate() {
	bus, err := urna.ReadAllBu(*folder)
	if err != nil {
		log.Fatal(err)
	}

	err = urna.ValidateVotos(bus)
	if err != nil {
		log.Fatal(err)
	}

	urna.ProcessAllZip(*folder, func(ebu urna.EntidadeBoletimUrna) error {
		err = urna.ValidateVotosBu(ebu)
		if err != nil {
			log.Fatal(err)
		}

		return nil
	})

	log.Println("validation successful")
}
