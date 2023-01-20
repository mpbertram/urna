package ue

import (
	"crypto/ed25519"
	"crypto/sha512"
	"encoding/asn1"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/exp/slices"
	"golang.org/x/text/encoding/charmap"
)

type BuEntry struct {
	Path string
}

func (entry BuEntry) ReadBu() (EntidadeBoletimUrna, error) {
	f, err := os.ReadFile(entry.Path)
	if err != nil {
		return EntidadeBoletimUrna{}, err
	}

	b, err := readBuFromBytes(f)
	if err != nil {
		return EntidadeBoletimUrna{}, err
	}

	return b, nil
}

func readBuFromBytes(bytes []byte) (EntidadeBoletimUrna, error) {
	var e EntidadeEnvelopeGenerico
	_, err := asn1.Unmarshal(bytes, &e)
	if err != nil {
		return EntidadeBoletimUrna{}, err
	}

	return e.ReadBu()
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

func CountVotos(entries []BuEntry, cargos []CargoConstitucional) map[CargoConstitucional]map[string]int {
	if len(cargos) == 0 {
		cargos = ValidCargoConstitucional()
	}

	votosPorCargo := make(map[CargoConstitucional]map[string]int)
	for _, cargo := range cargos {
		votosPorCargo[cargo] = map[string]int{}
	}

	for _, entry := range entries {
		b, err := entry.ReadBu()
		if err != nil {
			log.Println("error processing entry", entry)
		}
		for cargo, candidato := range CountVotosBu(b, cargos) {
			for candidato, numVotos := range candidato {
				votosPorCargo[cargo][candidato] = votosPorCargo[cargo][candidato] + numVotos
			}
		}
	}

	return votosPorCargo
}

func CountVotosBu(b EntidadeBoletimUrna, cargos []CargoConstitucional) map[CargoConstitucional]map[string]int {
	votosPorCargo := make(map[CargoConstitucional]map[string]int)
	for _, cargo := range cargos {
		votosPorCargo[cargo] = map[string]int{}
	}

	for _, votacaoPorEleicao := range b.ResultadosVotacaoPorEleicao {
		for _, votacao := range votacaoPorEleicao.ResultadosVotacao {
			for _, votoCargo := range votacao.TotaisVotosCargo {
				for _, votoVotavel := range votoCargo.VotosVotaveis {
					cargo, _ := CargoConstitucionalFromData(votoCargo.CodigoCargo.Bytes)

					if slices.Contains(cargos, cargo) {
						switch TipoVoto(votoVotavel.TipoVoto) {
						case Nominal, Legenda:
							candidato := fmt.Sprint(votoVotavel.IdentificacaoVotavel.Codigo)
							votosPorCargo[cargo][candidato] += votoVotavel.QuantidadeVotos
						case Branco:
							votosPorCargo[cargo][Branco.String()] += votoVotavel.QuantidadeVotos
						case Nulo:
							votosPorCargo[cargo][Nulo.String()] += votoVotavel.QuantidadeVotos
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
			log.Println("could not read", entry.Path)
			continue
		}

		err = ValidateVotosBu(b)
		if err != nil {
			log.Println(err, b)
			return err
		}
	}

	return nil
}

func ValidateVotosBu(b EntidadeBoletimUrna) error {
	pub := ed25519.PublicKey(b.ChaveAssinaturaVotosVotavel)
	for _, votacaoPorEleicao := range b.ResultadosVotacaoPorEleicao {
		for _, votacao := range votacaoPorEleicao.ResultadosVotacao {
			for _, votoCargo := range votacao.TotaisVotosCargo {
				for _, votoVotavel := range votoCargo.VotosVotaveis {

					payload := buildPayload(votoCargo, votoVotavel, b.Urna.CorrespondenciaResultado.Carga)
					checksum := sha512.Sum512(payload)

					ok := ed25519.Verify(pub, checksum[:], votoVotavel.Assinatura)
					if !ok {
						return fmt.Errorf(
							"error in verification; municipio=%s, local=%d, secao=%d, payload=%s",
							b.IdentificacaoSecao.Municipio(),
							b.IdentificacaoSecao.Local,
							b.IdentificacaoSecao.Secao,
							payload,
						)
					}

					log.Printf(
						"verified; municipio=%s, local=%d, secao=%d, payload=%s",
						b.IdentificacaoSecao.Municipio(),
						b.IdentificacaoSecao.Local,
						b.IdentificacaoSecao.Secao,
						payload,
					)
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
