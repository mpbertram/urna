package ue

import (
	"archive/zip"
	"bytes"
	"crypto/sha512"
	"encoding/asn1"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/exp/slices"
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

func VerifyAssinaturas(dir string) {
	ProcessAllZip(dir, func(e EntidadeAssinaturaResultado) {
		verifyEntidadeAssinatura(dir, e.AssinaturaHW)
		verifyEntidadeAssinatura(dir, e.AssinaturaSW)
	})
}

func verifyEntidadeAssinatura(dir string, a EntidadeAssinatura) {
	as, err := a.ReadConteudoAutoAssinado()
	if err != nil {
		log.Fatal(err)
	}

	for _, a := range as.ArquivosAssinados {
		ProcessAllZipRaw(dir, func(f *zip.File) {
			if strings.EqualFold(f.Name, a.NomeArquivo) {
				rc, err := f.Open()
				if err != nil {
					log.Fatal(err)
				}

				var buf bytes.Buffer
				io.Copy(&buf, rc)

				checksum := sha512.Sum512(buf.Bytes())
				if !slices.Equal(checksum[:], a.Assinatura.Hash) {
					log.Fatalf("hash check failed for file %s", a.NomeArquivo)
				} else {
					log.Printf("hash check successful for %s", a.NomeArquivo)
				}

				rc.Close()
			}
		})
	}
}
