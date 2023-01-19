package ue

import (
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

	_, err = vscmr.AssinaturaHW.ParseCertificate()
	if err != nil {
		t.Error(err)
	}
}
