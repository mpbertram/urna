package urna

import (
	"reflect"
	"testing"
)

func TestBu(t *testing.T) {
	bus, err := ReadAllBu("test-data")
	bu := bus[0]

	if err != nil {
		t.Error("could not read BU")
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
