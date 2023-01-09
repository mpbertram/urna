package urna

import (
	"encoding/asn1"
	"fmt"
	"os"
	"testing"
)

func Test(t *testing.T) {
	f, err := os.ReadFile("test-data/urna.bu")
	checkError(err)

	var envelope EntidadeEnvelopeGenerico
	_, err = asn1.Unmarshal(f, &envelope)
	checkError(err)

	switch envelopeType := envelope.TipoEnvelope; envelopeType {
	case 1:
		var boletim EntidadeBoletimUrna
		_, err := asn1.Unmarshal(envelope.Conteudo, &boletim)
		checkError(err)

		for _, votacaoPorEleicao := range boletim.ResultadosVotacaoPorEleicao {
			for _, votacao := range votacaoPorEleicao.ResultadosVotacao {
				// fmt.Println("Tipo cargo:", votacao.TipoCargo)
				for _, votoCargo := range votacao.TotaisVotosCargo {
					// fmt.Println("Codigo cargo:", votoCargo.CodigoCargo)
					for _, votoVotavel := range votoCargo.VotosVotaveis {
						// fmt.Println("Tipo voto:", votoVotavel.TipoVoto)
						if votoVotavel.TipoVoto == asn1.Enumerated(Nominal) {
							cc, err := CargoConstitucionalFromData(votoCargo.CodigoCargo.Bytes)
							checkError(err)

							fmt.Println(
								"Cargo:", cc.String(),
								", Voto:", votoVotavel.IdentificacaoVotavel.Codigo,
							)
						}
					}
				}
			}
		}

	default:
		fmt.Println(envelope, err)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
