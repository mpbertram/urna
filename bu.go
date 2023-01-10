// Cf. https://www.tse.jus.br/eleicoes/eleicoes-2022/documentacao-tecnica-do-software-da-urna-eletronica

package urna

import (
	"encoding/asn1"
	"errors"
)

type CodigoMunicipio int              // Código do município fornecido pelo cadastro da Justiça Eleitoral.
type DataHoraJE string                // Data e hora utilizada pela Justiça Eleitoral no formato YYYYMMDDThhmmss.
type IDEleicao int                    // Código numérico identificador da <glossario id='eleicao'>eleição</glossario> (Atribuído pelo Sistema Configurador de Eleições).
type IDPleito int                     // Código numérico identificador do <glossario id='pleito'>pleito</glossario> (Atribuído pelo Sistema Configurador de Eleições).
type IDProcessoEleitoral int          // Código numérico identificador do <glossario id='processo-eleitoral'>processo eleitoral</glossario> (Atribuído pelo Sistema Configurador de Eleções).
type NumeroCargoConsultaLivre int     // Número livre de cargo ou consulta definido no cadastramento da <glossario id='eleicao'>eleição</glossario>.
type NumeroInternoUrna int            // Número interno da urna eletrônica.
type NumeroLocal int                  // Número do local de votação da <glossario id='secao-eleitoral'>seção eleitoral</glossario> de acordo com o cadastro da Justiça Eleitoral.
type NumeroMesa int                   // Número da mesa de justificativa de acordo com o cadastro da Justiça Eleitoral (Informação referente ao Número da mesa utilizada para justificativa dos eleitores que não irão votar no seu domicílio eleitoral).
type NumeroPartido int                // Número do <glossario id='partido'>partido</glossario> fornecido pelo Sistema de Candidaturas da Justiça Eleitoral (Número do partido que compõe a <glossario id='coligacao'>coligação</glossario> ou do partido isolado).
type NumeroSecao int                  // Número da <glossario id='secao-eleitoral'>seçõe eleitoral</glossario> de acordo com o cadastro da Justiça Eleitoral.
type NumeroSerieFlash asn1.RawContent // Número de série da Flash (Representa um número de 4 bytes (0..2^32-1)).
type NumeroUrna int                   // Número da urna utilizada na mesa de justificativa de acordo com o cadastro da Justiça Eleitoral.
type NumeroVotavel int                // Número do <glossario id='votavel'>votável</glossario> fornecido pelo Sistema de Candidaturas da Justiça Eleitoral.
type NumeroZona int                   // Número da <glossario id='zona-eleitoral'>zona eleitoral</glossario> fornecido pelo cadastro da Justiça Eleitoral.

type CargoConstitucional byte

const (
	Presidente                  CargoConstitucional = 0x01
	VicePresidente              CargoConstitucional = 0x02
	Governador                  CargoConstitucional = 0x03
	ViceGovernador              CargoConstitucional = 0x04
	Senador                     CargoConstitucional = 0x05
	DeputadoFederal             CargoConstitucional = 0x06
	DeputadoEstadual            CargoConstitucional = 0x07
	DeputadoDistrital           CargoConstitucional = 0x08
	PrimeiroSuplenteSenador     CargoConstitucional = 0x09
	SegundoSuplenteSenador      CargoConstitucional = 0x0a
	Prefeito                    CargoConstitucional = 0x0b
	VicePrefeito                CargoConstitucional = 0x0c
	Vereador                    CargoConstitucional = 0x0d
	CargoConstitucionalInvalido CargoConstitucional = 0xff
)

