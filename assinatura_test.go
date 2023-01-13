package urna

import (
	"crypto/x509"
	"testing"
)

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
