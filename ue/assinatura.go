package ue

import (
	"archive/zip"
	"bytes"
	"crypto/sha512"
	"encoding/asn1"
	"fmt"
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

type VerificationResult struct {
	Ok  bool
	Msg string
}

func VerifyAssinaturaVscmr(path string) []VerificationResult {
	a, err := ReadAssinatura(path)
	if err != nil {
		log.Printf("error processing %s", path)
	}

	var results []VerificationResult

	results = append(results, verifyAssinaturaVscmr(path, a.AssinaturaHW)...)
	results = append(results, verifyAssinaturaVscmr(path, a.AssinaturaSW)...)

	return results
}

func verifyAssinaturaVscmr(path string, a EntidadeAssinatura) []VerificationResult {
	var results []VerificationResult

	conteudoAssinado, err := a.ReadConteudoAutoAssinado()
	if err != nil {
		log.Println(err)
	}

	checksum := sha512.Sum512(a.ConteudoAutoAssinado)
	if !slices.Equal(checksum[:], a.AutoAssinado.Assinatura.Hash) {
		results = append(results, VerificationResult{
			Ok: false,
			Msg: fmt.Sprintf(
				"hash check failed for auto content of %s, zona=%s, secao=%s",
				MunicipioByFile(path), ZonaByFile(path), SecaoByFile(path),
			),
		})
	} else {
		results = append(results, VerificationResult{
			Ok: true,
			Msg: fmt.Sprintf(
				"hash check successful for auto content of %s, zona=%s, secao=%s",
				MunicipioByFile(path), ZonaByFile(path), SecaoByFile(path),
			),
		})
	}

	err = a.VerifyAutoSignature()
	if err != nil {
		if err.Error() != "no certificate" {
			results = append(results, VerificationResult{
				Ok: false,
				Msg: fmt.Sprintf(
					"signature check failed for auto content of %s, zona=%s, secao=%s",
					MunicipioByFile(path), ZonaByFile(path), SecaoByFile(path),
				),
			})
		}
	} else {
		results = append(results, VerificationResult{
			Ok: true,
			Msg: fmt.Sprintf(
				"signature check successful for auto content of %s, zona=%s, secao=%s",
				MunicipioByFile(path), ZonaByFile(path), SecaoByFile(path),
			),
		})
	}

	for _, arquivo := range conteudoAssinado.ArquivosAssinados {
		file, err := os.ReadFile(filepath.Join(filepath.Dir(path), arquivo.NomeArquivo))
		if err != nil {
			continue
		}

		checksum := sha512.Sum512(file)
		if !slices.Equal(checksum[:], arquivo.Assinatura.Hash) {
			results = append(results, VerificationResult{
				Ok: false,
				Msg: fmt.Sprintf(
					"hash check failed for file %s, zona=%s, secao=%s",
					MunicipioByFile(arquivo.NomeArquivo), ZonaByFile(arquivo.NomeArquivo), SecaoByFile(arquivo.NomeArquivo),
				),
			})
		} else {
			results = append(results, VerificationResult{
				Ok: true,
				Msg: fmt.Sprintf(
					"hash check successful for file %s, zona=%s, secao=%s",
					MunicipioByFile(arquivo.NomeArquivo), ZonaByFile(arquivo.NomeArquivo), SecaoByFile(arquivo.NomeArquivo),
				),
			})
		}

		err = a.VerifySignature(arquivo)
		if err != nil {
			if err.Error() != "no certificate" {
				results = append(results, VerificationResult{
					Ok: false,
					Msg: fmt.Sprintf(
						"signature check failed for %s, zona=%s, secao=%s",
						MunicipioByFile(arquivo.NomeArquivo), ZonaByFile(arquivo.NomeArquivo), SecaoByFile(arquivo.NomeArquivo),
					),
				})
			}
		} else {
			results = append(results, VerificationResult{
				Ok: true,
				Msg: fmt.Sprintf(
					"signature check successful for %s, zona=%s, secao=%s",
					MunicipioByFile(arquivo.NomeArquivo), ZonaByFile(arquivo.NomeArquivo), SecaoByFile(arquivo.NomeArquivo),
				),
			})
		}
	}

	return results
}

func VerifyAssinaturaZip(path string) []VerificationResult {
	var results []VerificationResult

	ProcessZip(path, func(e EntidadeAssinaturaResultado, ctx ZipProcessCtx) {
		results = append(results, verifyAssinaturaZip(ctx, e.AssinaturaHW)...)
		results = append(results, verifyAssinaturaZip(ctx, e.AssinaturaSW)...)
	})

	return results
}

func verifyAssinaturaZip(ctx ZipProcessCtx, assinatura EntidadeAssinatura) []VerificationResult {
	var results []VerificationResult

	arquivosAssinados, err := assinatura.ReadConteudoAutoAssinado()
	if err != nil {
		log.Fatal(err)
	}

	checksum := sha512.Sum512(assinatura.ConteudoAutoAssinado)
	if !slices.Equal(checksum[:], assinatura.AutoAssinado.Assinatura.Hash) {
		results = append(results, VerificationResult{
			Ok: false,
			Msg: fmt.Sprintf(
				"hash check failed for auto content of %s, zona=%s, secao=%s",
				MunicipioByFile(ctx.Filename), ZonaByFile(ctx.Filename), SecaoByFile(ctx.Filename),
			),
		})
	} else {
		results = append(results, VerificationResult{
			Ok: true,
			Msg: fmt.Sprintf(
				"hash check successful for auto content of %s, zona=%s, secao=%s",
				MunicipioByFile(ctx.Filename), ZonaByFile(ctx.Filename), SecaoByFile(ctx.Filename),
			),
		})
	}

	err = assinatura.VerifyAutoSignature()
	if err != nil {
		if err.Error() != "no certificate" {
			results = append(results, VerificationResult{
				Ok: false,
				Msg: fmt.Sprintf(
					"signature check failed for auto content of %s, zona=%s, secao=%s",
					MunicipioByFile(ctx.Filename), ZonaByFile(ctx.Filename), SecaoByFile(ctx.Filename),
				),
			})
		}
	} else {
		results = append(results, VerificationResult{
			Ok: true,
			Msg: fmt.Sprintf(
				"signature check successful for auto content of %s, zona=%s, secao=%s",
				MunicipioByFile(ctx.Filename), ZonaByFile(ctx.Filename), SecaoByFile(ctx.Filename),
			),
		})
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
					results = append(results, VerificationResult{
						Ok: false,
						Msg: fmt.Sprintf(
							"hash check failed for %s, zona=%s, secao=%s",
							MunicipioByFile(arquivo.NomeArquivo), ZonaByFile(arquivo.NomeArquivo), SecaoByFile(arquivo.NomeArquivo),
						),
					})
				} else {
					results = append(results, VerificationResult{
						Ok: true,
						Msg: fmt.Sprintf(
							"hash check successful for %s, zona=%s, secao=%s",
							MunicipioByFile(arquivo.NomeArquivo), ZonaByFile(arquivo.NomeArquivo), SecaoByFile(arquivo.NomeArquivo),
						),
					})
				}

				err = assinatura.VerifySignature(arquivo)
				if err != nil {
					if err.Error() != "no certificate" {
						results = append(results, VerificationResult{
							Ok: false,
							Msg: fmt.Sprintf(
								"signature check failed for %s, zona=%s, secao=%s",
								MunicipioByFile(arquivo.NomeArquivo), ZonaByFile(arquivo.NomeArquivo), SecaoByFile(arquivo.NomeArquivo),
							),
						})
					}
				} else {
					results = append(results, VerificationResult{
						Ok: true,
						Msg: fmt.Sprintf(
							"signature check successful for %s, zona=%s, secao=%s",
							MunicipioByFile(arquivo.NomeArquivo), ZonaByFile(arquivo.NomeArquivo), SecaoByFile(arquivo.NomeArquivo),
						),
					})
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

	return results
}