func CargoConstitucionalFromData(data []byte) (CargoConstitucional, error) {
	if len(data) == 1 {
		switch v := data[0]; v {
		case 0x01:
			return Presidente, nil
		case 0x02:
			return VicePresidente, nil
		case 0x03:
			return Governador, nil
		case 0x04:
			return ViceGovernador, nil
		case 0x05:
			return Senador, nil
		case 0x06:
			return DeputadoFederal, nil
		case 0x07:
			return DeputadoEstadual, nil
		case 0x08:
			return DeputadoDistrital, nil
		case 0x09:
			return PrimeiroSuplenteSenador, nil
		case 0x0a:
			return SegundoSuplenteSenador, nil
		case 0x0b:
			return Prefeito, nil
		case 0x0c:
			return VicePrefeito, nil
		case 0x0d:
			return Vereador, nil
		}
	}

	return CargoConstitucionalInvalido, errors.New("invalid data")
}

func (cc CargoConstitucional) String() string {
	if cc <= 0x0d {
		return [...]string{
			"Presidente", "Vice Presidente", "Governador", "Vice Governador",
			"Senador", "Deputado Federal", "Deputado Estadual", "Deputado Distrital",
			"Primeiro Suplente de Senador", "Segundo Suplente de Senador",
			"Prefeito", "Vice Prefeito", "Vereador"}[cc-1]
	}

	return "Cargo Constitucional invalido"
}

func ValidCargoConstitucional() []CargoConstitucional {
	return []CargoConstitucional{
		Presidente, VicePresidente, Governador, ViceGovernador,
		Senador, DeputadoFederal, DeputadoEstadual, DeputadoDistrital,
		PrimeiroSuplenteSenador, SegundoSuplenteSenador,
		Prefeito, VicePrefeito, Vereador}
}

type Fase byte

const (
	Simulado     Fase = 0x01
	Oficial      Fase = 0x02
	Treinamento  Fase = 0x03
	FaseInvalida Fase = 0xff
)

func FaseFromData(data []byte) (Fase, error) {
	if len(data) == 1 {
		switch v := data[0]; v {
		case 0x01:
			return Simulado, nil
		case 0x02:
			return Oficial, nil
		case 0x03:
			return Treinamento, nil
		}
	}

	return FaseInvalida, errors.New("invalid data")
}

func (f Fase) String() string {
	if f <= 0x03 {
		return [...]string{"Simulado", "Oficial", "Treinamento"}[f-1]
	}

	return "Fase invalida"
}

type MotivoApuracaoManual byte

const (
	UrnaComDefeitoApMan    MotivoApuracaoManual = 0x01
	UrnaIndisponivelInicio MotivoApuracaoManual = 0x02
	UrnaOutraSecao         MotivoApuracaoManual = 0x03
	OutrosApMan            MotivoApuracaoManual = 0x63
)

func MotivoApuracaoManualFromData(data []byte) (MotivoApuracaoManual, error) {
	if len(data) == 1 {
		switch v := data[0]; v {
		case 0x01:
			return UrnaComDefeitoApMan, nil
		case 0x02:
			return UrnaIndisponivelInicio, nil
		case 0x03:
			return UrnaOutraSecao, nil
		default:
			return OutrosApMan, nil
		}
	}

	return OutrosApMan, errors.New("invalid data")
}

func (m MotivoApuracaoManual) String() string {
	if m <= 0x03 {
		return [...]string{"Urna com defeito", "Urna indisponivel no inicio", "Urna de outra secao"}[m-1]
	}

	return "Outro"
}

type MotivoApuracaoMistaComBU byte

const (
	UrnaDataHoraIncorreta       MotivoApuracaoMistaComBU = 0x01
	UrnaComDefeito              MotivoApuracaoMistaComBU = 0x02
	UrnaOutrasecao              MotivoApuracaoMistaComBU = 0x03
	UrnaPreparadaIncorretamente MotivoApuracaoMistaComBU = 0x04
	UrnaChegouAposInicioVotacao MotivoApuracaoMistaComBU = 0x05
	OutrosApMisBu               MotivoApuracaoMistaComBU = 0x63
)

