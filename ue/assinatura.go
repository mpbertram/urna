package ue

import (
	"archive/zip"
	"bytes"
	"crypto/sha512"
	"fmt"
	"github.com/google/certificate-transparency-go/asn1"
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

type VerificationResultStatus bool

const (
	Ok  VerificationResultStatus = true
	Nok VerificationResultStatus = false
)

func (s VerificationResultStatus) String() string {
	switch s {
	case Ok:
		return "ok"
	default:
		return "nok"
	}
}

type VerificationResultType uint8

const (
	Hash        VerificationResultType = 0
	Signature   VerificationResultType = 1
	Payload     VerificationResultType = 2
	Certificate VerificationResultType = 3
)

func (t VerificationResultType) String() string {
	switch t {
	case Hash:
		return "hash"
	case Signature:
		return "signature"
	case Payload:
		return "payload"
	case Certificate:
		return "cert"
	default:
		return ""
	}
}

type VerificationResult struct {
	Type      VerificationResultType
	Ok        VerificationResultStatus
	Err       error
	Filename  string
	Municipio string
	Zona      string
	Secao     string
	Payload   []byte
}

func (r VerificationResult) Msg() string {
	if r.Payload != nil {
		return fmt.Sprintf(
			"[%s] [%s] municipio=%s, zona=%s, secao=%s, payload=%s",
			r.Ok.String(),
			r.Type.String(),
			r.Municipio,
			r.Zona,
			r.Secao,
			r.Payload,
		)
	}

	return fmt.Sprintf(
		"[%s] [%s] file=%s municipio=%s, zona=%s, secao=%s",
		r.Ok.String(),
		r.Type.String(),
		r.Filename,
		r.Municipio,
		r.Zona,
		r.Secao,
	)
}

func VerifyCertsVscmr(path string) []VerificationResult {
	sig, err := ReadAssinatura(path)
	if err != nil {
		log.Printf("error processing %s", path)
	}

	var errors []VerificationResult

	errors = append(errors, parseCertificate(sig.AssinaturaHW, path)...)
	errors = append(errors, parseCertificate(sig.AssinaturaSW, path)...)

	return errors
}

func VerifyCertsZip(path string) []VerificationResult {
	var errors []VerificationResult

	ProcessZip(path, func(sig EntidadeAssinaturaResultado, ctx ZipProcessCtx) {
		errors = append(errors, parseCertificate(sig.AssinaturaHW, ctx.Filename)...)
		errors = append(errors, parseCertificate(sig.AssinaturaSW, ctx.Filename)...)
	})

	return errors
}

func parseCertificate(sig EntidadeAssinatura, path string) []VerificationResult {
	var errors []VerificationResult

	if len(sig.CertificadoDigital) > 0 {
		_, err := sig.ParseCertificateNative()
		if err != nil {
			errors = append(errors, newCertError(err, path))
		} else {
			errors = append(errors, newCertOk(path))
		}
	}

	return errors
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

func verifyAssinaturaVscmr(path string, sig EntidadeAssinatura) []VerificationResult {
	var results []VerificationResult

	conteudoAssinado, err := sig.ReadConteudoAssinado()
	if err != nil {
		log.Println(err)
	}

	results = append(results, verifyAutoContent(sig, path)...)

	for _, arquivo := range conteudoAssinado.ArquivosAssinados {
		file, err := os.ReadFile(filepath.Join(filepath.Dir(path), arquivo.NomeArquivo))
		if err != nil {
			continue
		}

		results = append(results, verifyHash(file, arquivo.Assinatura.Hash, arquivo.NomeArquivo))
		results = append(results, verifySignature(sig, arquivo))
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

func verifyAssinaturaZip(ctx ZipProcessCtx, sig EntidadeAssinatura) []VerificationResult {
	var results []VerificationResult

	conteudoAssinado, err := sig.ReadConteudoAssinado()
	if err != nil {
		log.Fatal(err)
	}

	results = append(results, verifyAutoContent(sig, ctx.Filename)...)

	var count int
	ProcessZipRaw(ctx.ZipFilename, func(f *zip.File) bool {
		for _, arquivo := range conteudoAssinado.ArquivosAssinados {
			if f.Name == arquivo.NomeArquivo {
				rc, err := f.Open()
				if err != nil {
					log.Fatal(err)
				}

				var buf bytes.Buffer
				io.Copy(&buf, rc)

				results = append(results, verifyHash(buf.Bytes(), arquivo.Assinatura.Hash, arquivo.NomeArquivo))
				results = append(results, verifySignature(sig, arquivo))

				count++
				rc.Close()
			}
		}

		if count == len(conteudoAssinado.ArquivosAssinados) {
			count = 0
			return true
		}

		return false
	})

	return results
}

func verifyAutoContent(sig EntidadeAssinatura, filename string) []VerificationResult {
	var results []VerificationResult
	results = append(results, verifyHash(sig.ConteudoAutoAssinado, sig.AutoAssinado.Assinatura.Hash, filename))
	results = append(results, verifyAutoSignature(sig, filename))
	return results
}

func verifyHash(content []byte, hash []byte, filename string) VerificationResult {
	checksum := sha512.Sum512(content)
	if !slices.Equal(checksum[:], hash) {
		return newHashError(filename)
	}

	return newHashOk(filename)
}

func verifyAutoSignature(assinatura EntidadeAssinatura, filename string) VerificationResult {
	err := assinatura.VerifyAutoSignature()
	if err != nil {
		if err.Error() != "no certificate" {
			return newSigError(err, filename)
		}
	}

	return newSigOk(filename)
}

func verifySignature(assinatura EntidadeAssinatura, arquivo AssinaturaArquivo) VerificationResult {
	err := assinatura.VerifySignature(arquivo)
	if err != nil {
		if err.Error() != "no certificate" {
			return newSigError(err, arquivo.NomeArquivo)
		}
	}

	return newSigOk(arquivo.NomeArquivo)
}

func newSigOk(filename string) VerificationResult {
	return VerificationResult{
		Type:      Signature,
		Ok:        Ok,
		Filename:  filename,
		Municipio: MunicipioByFile(filename),
		Zona:      ZonaByFile(filename),
		Secao:     SecaoByFile(filename),
	}
}

func newSigError(err error, filename string) VerificationResult {
	return VerificationResult{
		Type:      Signature,
		Ok:        Nok,
		Err:       err,
		Filename:  filename,
		Municipio: MunicipioByFile(filename),
		Zona:      ZonaByFile(filename),
		Secao:     SecaoByFile(filename),
	}
}

func newHashOk(filename string) VerificationResult {
	return VerificationResult{
		Type:      Hash,
		Ok:        Ok,
		Filename:  filename,
		Municipio: MunicipioByFile(filename),
		Zona:      ZonaByFile(filename),
		Secao:     SecaoByFile(filename),
	}
}

func newHashError(filename string) VerificationResult {
	return VerificationResult{
		Type:      Hash,
		Ok:        Nok,
		Filename:  filename,
		Municipio: MunicipioByFile(filename),
		Zona:      ZonaByFile(filename),
		Secao:     SecaoByFile(filename),
	}
}

func newCertOk(filename string) VerificationResult {
	return VerificationResult{
		Type:      Certificate,
		Ok:        Ok,
		Filename:  filename,
		Municipio: MunicipioByFile(filename),
		Zona:      ZonaByFile(filename),
		Secao:     SecaoByFile(filename),
	}
}

func newCertError(err error, filename string) VerificationResult {
	return VerificationResult{
		Type:      Certificate,
		Ok:        Nok,
		Err:       err,
		Filename:  filename,
		Municipio: MunicipioByFile(filename),
		Zona:      ZonaByFile(filename),
		Secao:     SecaoByFile(filename),
	}
}
