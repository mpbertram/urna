package ue

import (
	"archive/zip"
	"bytes"
	"crypto/sha512"
	"encoding/asn1"
	"io"
	"log"
	"os"
	"path/filepath"

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

func VerifyAssinaturaVscmr(path string) {
	a, err := ReadAssinatura(path)
	if err != nil {
		log.Printf("error processing %s", path)
	}

	verifyAssinaturaVscmr(path, a.AssinaturaHW)
	verifyAssinaturaVscmr(path, a.AssinaturaSW)
}

func verifyAssinaturaVscmr(path string, a EntidadeAssinatura) {
	conteudoAssinado, err := a.ReadConteudoAutoAssinado()
	if err != nil {
		log.Println(err)
	}

	checksum := sha512.Sum512(a.ConteudoAutoAssinado)
	if !slices.Equal(checksum[:], a.AutoAssinado.Assinatura.Hash) {
		log.Printf("hash check failed for auto content of %s", path)
	} else {
		log.Printf("hash check successful for auto content of %s", path)
	}

	err = a.VerifyAutoSignature()
	if err != nil {
		if err.Error() != "no certificate" {
			log.Printf("signature check failed for auto content of %s", path)
		}
	} else {
		log.Printf("signature check successful for auto content of %s", path)
	}

	for _, arquivo := range conteudoAssinado.ArquivosAssinados {
		file, err := os.ReadFile(filepath.Join(filepath.Dir(path), arquivo.NomeArquivo))
		if err != nil {
			continue
		}

		checksum := sha512.Sum512(file)
		if !slices.Equal(checksum[:], arquivo.Assinatura.Hash) {
			log.Printf("hash check failed for file %s", arquivo.NomeArquivo)
		} else {
			log.Printf("hash check successful for %s", arquivo.NomeArquivo)
		}

		err = a.VerifySignature(arquivo)
		if err != nil {
			if err.Error() != "no certificate" {
				log.Printf("signature check failed for %s", arquivo.NomeArquivo)
			}
		} else {
			log.Printf("signature check successful for %s", arquivo.NomeArquivo)
		}
	}
}

func VerifyAssinaturaZip(path string) {
	ProcessZip(path, func(e EntidadeAssinaturaResultado, ctx ZipProcessCtx) {
		verifyAssinaturaZip(ctx, e.AssinaturaHW)
		verifyAssinaturaZip(ctx, e.AssinaturaSW)
	})
}

func verifyAssinaturaZip(ctx ZipProcessCtx, assinatura EntidadeAssinatura) {
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
		if err.Error() != "no certificate" {
			log.Printf("signature check failed for auto content of %s", ctx.Filename)
		}
	} else {
		log.Printf("signature check successful for auto content of %s", ctx.Filename)
	}

	var count int
	ProcessZipRaw(ctx.ZipFilename, func(f *zip.File) bool {
		for _, arquivo := range arquivosAssinados.ArquivosAssinados {
			if f.Name == arquivo.NomeArquivo {
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
					if err.Error() != "no certificate" {
						log.Printf("signature check failed for %s", arquivo.NomeArquivo)
					}
				} else {
					log.Printf("signature check successful for %s", arquivo.NomeArquivo)
				}

				count++
				rc.Close()
			}
		}

		if count == len(arquivosAssinados.ArquivosAssinados) {
			count = 0
			return true
		}

		return false
	})
}