func MotivoApuracaoMistaComBUFromData(data []byte) (MotivoApuracaoMistaComBU, error) {
	if len(data) == 1 {
		switch v := data[0]; v {
		case 0x01:
			return UrnaDataHoraIncorreta, nil
		case 0x02:
			return UrnaComDefeito, nil
		case 0x03:
			return UrnaOutrasecao, nil
		case 0x04:
			return UrnaPreparadaIncorretamente, nil
		case 0x05:
			return UrnaChegouAposInicioVotacao, nil
		default:
			return OutrosApMisBu, nil
		}
	}

	return OutrosApMisBu, errors.New("invalid data")
}

func (m MotivoApuracaoMistaComBU) String() string {
	if m <= 0x05 {
		return [...]string{"Urna com data/hora incorreta", "Urna com defeito", "Urna de outra secao", "Urna preparada incorretamente", "Urna chegou apos inicio da votacao"}[m-1]
	}

	return "Outro"
}

type MotivoApuracaoMistaComMR byte

const (
	NaoObteveExitoContingencia         MotivoApuracaoMistaComMR = 0x01
	IndisponibilidadeUrnaContingencia  MotivoApuracaoMistaComMR = 0x02
	IndisponibilidadeFlashContingencia MotivoApuracaoMistaComMR = 0x03
	ProblemaEnergiaEletrica            MotivoApuracaoMistaComMR = 0x04
	NaoFoiPossivelTrocarUrna           MotivoApuracaoMistaComMR = 0x05
	NaoFoiSolicitadaTrocaUrna          MotivoApuracaoMistaComMR = 0x06
	OutrosApMisMr                      MotivoApuracaoMistaComMR = 0x63
)

func MotivoApuracaoMistaComMRFromData(data []byte) (MotivoApuracaoMistaComMR, error) {
	if len(data) == 1 {
		switch v := data[0]; v {
		case 0x01:
			return NaoObteveExitoContingencia, nil
		case 0x02:
			return IndisponibilidadeUrnaContingencia, nil
		case 0x03:
			return IndisponibilidadeFlashContingencia, nil
		case 0x04:
			return ProblemaEnergiaEletrica, nil
		case 0x05:
			return NaoFoiPossivelTrocarUrna, nil
		case 0x06:
			return NaoFoiSolicitadaTrocaUrna, nil
		default:
			return OutrosApMisMr, nil
		}
	}

	return OutrosApMisMr, errors.New("invalid data")
}

func (m MotivoApuracaoMistaComMR) String() string {
	if m <= 0x06 {
		return [...]string{"Nao obteve exito contingencia", "IndisponibilidadeUrnaContingencia", "IndisponibilidadeFlashContingencia", "ProblemaEnergiaEletrica", "NaoFoiPossivelTrocarUrna", "NaoFoiSolicitadaTrocaUrna"}[m-1]
	}

	return "Outro"
}

type TipoApuracao byte

const (
	TotalmenteManual     TipoApuracao = 0x01
	TotalmenteEletronica TipoApuracao = 0x02
	MistaBU              TipoApuracao = 0x03
	MistaMR              TipoApuracao = 0x04
	TipoApuracaoInvalida TipoApuracao = 0xff
)

func TipoApuracaoFromData(data []byte) (TipoApuracao, error) {
	if len(data) == 1 {
		switch v := data[0]; v {
		case 0x01:
			return TotalmenteManual, nil
		case 0x02:
			return TotalmenteEletronica, nil
		case 0x03:
			return MistaBU, nil
		case 0x04:
			return MistaMR, nil
		}
	}

	return TipoApuracaoInvalida, errors.New("invalid data")
}

func (t TipoApuracao) String() string {
	if t <= 0x04 {
		return [...]string{"TotalmenteManual", "TotalmenteEletronica", "MistaBU", "MistaMR"}[t-1]
	}

	return "Invalido"
}

// Tipos de arquivos de votação.
type TipoArquivo byte

