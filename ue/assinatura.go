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
	ProcessAllZip(dir, func(e EntidadeAssinaturaResultado, ctx ZipProcessCtx) {
		verifyEntidadeAssinatura(ctx, e.AssinaturaHW)
		verifyEntidadeAssinatura(ctx, e.AssinaturaSW)
	})
}

func verifyEntidadeAssinatura(ctx ZipProcessCtx, a EntidadeAssinatura) {
	as, err := a.ReadConteudoAutoAssinado()
	if err != nil {
		log.Fatal(err)
	}

	checksum := sha512.Sum512(a.ConteudoAutoAssinado)
	if !slices.Equal(checksum[:], a.AutoAssinado.Assinatura.Hash) {
		log.Fatalf("hash check failed for auto content of %s", ctx.Filename)
	} else {
		log.Printf("hash check successful for auto content of %s", ctx.Filename)
	}

	ProcessZipRaw(ctx.ZipFilename, func(f *zip.File) {
		for _, a := range as.ArquivosAssinados {
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
		}
	})
}
