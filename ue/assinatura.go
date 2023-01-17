package ue

import (
	"encoding/asn1"
	"os"
)

func ReadAssinatura(file string) (EntidadeAssinaturaResultado, error) {
	f, err := os.ReadFile(file)
	if err != nil {
		return EntidadeAssinaturaResultado{}, err
	}

	var a EntidadeAssinaturaResultado
	_, err = asn1.Unmarshal(f, &a)
	if err != nil {
		return EntidadeAssinaturaResultado{}, err
	}

	return a, nil
}