const (
	VotacaoUE               TipoArquivo = 0x01 // Urna eletrônica.
	VotacaoRED              TipoArquivo = 0x02 // RED (Recuperador de Dados - Responsável por gerar uma nova memória de resultado a partir da urna originária).
	SaMistaMRParcialCedula  TipoArquivo = 0x03 // <glossario id='sistema-de-apuracao'>Sistema de Apuração</glossario> (Votação Mista - Memória de Resultado e Cédulas).
	SaMistaBUImpressoCedula TipoArquivo = 0x04 // <glossario id='sistema-de-apuracao'>Sistema de Apuração</glossario> (Votação Mista - <glossario id='boletim-de-urna'>BU</glossario> impresso e Cédulas).
	SaManual                TipoArquivo = 0x05 // <glossario id='sistema-de-apuracao'>Sistema de Apuração</glossario> (Votação totalmente manual - Cédulas).
	SaEletronica            TipoArquivo = 0x06 // <glossario id='sistema-de-apuracao'>Sistema de Apuração</glossario> (Votação totalmente eletrônica).
	TipoArquivoInvalido     TipoArquivo = 0xff
)

func TipoArquivoFromData(data []byte) (TipoArquivo, error) {
	if len(data) == 1 {
		switch v := data[0]; v {
		case 0x01:
			return VotacaoUE, nil
		case 0x02:
			return VotacaoRED, nil
		case 0x03:
			return SaMistaMRParcialCedula, nil
		case 0x04:
			return SaMistaBUImpressoCedula, nil
		case 0x05:
			return SaManual, nil
		case 0x06:
			return SaEletronica, nil
		}
	}

	return TipoArquivoInvalido, errors.New("invalid data")
}

func (t TipoArquivo) String() string {
	if t <= 0x06 {
		return [...]string{"VotacaoUE", "VotacaoRED", "SaMistaMRParcialCedula", "SaMistaBUImpressoCedula", "SaManual", "SaEletronica"}[t-1]
	}

	return "Invalido"
}

// Tipos de cargos ou consultas da <glossario id='eleicao'>eleição</glossario>.
type TipoCargoConsulta byte

const (
	Majoritario               TipoCargoConsulta = 0x01 // <glossario id='cargo-majoritario'>Cargos majoritários</glossario>.
	Proporcional              TipoCargoConsulta = 0x02 // <glossario id='cargo-proporcional'>Cargos proporcionais</glossario>.
	Consulta                  TipoCargoConsulta = 0x03 // São as perguntas da <glossario id='consulta-popular'>consulta popular</glossario>.
	TipoCargoConsultaInvalido TipoCargoConsulta = 0xff
)

func TipoCargoConsultaFromData(data []byte) (TipoCargoConsulta, error) {
	if len(data) == 1 {
		switch v := data[0]; v {
		case 0x01:
			return Majoritario, nil
		case 0x02:
			return Proporcional, nil
		case 0x03:
			return Consulta, nil
		}
	}

	return TipoCargoConsultaInvalido, errors.New("invalid data")
}

func (t TipoCargoConsulta) String() string {
	if t <= 0x03 {
		return [...]string{"Majoritario", "Proporcional", "Consulta"}[t-1]
	}

	return "Invalido"
}

// Tipos de envelopes dos arquivos.
type TipoEnvelope byte

const (
	EnvelopeBoletimUrna         TipoEnvelope = 0x01 // Boletim de Urna (BU).
	EnvelopeRegistroDigitalVoto TipoEnvelope = 0x02 // Registro Digital do Voto (RDV).
	EnvelopeBoletimUrnaImpresso TipoEnvelope = 0x04 // Boletim de Urna impresso.
	EnvelopeImagemBiometria     TipoEnvelope = 0x05 // Arquivo de imagem de biometria.
	EnvelopeInvalido            TipoEnvelope = 0xff
)

