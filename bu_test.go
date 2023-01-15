package urna

import (
	"reflect"
	"testing"
)

func TestZips(t *testing.T) {
	all := make(map[string]map[string]int)
	ProcessAllZips("test-data", func(be []BuEntry) error {
		votes := ComputeVotos(be, []CargoConstitucional{Presidente})
		for k1, v1 := range votes {
			if all[k1] == nil {
				all[k1] = v1
			} else {
				for k2, v2 := range v1 {
					all[k1][k2] = all[k1][k2] + v2
				}
			}
		}
		return nil
	})
	if all[Presidente.String()][Nulo.String()] != 6 {
		t.Errorf("Wrong count for Nulo (%d)", all[Presidente.String()][Nulo.String()])
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

	err = ValidateVotos(bus)
	if err != nil {
		t.Error(err)
	}

	v := ComputeVotos(bus, []CargoConstitucional{Presidente})
	if v[Presidente.String()][Nulo.String()] != 6 {
		t.Errorf("wrong count for Nulo (%d)", v[Presidente.String()][Nulo.String()])
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
