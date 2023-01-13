package urna

import (
	"crypto/x509"
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
	a, err := ReadAssinatura("test-data/urna.vscmr")
	if err != nil {
		t.Error(err)
	}
	_, err = x509.ParseCertificate(a.AssinaturaHW.CertificadoDigital)
	if err != nil {
		t.Error(err)
	}
}