func TipoEnvelopeFromData(data []byte) (TipoEnvelope, error) {
	if len(data) == 1 {
		switch v := data[0]; v {
		case 0x01:
			return EnvelopeBoletimUrna, nil
		case 0x02:
			return EnvelopeRegistroDigitalVoto, nil
		case 0x03:
			return EnvelopeBoletimUrnaImpresso, nil
		case 0x04:
			return EnvelopeImagemBiometria, nil
		}
	}

	return EnvelopeInvalido, errors.New("invalid data")
}

func (t TipoEnvelope) String() string {
	if t <= 0x04 {
		return [...]string{"EnvelopeBoletimUrna", "EnvelopeRegistroDigitalVoto", "EnvelopeBoletimUrnaImpresso", "EnvelopeImagemBiometria"}[t-1]
	}

	return "Invalido"
}

// Tipos de urna eletrônica.
type TipoUrna byte

const (
	Secao                  TipoUrna = 0x01 // Urna de seção.
	Contingencia           TipoUrna = 0x03 // <glossario id='urna-de-contingencia'>Urna de contingência</glossario>.
	ReservaSecao           TipoUrna = 0x04 // Resultado de urna de contingência que passou a ser de seção.
	ReservaEncerrandoSecao TipoUrna = 0x06 // Barriga de aluguel para seção (Procedimento de recuperação dos dados de uma urna de seção, a partir da inserção de seu cartão de memória externo em uma <glossario id='urna-de-contingencia'>urna de contingência</glossario>).
	UrnaInvalida           TipoUrna = 0xff
)

func TipoUrnaFromData(data []byte) (TipoUrna, error) {
	if len(data) == 1 {
		switch v := data[0]; v {
		case 0x01:
			return Secao, nil
		case 0x02:
			return Contingencia, nil
		case 0x03:
			return ReservaSecao, nil
		case 0x04:
			return ReservaEncerrandoSecao, nil
		}
	}

	return UrnaInvalida, errors.New("invalid data")
}

func (t TipoUrna) String() string {
	if t <= 0x04 {
		return [...]string{"Secao", "Contingencia", "ReservaSecao", "ReservaEncerrandoSecao"}[t-1]
	}

	return "Invalido"
}

// Tipos de votos existentes na urna eletrônica.
type TipoVoto byte

const (
	Nominal           TipoVoto = 0x01 // <glossario id='votos-nominais'>Voto nominal.</glossario>
	Branco            TipoVoto = 0x02 // Voto branco.
	Nulo              TipoVoto = 0x03 // Voto nulo.
	Legenda           TipoVoto = 0x04 // <glossario id='votos-de-legenda'>Voto de legenda.</glossario>
	CargoSemCandidato TipoVoto = 0x05 // Nenhum candidato para ser votado no cargo.
	TipoVotoInvalido  TipoVoto = 0xff
)

func TipoVotoFromData(data []byte) (TipoVoto, error) {
	if len(data) == 1 {
		switch v := data[0]; v {
		case 0x01:
			return Nominal, nil
		case 0x02:
			return Branco, nil
		case 0x03:
			return Nulo, nil
		case 0x04:
			return Legenda, nil
		case 0x05:
			return CargoSemCandidato, nil
		}
	}

	return TipoVotoInvalido, errors.New("invalid data")
}

func (t TipoVoto) String() string {
	if t <= 0x06 {
		return [...]string{"Nominal", "Branco", "Nulo", "Legenda", "CargoSemCandidato"}[t-1]
	}

	return "Invalido"
}

// Identificador que contém informações do cabeçalho da entidade (Arquivos ASN.1).
type CabecalhoEntidade struct {
	DataGeracao DataHoraJE    // Data da geração da entidade.
	IdEleitoral asn1.RawValue // Identificador Eleitoral (<glossario id='processo'>Processo</glossario> <glossario id='pleito'>pleito</glossario> ou <glossario id='eleicao'>eleição</glossario>).
}

