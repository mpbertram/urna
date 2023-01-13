package urna

import (
	"encoding/asn1"
	"os"
	"testing"
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

func TestAssinatura(t *testing.T) {
	_, err := ReadAssinatura("test-data/urna.vscmr")
	if err != nil {
		t.Error(err)
	}
}
