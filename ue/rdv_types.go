// Cf. https://www.tse.jus.br/eleicoes/eleicoes-2022/documentacao-tecnica-do-software-da-urna-eletronica

package ue

import (
	"errors"
	"github.com/google/certificate-transparency-go/asn1"
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

func MotivoApuracaoEletronicaFromData(data asn1.Enumerated) (MotivoApuracaoEletronica, error) {
	switch data {
	case 0x01:
		return NaoFoiPossivelReuperarResultado, nil
	case 0x02:
		return UrnaNaoChegouMidiaDefeituosa, nil
	case 0x03:
		return UrnaNaoChegouMidiaExtraviada, nil
	default:
		return Outros, nil
	}
}

func (m MotivoApuracaoEletronica) String() string {
	if m <= 0x03 {
		return [...]string{"Nao foi possivel recuperar resultado", "Urna nao chegou (midia defeituosa)", "Urna nao chegou (midia extraviada)"}[m-1]
	}

	return "Outro"
}

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

type TipoVotoRdv int

const (
	LegendaRdv                            TipoVotoRdv = 0x01 // <glossario id='votos-nominais'>Voto nominal.</glossario>
	NominalRdv                            TipoVotoRdv = 0x02 // Voto branco.
	BrancoRdv                             TipoVotoRdv = 0x03 // Voto nulo.
	NuloRdv                               TipoVotoRdv = 0x04 // <glossario id='votos-de-legenda'>Voto de legenda.</glossario>
	BrancoAposSuspensaoRdv                TipoVotoRdv = 0x05 // Nenhum candidato para ser votado no cargo.
	NuloAposSuspensaoRdv                  TipoVotoRdv = 0x06
	NuloPorRepeticaoRdv                   TipoVotoRdv = 0x07
	NuloCargoSemCandidatoRdv              TipoVotoRdv = 0x08
	NuloAposSuspensaoCargoSemCandidatoRdv TipoVotoRdv = 0x09
	TipoVotoInvalidoRdv                   TipoVotoRdv = 0xff
)

func TipoVotoRdvFromData(data int) (TipoVotoRdv, error) {

	switch v := data; v {
	case 0x01:
		return LegendaRdv, nil
	case 0x02:
		return NominalRdv, nil
	case 0x03:
		return BrancoRdv, nil
	case 0x04:
		return NuloRdv, nil
	case 0x05:
		return BrancoAposSuspensaoRdv, nil
	case 0x06:
		return NuloAposSuspensaoRdv, nil
	case 0x07:
		return NuloPorRepeticaoRdv, nil
	case 0x08:
		return NuloCargoSemCandidatoRdv, nil
	case 0x09:
		return NuloAposSuspensaoCargoSemCandidatoRdv, nil
	case 0xff:
		return TipoVotoInvalidoRdv, nil
	}

	return TipoVotoInvalidoRdv, errors.New("invalid data")
}

func (t TipoVotoRdv) String() string {
	if t <= 0x09 {
		return [...]string{
			"Legenda",
			"Nominal",
			"Branco",
			"Nulo",
			"Branco apos suspensao",
			"Nulo apos suspensao",
			"Nulo por repeticao",
			"Nulo cargo sem candidato",
			"Nulo apos suspensao cargo sem candidato"}[t-1]
	}

	return "Invalido"
}

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

func (e EntidadeResultadoRDV) Extension() string {
	return ".rdv"
}

type EleicaoGenerica interface {
	GetId() int
	GetVotosCargos() []VotosCargo
}

// Votos para todos os cargos de uma eleição.
type EleicaoVota struct {
	IdEleicao   int          // Identificador da eleição.
	VotosCargos []VotosCargo // Grupo de cédulas da eleição.
}

func (e EleicaoVota) GetId() int {
	return e.IdEleicao
}

func (e EleicaoVota) GetVotosCargos() []VotosCargo {
	return e.VotosCargos
}

// Votos para todos os cargos de uma eleição.
type EleicaoSA struct {
	IdEleicao     int             // Identificador da eleição.
	TipoCedulaSA  asn1.Enumerated // Tipo da cédula de papel apurada pelo SA.
	OrigemVotosSA asn1.Enumerated // A origem dos votos inseridos no SA.
	VotosCargos   []VotosCargo    // Grupo de cédulas da eleição.
}

func (e EleicaoSA) GetId() int {
	return e.IdEleicao
}

func (e EleicaoSA) GetVotosCargos() []VotosCargo {
	return e.VotosCargos
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
