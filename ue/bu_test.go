package ue

import (
	"reflect"
	"testing"
)

func TestZip(t *testing.T) {
	all := make(map[CargoConstitucional]map[string]int)
	ProcessAllZip("test-data", func(eeg EntidadeEnvelopeGenerico) error {
		bu, err := eeg.ReadBu()
		if err != nil {
			t.Error(err)
		}

		results := ValidateVotosBu(bu)
		for _, r := range results {
			if !r.Ok {
				t.Error(r.Msg())
			}
		}

		votos := CountVotosBu(bu, []CargoConstitucional{Presidente})
		for cargo, candidato := range votos {
			if all[cargo] == nil {
				all[cargo] = candidato
			} else {
				for candidato, numVotos := range candidato {
					all[cargo][candidato] = all[cargo][candidato] + numVotos
				}
			}
		}
		return nil
	})
	if all[Presidente][Nulo.String()] != 6 {
		t.Errorf("Wrong count for Nulo (%d)", all[Presidente][Nulo.String()])
	}
}

func TestBu(t *testing.T) {
	bus, err := ReadAllBu("test-data")
	if err != nil {
		t.Error("could not read BUs", err)
	}

	bu, err := bus[0].ReadBu()
	if err != nil {
		t.Error("could not read BU", err)
	}

	result := ValidateVotosBu(bu)
	for _, r := range result {
		if !r.Ok {
			t.Error(r.Msg())
		}
	}

	v := CountVotos(bus, []CargoConstitucional{Presidente})
	if v[Presidente][Nulo.String()] != 6 {
		t.Errorf("wrong count for Nulo (%d)", v[Presidente][Nulo.String()])
	}

	d, err := bu.ReadDadosSecaoSA()
	if err != nil {
		t.Error("could not read DadosSecao", bu)
	}

	if reflect.TypeOf(d) != reflect.TypeOf(DadosSecao{}) {
		t.Error("not dados secao", d)
	}

	if d.(DadosSecao).DataHoraAbertura != "20221002T080001" {
		t.Error("wrong DataHoraAbertura", d.(DadosSecao).DataHoraAbertura)
	}
	if d.(DadosSecao).DataHoraEncerramento != "20221002T170204" {
		t.Error("wrong DataHoraEncerramento", d.(DadosSecao).DataHoraEncerramento)
	}

	i, err := bu.Urna.CorrespondenciaResultado.ReadIdentificacao()
	if err != nil {
		t.Error("could not read Identificacao", bu.Urna.CorrespondenciaResultado)
	}

	if reflect.TypeOf(i) != reflect.TypeOf(IdentificacaoSecaoEleitoral{}) {
		t.Error("not identificacao secao eleitoral")
	}

	if i.(IdentificacaoSecaoEleitoral).MunicipioZona.Municipio != 88986 {
		t.Error("wrong municipio", i.(IdentificacaoSecaoEleitoral).MunicipioZona.Municipio)
	}
	if i.(IdentificacaoSecaoEleitoral).MunicipioZona.Zona != 7 {
		t.Error("wrong zona", i.(IdentificacaoSecaoEleitoral).MunicipioZona.Zona)
	}
	if i.(IdentificacaoSecaoEleitoral).Local != 1 {
		t.Error("wrong local", i.(IdentificacaoSecaoEleitoral).Local)
	}
	if i.(IdentificacaoSecaoEleitoral).Secao != 55 {
		t.Error("wrong secao", i.(IdentificacaoSecaoEleitoral).Secao)
	}
}