// Identificador com informações da urna eletrônica.
type Urna struct {
	TipoUrna                 asn1.Enumerated          // Tipo da urna eletrônica.
	VersaoVotacao            string                   // Versão do software de votação da urna eletrônica.
	CorrespondenciaResultado CorrespondenciaResultado // Informações da <glossario id='correspondencia'>correspondência</glossario> da urna eletrônica.
	TipoArquivo              asn1.Enumerated          // Tipo do arquivo gerado pela urna eletrônica.
	NumeroSerieFV            NumeroSerieFlash         // Número de série da Flash de Votação.
	MotivoUtilizacaoSA       asn1.RawValue            `asn1:"optional"` // Identificador numérico para o motivo de utilização do <glossario id='sistema-de-apuracao'>Sistema de Apuração</glossario> para a urna eletrônica.
}

// Envelope
// Entidade responsável por envelopar os arquivos ou dados binários da urna eletrônica.
// Transforma os arquivos da urna eletrônica em arquivos no padrão ASN.1 assinados e algumas vezes criptografados.
type EntidadeEnvelopeGenerico struct {
	Cabecalho     CabecalhoEntidade // Informações do cabeçalho da entidade.
	Fase          asn1.Enumerated   // Fase em que foi gerado o arquivo.
	Urna          Urna              `asn1:"optional"` // Informações da urna eletrônica (Deve existir para RDV e ser omitido no BU).
	Identificacao asn1.RawValue     // Identificação se é urna de seção eleitoral ou de Mesa Receptora de Justificativa.
	TipoEnvelope  asn1.Enumerated   // Tipo de envelope que será criado.
	Seguranca     Seguranca         `asn1:"optional"` // Informações de segurança solicitados pela biblioteca do CEPESC (Existindo o conteúdo estará cifrado).
	Conteudo      []byte            // Conteúdo do envelope gerado.
}

// Identificador com informações de segurança solicitados pela biblioteca do CEPESC.
type Seguranca struct {
	IdTipoArquivo  int    // Identificador que corresponde ao arquivo solicitado pela biblioteca do CEPESC.
	IdCriptografia int    // Identificador que corresponde ao Turno solicitado pela biblioteca do CEPESC.
	IdArquivoCD    int    // Identificador do arquivo solicitado pela biblioteca do CEPESC.
	IdArquivoChave []byte // Chave pública para cifrar o arquivo.
}

// SEQUENCE RAIZ
// Entidade responsável por apresentar as informações do <glossario id='boletim-de-urna'>boletim de urna</glossario>.
type EntidadeBoletimUrna struct {
	Cabecalho                   CabecalhoEntidade            // Informações do cabeçalho da entidade.
	Fase                        asn1.Enumerated              // Fase em que foi gerado o arquivo.
	Urna                        Urna                         // Informações da urna eletrônica.
	IdentificacaoSecao          IdentificacaoSecaoEleitoral  // Informações da <glossario id='secao-eleitoral'>seção eleitoral</glossario> que está instalada a urna eletrônica.
	DataHoraEmissao             DataHoraJE                   // Data e hora da emissão do boletim de urna.
	DadosSecaoSA                asn1.RawValue                // Identificação para resultado de urna de seção ou <glossario id='sistema-de-apuracao'>de Sistema de Apuração</glossario>.
	QtdEleitoresLibCodigo       int                          `asn1:"tag:1,optional"` // Quantidade de eleitores que compareceram que foram habilitados manualmente.
	QtdEleitoresCompBiometrico  int                          `asn1:"tag:2,optional"` // Quantidade de eleitores que compareceram que utilizaram <glossario id='identificacao-biometrica'>biometria</glossario>.
	ResultadosVotacaoPorEleicao []ResultadoVotacaoPorEleicao `asn1:"tag:3"`          // Lista com os resultados da votação para cada eleição.
	HistoricoCorrespondencias   []CorrespondenciaResultado   `asn1:"tag:4,optional"` // Lista com informações de histórico das <glossario id='correspondencia'>correspondências</glossario> (Pode ser opcional porque quando o BU é da urna original não existe esse histórico).
	HistoricoVotoImpresso       []HistoricoVotoImpresso      `asn1:"tag:5,optional"` // Lista com informações de histórico de voto impresso.
	ChaveAssinaturaVotosVotavel []byte                       // Chave de assinatura pública das tuplas dos votáveis.
}

