package urna

import (
	"encoding/asn1"
	"errors"
	"fmt"
	"os"
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

	fmt.Println(votosPorCargo)
	return votosPorCargo
}

func Test(t *testing.T) {
	bu, err := ReadBu("test-data/urna.bu")
	checkError(err)
	ComputeVotes(bu)
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
