// Cf. https://www.tse.jus.br/eleicoes/eleicoes-2022/documentacao-tecnica-do-software-da-urna-eletronica

package ue

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"encoding/pem"
	"errors"
	"math/big"
	"strings"

	"github.com/google/certificate-transparency-go/asn1"
	"github.com/google/certificate-transparency-go/x509"
)

// ENUMS
// Tipos de algoritmos de assinatura (cepesc é o algoritmo padrão (ainda não há previsão de uso dos demais)).

type AlgoritmoAssinatura int

const (
	Rsa                        AlgoritmoAssinatura = 1
	Ecdsa                      AlgoritmoAssinatura = 2
	Cepesc                     AlgoritmoAssinatura = 3
	UnknownAlgoritmoAssinatura AlgoritmoAssinatura = 0xfffffffffffffff
)

func AlgoritmoAssinaturaFromData(data int) (AlgoritmoAssinatura, error) {

	switch v := data; v {
	case 0x01:
		return Rsa, nil
	case 0x02:
		return Ecdsa, nil
	case 0x03:
		return Cepesc, nil
	}

	return UnknownAlgoritmoAssinatura, errors.New("invalid data")
}

// Tipos de algoritmos de hash (Todos os algoritmos devem ser suportados mas sha512 é o padrão).
type AlgoritmoHash int

const (
	Sha1                 AlgoritmoHash = 1
	Sha256               AlgoritmoHash = 2
	Sha384               AlgoritmoHash = 3
	Sha512               AlgoritmoHash = 4
	UnknownAlgoritmoHash AlgoritmoHash = 0xfffffffffffffff
)

func AlgoritmoHashFromData(data int) (AlgoritmoHash, error) {

	switch v := data; v {
	case 0x01:
		return Sha1, nil
	case 0x02:
		return Sha256, nil
	case 0x03:
		return Sha384, nil
	case 0x04:
		return Sha512, nil
	}

	return UnknownAlgoritmoHash, errors.New("invalid data")
}

// Tipos de modelos de urna eletrônica.
type ModeloUrna byte

const (
	Ue2009 ModeloUrna = 9  // Urna modelo 2009.
	Ue2010 ModeloUrna = 10 // Urna modelo 2010.
	Ue2011 ModeloUrna = 11 // Urna modelo 2011.
	Ue2013 ModeloUrna = 13 // Urna modelo 2013.
	Ue2015 ModeloUrna = 15 // Urna modelo 2015.
	Ue2020 ModeloUrna = 20 // Urna modelo 2020.
)

// ENVELOPE
// Entidade que engloba a lista de assinaturas utilizadas para assinar os arquivos para manter a integridade e segurança dos dados.
type EntidadeAssinatura struct {
	DataHoraCriacao      DataHoraJE            // Data e Hora da criacao do arquivo.
	Versao               int                   // Versao do protocolo (Alterações devem gerar novo valor. Nas eleições de 2012 foi utilizado o enumerado de valor 1 a partir de 2014 utilizar o valor 2).
	AutoAssinado         AutoAssinaturaDigital // Informações da auto assinatura digital.
	ConteudoAutoAssinado []byte                // Conteúdo da assinatura do próprio arquivo.
	CertificadoDigital   []byte                `asn1:"optional"` // Certificado digital da urna eletrônica.
	ConjuntoChave        string                `asn1:"optional"` // Identificador do conjunto de chaves usado para assinar o pacote.
}

func (alg AlgoritmoHash) GetHashFunction() (crypto.Hash, error) {
	switch alg {
	case Sha1:
		return crypto.SHA1, nil
	case Sha256:
		return crypto.SHA256, nil
	case Sha384:
		return crypto.SHA384, nil
	case Sha512:
		return crypto.SHA512, nil
	default:
		return 0, errors.New("invalid hash function")
	}
}

func (sig EntidadeAssinatura) ReadConteudoAutoAssinado() (Assinatura, error) {
	var a Assinatura
	_, err := asn1.Unmarshal(sig.ConteudoAutoAssinado, &a)
	if err != nil {
		return Assinatura{}, err
	}
	return a, nil
}

func (sig EntidadeAssinatura) ParseCertificate() (*x509.Certificate, error) {

	cert, err := x509.ParseCertificate(sig.CertificadoDigital)

	if err != nil {
		if err.Error() == "asn1: syntax error: trailing data" {
			digSig := bytes.Clone(sig.CertificadoDigital)

			for {
				if bytes.HasSuffix(digSig, []byte{0x00}) {
					digSig = digSig[:len(digSig)-1]
					cert, err = x509.ParseCertificate(digSig)
					if err == nil {
						break
					}
				} else {
					break
				}
			}

			if err != nil {
				return &x509.Certificate{}, err
			}

			return cert, nil
		}

		if strings.Contains(err.Error(), "tags don't match") {
			pemCert, _ := pem.Decode(sig.CertificadoDigital)

			if cert, err = x509.ParseCertificate(pemCert.Bytes); err != nil {
				return &x509.Certificate{}, err
			} else {
				return cert, nil
			}
		}

		return &x509.Certificate{}, err
	}

	return cert, nil
}

