package urna

import (
	"crypto/ed25519"
	"crypto/sha512"
	"encoding/asn1"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/exp/slices"
	"golang.org/x/text/encoding/charmap"
)

type BuEntry struct {
	path string
}

func (b BuEntry) ReadBu() (EntidadeBoletimUrna, error) {
	f, err := os.ReadFile(b.path)
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

	var bu EntidadeBoletimUrna
	_, err = asn1.Unmarshal(e.Conteudo, &bu)
	if err != nil {
		return EntidadeBoletimUrna{}, err
	}

	return bu, nil
}

func ReadAllBu(dir string) ([]BuEntry, error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return []BuEntry{}, err
	}

	var bus []BuEntry
	for _, e := range dirEntries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".bu") {
			bu := BuEntry{strings.Join([]string{dir, e.Name()}, "/")}
			if err != nil {
				log.Println("error reading file", err)
			}
			bus = append(bus, bu)
		}
	}

	return bus, nil
}

func ComputeVotos(buEntries []BuEntry, cargos []CargoConstitucional) map[string]map[string]int {
	if len(cargos) == 0 {
		cargos = ValidCargoConstitucional()
	}

	votosPorCargo := make(map[string]map[string]int)
	for _, cargo := range cargos {
		votosPorCargo[cargo.String()] = map[string]int{}
	}

	for _, entry := range buEntries {
		b, err := entry.ReadBu()
		if err != nil {
			log.Println("could not read", entry.path)
		}

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

func ValidateVotos(buEntries []BuEntry) error {
	for _, entry := range buEntries {
		b, err := entry.ReadBu()
		if err != nil {
			log.Println("could not read", entry.path)
		}

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
