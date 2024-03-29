package ue

import (
	"testing"
)

func TestAssinaturaZip(t *testing.T) {
	VerifyAssinaturaZip("test-data/o00407-0100700090001.zip")
}

func TestAssinatura(t *testing.T) {
	vscmr, err := ReadAssinatura("test-data/urna.vscmr")
	if err != nil {
		t.Error(err)
	}

	_, err = vscmr.AssinaturaHW.ReadConteudoAssinado()
	if err != nil {
		t.Error(err)
	}

	_, err = vscmr.AssinaturaHW.ParseCertificate()
	if err != nil {
		t.Error(err)
	}
}
