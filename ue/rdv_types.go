// Cf. https://www.tse.jus.br/eleicoes/eleicoes-2022/documentacao-tecnica-do-software-da-urna-eletronica

package ue

import (
	"encoding/asn1"
	"errors"
)

// TIPOS
type QuantidadeEscolhas int // Quantidade máxima de escolhas para um mesmo cargo.
type VotoDigitado string    // Digitação como feita pelo eleitor na urna.

// ENUMS
type MotivoApuracaoEletronica byte

const (
	NaoFoiPossivelReuperarResultado MotivoApuracaoEletronica = 1
	UrnaNaoChegouMidiaDefeituosa    MotivoApuracaoEletronica = 2
	UrnaNaoChegouMidiaExtraviada    MotivoApuracaoEletronica = 3
	Outros                          MotivoApuracaoEletronica = 99
)

// Origem dos votos inseridos no SA.
type OrigemVotosSA byte

const (
	Cedula OrigemVotosSA = 1
	Rdv    OrigemVotosSA = 2
	Bu     OrigemVotosSA = 3
)

// Tipo do sistema eleitoral.
type TipoCedulaSA byte

const (
	CedulaSAMajoritario  TipoCedulaSA = 1
	CedulaSAProporcional TipoCedulaSA = 2
)

// SEQUENCES e CHOICES
// Entidade usada para a geração do RDV na memória de resultado.
type EntidadeResultadoRDV struct {
	Cabecalho CabecalhoEntidade           // Informações do cabeçalho da entidade.
	Urna      Urna                        // Informações da urna eletrônica.
	Rdv       EntidadeRegistroDigitalVoto // Registro digital do voto.
}

// Entidade usada para o armazenamento do RDV nas mídias interna e externa da urna.
type EntidadeRegistroDigitalVoto struct {
	Pleito        IDPleito                    // Identificador do pleito corrente.
	Fase          asn1.Enumerated             // Fase em que foi gerado o arquivo.
	Identificacao IdentificacaoSecaoEleitoral // Identificação da seção eleitoral.
	Eleicoes      asn1.RawValue               // Grupo de votos de todas as eleições.
}

// Result is one of ([]EleicaoVota, []EleicaoSA)
func (rdv EntidadeRegistroDigitalVoto) ReadEleicoes() (interface{}, error) {
	switch rdv.Eleicoes.Tag {
	case 0:
		var e []EleicaoVota
		err := FillSlice(rdv.Eleicoes.Bytes, &e)
		if err != nil {
			return nil, err
		}
		return e, nil
	case 1:
		var e []EleicaoSA
		err := FillSlice(rdv.Eleicoes.Bytes, &e)
		if err != nil {
			return nil, err
		}
		return e, nil
	}

	return nil, errors.New("could not read dados secao/SA")
}

// Votos para todos os cargos de uma eleição.
type EleicaoVota struct {
	IdEleicao   int          // Identificador da eleição.
	VotosCargos []VotosCargo // Grupo de cédulas da eleição.
}

// Votos para todos os cargos de uma eleição.
type EleicaoSA struct {
	IdEleicao     IDEleicao     // Identificador da eleição.
	TipoCedulaSA  TipoCedulaSA  // Tipo da cédula de papel apurada pelo SA.
	OrigemVotosSA OrigemVotosSA // A origem dos votos inseridos no SA.
	VotosCargos   []VotosCargo  // Grupo de cédulas da eleição.
}

// Votos de um eleitor para todas as escolhas de um cargo.
type Voto struct {
	TipoVoto  asn1.Enumerated // Tipo do voto registrado.
	Digitacao VotoDigitado    `asn1:"optional"` // Número como digitado pelo eleitor (não existe para TipoVoto = 3, 5, 6, 8 e 9).
}

// Todos os votos para um cargo específico.
type VotosCargo struct {
	IdCargo            asn1.RawValue      // Código do cargo votado.
	QuantidadeEscolhas QuantidadeEscolhas // Quantidade de escolhas para o cargo.
	Votos              []Voto             // Votos do cargo.
}

// Result is one of CargoConstitucional or NumeroCargoConsultaLivre
func (vc VotosCargo) ReadIdCargo() (interface{}, error) {
	switch vc.IdCargo.Tag {
	case 1:
		cc, err := CargoConstitucionalFromData(vc.IdCargo.Bytes)
		if err != nil {
			return nil, err
		}
		return cc, nil
	case 2:
		var n NumeroCargoConsultaLivre
		_, err := asn1.Unmarshal(vc.IdCargo.Bytes, &n)
		if err != nil {
			return nil, err
		}
		return n, nil
	}

	return nil, errors.New("could not read cargo")
}
