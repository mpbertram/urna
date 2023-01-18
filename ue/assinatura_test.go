package ue

import (
	"crypto/x509"
	"log"
	"testing"
)

func TestAssinaturaZip(t *testing.T) {
	VerifyAssinaturas("test-data/signature")
}

func TestAssinatura(t *testing.T) {
	vscmr, err := ReadAssinatura("test-data/urna.vscmr")
	if err != nil {
		t.Error(err)
	}

	as, err := vscmr.AssinaturaHW.ReadConteudoAutoAssinado()
	if err != nil {
		t.Error(err)
	}
	for _, a := range as.ArquivosAssinados {
		log.Println(a.NomeArquivo)
	}

	_, err = x509.ParseCertificate(vscmr.AssinaturaHW.CertificadoDigital)
	if err != nil {
		t.Error(err)
	}
}
