// Cf. https://www.tse.jus.br/eleicoes/eleicoes-2022/documentacao-tecnica-do-software-da-urna-eletronica

package ue

import "encoding/asn1"

// ENUMS
// Tipos de algoritmos de assinatura (cepesc é o algoritmo padrão (ainda não há previsão de uso dos demais)).

type AlgoritmoAssinatura byte

const (
	Rsa    AlgoritmoAssinatura = 1
	Ecdsa  AlgoritmoAssinatura = 2
	Cepesc AlgoritmoAssinatura = 3
)

// Tipos de algoritmos de hash (Todos os algoritmos devem ser suportados mas sha512 é o padrão).
type AlgoritmoHash byte

const (
	Sha1   AlgoritmoHash = 1
	Sha256 AlgoritmoHash = 2
	Sha384 AlgoritmoHash = 3
	Sha512 AlgoritmoHash = 4
)

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

// Entidade responsável por gerar o arquivo de assinatura de todos os arquivos de resultados da urna.
// Podendo ter dois tipos de assinatura (Hardware (HW) e Software (SW)).
// Esses arquivos são informados na Mídia de Resultado quando a urna eletrônica é encerrada.
type EntidadeAssinaturaResultado struct {
	ModeloUrna   asn1.Enumerated    // Modelo da urna eletrônica.
	AssinaturaSW EntidadeAssinatura // Assinatura realizada via software (normalmente CEPESC).
	AssinaturaHW EntidadeAssinatura // Assinatura realizada via hardware de segurança da urna eletrônica.
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
	Assinatura []byte // Assinatura (Gerado/verificado a partir do hash acima).
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
