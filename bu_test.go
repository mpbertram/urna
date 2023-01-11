package urna

import (
	"encoding/asn1"
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func ReadBu(file string) (EntidadeBoletimUrna, error) {
	f, err := os.ReadFile(file)
	if err != nil {
		return EntidadeBoletimUrna{}, errors.New("invalid input file")
	}

	var e EntidadeEnvelopeGenerico
	_, err = asn1.Unmarshal(f, &e)
	if err != nil {
		return EntidadeBoletimUrna{}, errors.New("error unmarshalling input file")
	}

	if TipoEnvelope(e.TipoEnvelope) != EnvelopeBoletimUrna {
		return EntidadeBoletimUrna{}, errors.New("envelope is not a bu")
	}

	var b EntidadeBoletimUrna
	_, err = asn1.Unmarshal(e.Conteudo, &b)
	if err != nil {
		return EntidadeBoletimUrna{}, errors.New("error unmarshalling bu from envelope")
	}

	return b, nil
}

func ComputeVotes(b EntidadeBoletimUrna) map[string]map[string]int {
	votosPorCargo := make(map[string]map[string]int)
	for _, cargo := range ValidCargoConstitucional() {
		votosPorCargo[cargo.String()] = map[string]int{}
	}

	for _, votacaoPorEleicao := range b.ResultadosVotacaoPorEleicao {
		for _, votacao := range votacaoPorEleicao.ResultadosVotacao {
			for _, votoCargo := range votacao.TotaisVotosCargo {
				for _, votoVotavel := range votoCargo.VotosVotaveis {
					cc, _ := CargoConstitucionalFromData(votoCargo.CodigoCargo.Bytes)

					switch TipoVoto(votoVotavel.TipoVoto) {
					case Nominal, Legenda:
						votosPorCargo[cc.String()][fmt.Sprint(votoVotavel.IdentificacaoVotavel.Codigo)] = votoVotavel.QuantidadeVotos
					case Branco:
						votosPorCargo[cc.String()][Branco.String()] = votoVotavel.QuantidadeVotos
					case Nulo:
						votosPorCargo[cc.String()][Nulo.String()] = votoVotavel.QuantidadeVotos
					}
				}
			}
		}
	}

	return votosPorCargo
}

func Test(t *testing.T) {
	bu, err := ReadBu("test-data/urna.bu")
	if err != nil {
		t.Error("could not read BU")
	}

	fmt.Println(ComputeVotes(bu))

	d, err := bu.ReadDadosSecaoSA()
	if err != nil {
		t.Error("could not read DadosSecao", bu)
	}

	if reflect.TypeOf(d) != reflect.TypeOf(DadosSecao{}) {
		t.Error("not dados secao", d)
	}

	if d.(DadosSecao).DataHoraAbertura != "20221002T080001" {
		t.Error("wrong DataHoraAbertura", d.(DadosSecao).DataHoraAbertura)
	}
	if d.(DadosSecao).DataHoraEncerramento != "20221002T170204" {
		t.Error("wrong DataHoraEncerramento", d.(DadosSecao).DataHoraEncerramento)
	}

	i, err := bu.Urna.CorrespondenciaResultado.ReadIdentificacao()
	if err != nil {
		t.Error("could not read Identificacao", bu.Urna.CorrespondenciaResultado)
	}

	if reflect.TypeOf(i) != reflect.TypeOf(IdentificacaoSecaoEleitoral{}) {
		t.Error("not identificacao secao eleitoral")
	}

	if i.(IdentificacaoSecaoEleitoral).MunicipioZona.Municipio != 88986 {
		t.Error("wrong municipio", i.(IdentificacaoSecaoEleitoral).MunicipioZona.Municipio)
	}
	if i.(IdentificacaoSecaoEleitoral).MunicipioZona.Zona != 7 {
		t.Error("wrong zona", i.(IdentificacaoSecaoEleitoral).MunicipioZona.Zona)
	}
	if i.(IdentificacaoSecaoEleitoral).Local != 1 {
		t.Error("wrong local", i.(IdentificacaoSecaoEleitoral).Local)
	}
	if i.(IdentificacaoSecaoEleitoral).Secao != 55 {
		t.Error("wrong local", i.(IdentificacaoSecaoEleitoral).Secao)
	}
}
