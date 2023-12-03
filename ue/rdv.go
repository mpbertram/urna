package ue

import (
	"github.com/google/certificate-transparency-go/asn1"
	"os"
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
