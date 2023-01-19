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

func verifyEntidadeAssinatura(ctx ZipProcessCtx, assinatura EntidadeAssinatura) {
	arquivosAssinados, err := assinatura.ReadConteudoAutoAssinado()
	if err != nil {
		log.Fatal(err)
	}

	checksum := sha512.Sum512(assinatura.ConteudoAutoAssinado)
	if !slices.Equal(checksum[:], assinatura.AutoAssinado.Assinatura.Hash) {
		log.Printf("hash check failed for auto content of %s", ctx.Filename)
	} else {
		log.Printf("hash check successful for auto content of %s", ctx.Filename)
	}

	err = assinatura.VerifyAutoSignature()
	if err != nil {
		if !strings.EqualFold(err.Error(), "no certificate") {
			log.Printf("signature check failed for auto content of %s", ctx.Filename)
		}
	} else {
		log.Printf("signature check successful for auto content of %s", ctx.Filename)
	}

	ProcessZipRaw(ctx.ZipFilename, func(f *zip.File) {
		for _, arquivo := range arquivosAssinados.ArquivosAssinados {
			if strings.EqualFold(f.Name, arquivo.NomeArquivo) {
				rc, err := f.Open()
				if err != nil {
					log.Fatal(err)
				}

				var buf bytes.Buffer
				io.Copy(&buf, rc)

				checksum := sha512.Sum512(buf.Bytes())
				if !slices.Equal(checksum[:], arquivo.Assinatura.Hash) {
					log.Printf("hash check failed for file %s", arquivo.NomeArquivo)
				} else {
					log.Printf("hash check successful for %s", arquivo.NomeArquivo)
				}

				err = assinatura.VerifySignature(arquivo)
				if err != nil {
					if !strings.EqualFold(err.Error(), "no certificate") {
						log.Printf("signature check failed for %s", arquivo.NomeArquivo)
					}
				} else {
					log.Printf("signature check successful for %s", arquivo.NomeArquivo)
				}

				rc.Close()
			}
		}
	})
}