// DEMAIS SEQUENCES E CHOICES (ordem alfabética)
type ApuracaoEletronica struct {
	Tipoapuracao   asn1.Enumerated
	MotivoApuracao asn1.Enumerated
}

type ApuracaoMistaBUAE struct {
	Tipoapuracao   asn1.Enumerated
	MotivoApuracao asn1.Enumerated
}

type ApuracaoMistaMR struct {
	TipoApuracao   asn1.Enumerated
	MotivoApuracao asn1.Enumerated
}

type ApuracaoTotalmenteManualDigitacaoAE struct {
	Tipoapuracao   asn1.Enumerated
	MotivoApuracao asn1.Enumerated
}

// Identificador com informações da carga da urna eletrônica.
type Carga struct {
	NumeroInternoUrna NumeroInternoUrna // Número interno da urna eletrônica.
	NumeroSerieFC     NumeroSerieFlash  // Número de série da unidade de Flash Card.
	DataHoraCarga     DataHoraJE        // Data e hora da carga no formato utilizado pela Justiça Eleitoral (YYYYMMDDThhmmss).
	CodigoCarga       string            // Código da carga da urna eletrônica.
}

// Identificador com informações da urna e da carga.
type CorrespondenciaResultado struct {
	Identificacao asn1.RawValue // Identificação se  tem carga de seção ou de mesa receptora de justificativa.
	Carga         Carga         // Informações da carga da urna eletrônica.
}

// Identificador com informações do <glossario id='boletim-de-urna'>BU</glossario>) de <glossario id='sistema-de-apuracao'>SA</glossario>).
type DadosSA struct {
	JuntaApuradora          int               // Número da junta eleitoral responsával pela apuração dos votos.
	TurmaApuradora          int               // Número da turma apuradora responsával pela apuração dos votos.
	NumeroInternoUrnaOrigem NumeroInternoUrna // Número interno da urna eletrônica com impossibilidade de utilização.
}

// Identificador com informações do <glossario id='boletim-de-urna'>BU</glossario>) de seção.
type DadosSecao struct {
	DataHoraAbertura                 DataHoraJE // Data e hora do início da aquisição do voto (Primeiro voto) no formato adotado pela Justiça Eleitoral (YYYYMMDDThhmmss).
	DataHoraEncerramento             DataHoraJE // Data e hora do término da aquisição do voto (Último voto) no formato adotado pela Justiça Eleitoral (YYYYMMDDThhmmss).
	DataHoraDesligamentoVotoImpresso DataHoraJE // Data e hora do desligamento da impressão do voto (somente se tinha voto impresso na seção e se ocorreu o cancelamento) (YYYYMMDDThhmmss).
}

// Identificador com informações de histórico de voto impresso
type HistoricoVotoImpresso struct {
	IdImpressoraVotos  int        // Número interno da impressora de votos
	IdRepositorioVotos int        // Número interno do repositório de votos
	DataHoraLigamento  DataHoraJE // Data e hora do momento que o dispositivo for ligado
}

// Identificador com informações de <glossario id='contingencia'>contingência</glossario>.
type IdentificacaoContingencia struct {
	MunicipioZona MunicipioZona // Número do município e Número da <glossario id='zona-eleitoral'>zona eleitoral</glossario> a qual pertence a urna.
}