type dsaSignature struct {
	R, S *big.Int
}

func (sig EntidadeAssinatura) VerifyAutoSignature() error {
	return sig.verifySignature(sig.AutoAssinado.Assinatura)
}

func (sig EntidadeAssinatura) VerifySignature(file AssinaturaArquivo) error {
	return sig.verifySignature(file.Assinatura)
}

func (sig EntidadeAssinatura) verifySignature(digSig AssinaturaDigital) error {
	if len(sig.CertificadoDigital) == 0 {
		return errors.New("no certificate")
	}

	cert, err := sig.ParseCertificate()
	if err != nil {
		return err
	}

	switch publicKey := cert.PublicKey.(type) {
	case *ecdsa.PublicKey:
		algHash, err := AlgoritmoHashFromData(int(sig.AutoAssinado.AlgoritmoHash.Algoritmo))
		if err != nil {
			return err
		}

		signed := digSig.Hash

		hashType, err := algHash.GetHashFunction()
		if err != nil {
			return err
		}

		hash := hashType.New()
		hash.Write(signed)
		signed = hash.Sum(nil)

		ecdsaSig := new(dsaSignature)
		if rest, err := asn1.UnmarshalWithParams(digSig.Assinatura, ecdsaSig, "lax"); err != nil {
			return err
		} else if len(rest) != 0 {
			return errors.New("x509: trailing data after ECDSA signature")
		}
		if ecdsaSig.R.Sign() <= 0 || ecdsaSig.S.Sign() <= 0 {
			return errors.New("x509: ECDSA signature contained zero or negative values")
		}
		if !ecdsa.Verify(publicKey, signed, ecdsaSig.R, ecdsaSig.S) {
			return errors.New("x509: ECDSA verification failure")
		}
	default:
		err = cert.CheckSignature(
			cert.SignatureAlgorithm,
			digSig.Hash,
			digSig.Assinatura,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// Entidade responsável por gerar o arquivo de assinatura de todos os arquivos de resultados da urna.
// Podendo ter dois tipos de assinatura (Hardware (HW) e Software (SW)).
// Esses arquivos são informados na Mídia de Resultado quando a urna eletrônica é encerrada.
type EntidadeAssinaturaResultado struct {
	ModeloUrna   asn1.Enumerated    // Modelo da urna eletrônica.
	AssinaturaSW EntidadeAssinatura // Assinatura realizada via software (normalmente CEPESC).
	AssinaturaHW EntidadeAssinatura // Assinatura realizada via hardware de segurança da urna eletrônica.
}

func (e EntidadeAssinaturaResultado) Extension() string {
	return ".vscmr"
}

// Demais SEQUENCES
// Informações do algoritmo de hash.
// Informações do algoritmo de assinatura .
type AlgoritmoAssinaturaInfo struct {
	Algoritmo asn1.Enumerated // Tipo do algoritmo de assinatura.
	Bits      int             // Tamanho da assinatura.
}

type AlgoritmoHashInfo struct {
	Algoritmo asn1.Enumerated // Tipo do algoritmo de hash.
}

// Informações dos arquivos assinados.
type Assinatura struct {
	ArquivosAssinados []AssinaturaArquivo // Lista com Informações dos arquivos assinados.
}

// Informações do arquivo e da assinatura.
type AssinaturaArquivo struct {
	NomeArquivo string            // Nome do arquivo.
	Assinatura  AssinaturaDigital // Assinatura digital do arquivo.
}

// Informações da assinatura digital
type AssinaturaDigital struct {
	Tamanho    int    // Tamanho da assinatura.
	Hash       []byte // Hash da assinatura (Deve ser calculado uma única vez e ser utilizado também para o cálculo da assinatura).
	Assinatura []byte `asn1:"lax"` // Assinatura (Gerado/verificado a partir do hash acima).
}

// Informações da auto assinatura digital.
type AutoAssinaturaDigital struct {
	Usuario             DescritorChave          // Nome do usuário (Geralmente uma seção) que realizou a assinatura do arquivo.
	AlgoritmoHash       AlgoritmoHashInfo       // Algoritmo de hash utilizado para realizar a assinatura (Será o mesmo para as assinaturas de arquivos).
	AlgoritmoAssinatura AlgoritmoAssinaturaInfo // Algoritmo utilizado para realizar a assinatura (Será o mesmo para as assinaturas de arquivos).
	Assinatura          AssinaturaDigital       // Informações da assinatura digital.
}

// Identificador com informações da assinatura.
type DescritorChave struct {
	NomeUsuario string // Nome do usuário (Geralmente uma seção) que realizou a assinatura no arquivo.
	Serial      int    // Data em que foi gerado o conjunto de chaves.
}
