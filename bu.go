package urna

import (
	"crypto/ed25519"
	"crypto/sha512"
	"encoding/asn1"
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/exp/slices"
	"golang.org/x/text/encoding/charmap"
)

func ReadAllBu(dir string) ([]EntidadeBoletimUrna, error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return []EntidadeBoletimUrna{}, err
	}

	var bus []EntidadeBoletimUrna
	for _, e := range dirEntries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".bu") {
			bu, err := ReadBu(strings.Join([]string{dir, e.Name()}, "/"))
			if err != nil {
				fmt.Println("error reading file", err)
			}
			bus = append(bus, bu)
		}
	}

	return bus, nil
}

func ReadBu(file string) (EntidadeBoletimUrna, error) {
	f, err := os.ReadFile(file)
	if err != nil {
		return EntidadeBoletimUrna{}, err
	}

	var e EntidadeEnvelopeGenerico
	_, err = asn1.Unmarshal(f, &e)
	if err != nil {
		return EntidadeBoletimUrna{}, err
	}

	if TipoEnvelope(e.TipoEnvelope) != EnvelopeBoletimUrna {
		return EntidadeBoletimUrna{}, errors.New("envelope is not a bu")
	}

	var b EntidadeBoletimUrna
	_, err = asn1.Unmarshal(e.Conteudo, &b)
	if err != nil {
		return EntidadeBoletimUrna{}, err
	}

	return b, nil
}

func ComputeVotos(boletins []EntidadeBoletimUrna, cargos []CargoConstitucional) map[string]map[string]int {
	if len(cargos) == 0 {
		cargos = ValidCargoConstitucional()
	}

	votosPorCargo := make(map[string]map[string]int)
	for _, cargo := range cargos {
		votosPorCargo[cargo.String()] = map[string]int{}
	}

	for _, b := range boletins {
		for _, votacaoPorEleicao := range b.ResultadosVotacaoPorEleicao {
			for _, votacao := range votacaoPorEleicao.ResultadosVotacao {
				for _, votoCargo := range votacao.TotaisVotosCargo {
					for _, votoVotavel := range votoCargo.VotosVotaveis {
						cc, _ := CargoConstitucionalFromData(votoCargo.CodigoCargo.Bytes)

						if slices.Contains(cargos, cc) {
							switch TipoVoto(votoVotavel.TipoVoto) {
							case Nominal, Legenda:
								votosPorCargo[cc.String()][fmt.Sprint(votoVotavel.IdentificacaoVotavel.Codigo)] += votoVotavel.QuantidadeVotos
							case Branco:
								votosPorCargo[cc.String()][Branco.String()] += votoVotavel.QuantidadeVotos
							case Nulo:
								votosPorCargo[cc.String()][Nulo.String()] += votoVotavel.QuantidadeVotos
							}
						}
					}
				}
			}
		}
	}

	return votosPorCargo
}

func ValidateVotos(boletins []EntidadeBoletimUrna) error {
	for _, b := range boletins {
		pub := ed25519.PublicKey(b.ChaveAssinaturaVotosVotavel)
		for _, votacaoPorEleicao := range b.ResultadosVotacaoPorEleicao {
			for _, votacao := range votacaoPorEleicao.ResultadosVotacao {
				for _, votoCargo := range votacao.TotaisVotosCargo {
					for _, votoVotavel := range votoCargo.VotosVotaveis {

						checksum := sha512.Sum512(
							buildPayload(votoCargo, votoVotavel, b.Urna.CorrespondenciaResultado.Carga))
						ok := ed25519.Verify(pub, checksum[:], votoVotavel.Assinatura)
						if !ok {
							return errors.New("error in verification")
						}
					}
				}
			}
		}
	}

	return nil
}

func buildPayload(vc TotalVotosCargo, vv TotalVotosVotavel, c Carga) []byte {
	cc, _ := CargoConstitucionalFromData(vc.CodigoCargo.Bytes)
	codigoCargo := fmt.Sprint(int(cc))

	tipoVoto := fmt.Sprint(vv.TipoVoto)
	qtdVotos := fmt.Sprint(vv.QuantidadeVotos)

	codigo := fmt.Sprint(vv.IdentificacaoVotavel.Codigo)
	if codigo == "0" {
		codigo = ""
	}

	partido := fmt.Sprint(vv.IdentificacaoVotavel.Partido)
	if partido == "0" {
		partido = ""
	}

	carga := c.CodigoCarga

	payload := codigoCargo + tipoVoto + qtdVotos + codigo + partido + carga
	payloadBytes := []byte{}
	for _, r := range payload {
		encodedRune, _ := charmap.ISO8859_1.EncodeRune(r)
		payloadBytes = append(payloadBytes, encodedRune)
	}

	return payloadBytes
}
