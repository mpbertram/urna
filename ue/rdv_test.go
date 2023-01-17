package ue

import (
	"testing"
)

func TestRdv(t *testing.T) {
	rdv, err := ReadRdv("test-data/urna.rdv")
	if err != nil {
		t.Error(err)
	}

	e, err := rdv.Rdv.ReadEleicoes()
	if err != nil {
		t.Error(err)
	}

	for _, eleicao := range e.([]EleicaoVota) {
		for _, votosCargo := range eleicao.VotosCargos {
			_, err := votosCargo.ReadIdCargo()
			if err != nil {
				t.Error(err)
			}
		}
	}
}
