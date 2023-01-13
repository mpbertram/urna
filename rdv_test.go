package urna

import (
	"encoding/asn1"
	"os"
	"testing"
)

func ReadRdv(file string) (EntidadeResultadoRDV, error) {
	f, err := os.ReadFile(file)
	if err != nil {
		return EntidadeResultadoRDV{}, err
	}

	var rdv EntidadeResultadoRDV
	_, err = asn1.Unmarshal(f, &rdv)
	if err != nil {
		return EntidadeResultadoRDV{}, err
	}

	return rdv, nil
}

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
