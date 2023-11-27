package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	urna "github.com/mpbertram/urna/ue"
)

func Rdv() {
	var function string
	if len(os.Args) > 2 {
		function = os.Args[2]
	}

	switch function {
	case "csv":
		rdvToCsv(verifyRdvFlags())
	default:
		fmt.Println("usage: urna rdv csv <file_1> ... <file_n>")
	}
}

func rdvToCsv(files []string) {
	w := csv.NewWriter(os.Stdout)
	w.Write(
		[]string{
			"ID eleicao",
			"Data geracao",
			"Cargo",
			"Quantidade escolhas",
			"Tipo voto",
			"Voto digitado"})

	for _, f := range files {
		log.Printf("processing file %s", f)

		if strings.HasSuffix(f, ".rdv") {
			rdv, err := urna.ReadRdv(f)

			if err != nil {
				log.Fatal("error reading RDV data ", err)
			}

			processRdv(rdv, w)
		}

		if strings.HasSuffix(f, ".zip") {
			urna.ProcessZip(f, func(rdv urna.EntidadeResultadoRDV) error {
				processRdv(rdv, w)
				return nil
			})
		}
	}
}

func processRdv(rdv urna.EntidadeResultadoRDV, w *csv.Writer) {
	el, err := rdv.Rdv.ReadEleicoes()
	if err != nil {
		log.Fatal("error reading Eleicoes ", err)
	}

	if reflect.TypeOf(el) == reflect.TypeOf([]urna.EleicaoVota{}) {
		for _, e := range el.([]urna.EleicaoVota) {
			for _, vc := range e.VotosCargos {
				processVotos(vc, w, e.IdEleicao, rdv.Cabecalho.DataGeracao)
			}
		}
	}

	if reflect.TypeOf(el) == reflect.TypeOf([]urna.EleicaoSA{}) {
		for _, e := range el.([]urna.EleicaoSA) {
			for _, vc := range e.VotosCargos {
				processVotos(vc, w, e.IdEleicao, rdv.Cabecalho.DataGeracao)
			}
		}
	}
}

func processVotos(vc urna.VotosCargo, w *csv.Writer, id int, date urna.DataHoraJE) {
	var cargo string
	var escolhas int
	var tipoVoto string
	var votoDigitado string

	escolhas = int(vc.QuantidadeEscolhas)

	c, err := vc.ReadIdCargo()
	if err != nil {
		log.Fatal("error reading ID cargo ", err)
	}

	c, ok := c.(urna.CargoConstitucional)
	if ok {
		cargo = c.(urna.CargoConstitucional).String()
	} else {
		cargo = fmt.Sprint(c.(urna.NumeroCargoConsultaLivre))
	}

	for _, v := range vc.Votos {
		votoDigitado = string(v.Digitacao)

		tv, err := urna.TipoVotoRdvFromData(int(v.TipoVoto))
		if err != nil {
			log.Fatal("error reading tipo voto ", err)
		}

		tipoVoto = tv.String()

		w.Write([]string{
			fmt.Sprint(id),
			string(date),
			cargo,
			fmt.Sprint(escolhas),
			tipoVoto,
			votoDigitado})
	}

	w.Flush()
}

func verifyRdvFlags() []string {
	if len(os.Args) > 3 {
		return os.Args[3:]
	}

	fmt.Println("usage: urna rdv csv <file_1> ... <file_n>")
	os.Exit(1)

	return []string{}
}