// Identificador com informações da mesa receptora de justificativa.
type IdentificacaoMesaJustificativa struct {
	MunicipioZona MunicipioZona // Número do município e Número da <glossario id='zona-eleitoral'>zona eleitoral</glossario>.
	Mesa          NumeroMesa    // Número da mesa de justificativa.
	Urna          NumeroUrna    // Número da urna de justificativa.
}

// Identificador com informações da <glossario id='secao-eleitoral'>seção eleitoral</glossario>.
type IdentificacaoSecaoEleitoral struct {
	MunicipioZona MunicipioZona // Número do município e Número da <glossario id='zona-eleitoral'>zona eleitoral</glossario> a qual pertence a <glossario id='secao-eleitoral'>seção eleitoral</glossario>.
	Local         NumeroLocal   // Número do local de votação da seção eleitoral.
	Secao         NumeroSecao   // Número identificador da <glossario id='secao-eleitoral'>seção eleitoral</glossario>.
}

// Identificação de um votável que pode ser um candidato ou uma pergunta de consulta popular.
type IdentificacaoVotavel struct {
	Partido NumeroPartido // Número do partido.
	Codigo  NumeroVotavel // Número do votável.
}

// Identificador que contém informações de município e <glossario id='zona eleitoral'>zona eleitoral</glossario> que são relacionados entre si.
type MunicipioZona struct {
	Municipio CodigoMunicipio // Código do município de acordo com o cadastro da Justiça Eleitoral.
	Zona      NumeroZona      // Número da <glossario id='zona eleitoral'>zona eleitoral</glossario> de acordo com o cadastro da Justiça Eleitoral.
}

// Identificador com informações do resultado de votação da urna eletrônica.
type ResultadoVotacao struct {
	TipoCargo         asn1.Enumerated   // Tipo do cargo ou consulta.
	QtdComparecimento int               // Quantidade de eleitores que compareceram à seção para votação no cargo ou consulta.
	TotaisVotosCargo  []TotalVotosCargo // Quantidade total de votos para cada cargo ou consulta.
}

// Estrutura com os resultados da votação de uma eleição.
type ResultadoVotacaoPorEleicao struct {
	IdEleicao         IDEleicao          // Identificador numérico da <glossario id='eleicao'>eleição</glossario>.
	QtdEleitoresAptos int                // Quantidade de <glossario id='eleitor'>eleitores</glossario> aptos a votar na urna eletrônica da seção.
	ResultadosVotacao []ResultadoVotacao // Lista com informações do resultado da votação na urna eletrônica.
}

// Identificador com informações do total de votos para cada cargo ou consulta.
type TotalVotosCargo struct {
	CodigoCargo    asn1.RawValue       // Código do cargo ou da consulta.
	OrdemImpressao int                 // Ordem para impressão dos cargos ou consultas no <glossario id='voto-em-transito'>boletim de urna</glossario> e demais relatórios utilizados na Justiça Eleitoral.
	VotosVotaveis  []TotalVotosVotavel // Informações do total de votos agrupados por tipo de voto e número do <glossario id='votavel'>votável</glossario>.
}

// Identificador com informações da quantidade de votos agrupados por tipo de voto e número do <glossario id='votavel'>votável</glossario>.
type TotalVotosVotavel struct {
	TipoVoto             asn1.Enumerated      `asn1:"tag:1"`          // Tipo do voto.
	QuantidadeVotos      int                  `asn1:"tag:2"`          // Quantidade de votos por tipo e número do votável.
	IdentificacaoVotavel IdentificacaoVotavel `asn1:"tag:3,optional"` // Identificação do votável (Para tipo de voto "Branco" ou "Nulo" esse campo deverá ser omitido).
	Assinatura           []byte               // Assinatura dos dados compostos de votos do votável. Os seguintes campos são assinados: TotalVotosCargo::codigoCargo TotalVotosVotavel::tipoVoto TotalVotosVotavel::quantidadeVotos identificacaoVotavel::codigo identificacaoVotavel::partido Carga::codigoCarga
}
