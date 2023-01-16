package urna

import (
	"archive/zip"
	"bytes"
	"crypto/ed25519"
	"crypto/sha512"
	"encoding/asn1"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/exp/slices"
	"golang.org/x/text/encoding/charmap"
)

type BuEntry struct {
	path string
}

func (entry BuEntry) ReadBu() (EntidadeBoletimUrna, error) {
	f, err := os.ReadFile(entry.path)
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

func ProcessAllZip(dir string, process func(EntidadeBoletimUrna) error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range dirEntries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".zip") {
			ProcessZip((strings.Join([]string{dir, e.Name()}, "/")), process)
		}
	}
}

func ProcessZip(path string, process func(EntidadeBoletimUrna) error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	for _, f := range r.File {
		if strings.HasSuffix(f.Name, ".bu") {
			processZipFile(f, process)
		}
	}
}

func processZipFile(f *zip.File, process func(EntidadeBoletimUrna) error) {
	rc, err := f.Open()
	if err != nil {
		log.Println(err)
		return
	}

	var buf bytes.Buffer
	io.Copy(io.Writer(&buf), rc)

	if err != nil {
		log.Println(err)
		return
	}

	bu, err := readBuFromBytes(buf.Bytes())
	if err != nil {
		log.Println(err)
	}

	err = process(bu)
	if err != nil {
		log.Println(err)
	}

	rc.Close()
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
			log.Println("error processing entry", entry)
		}
		for cargo, candidato := range ComputeVotosBu(b, cargos) {
			for candidato, numVotos := range candidato {
				votosPorCargo[cargo][candidato] = votosPorCargo[cargo][candidato] + numVotos
			}
		}
	}

	return votosPorCargo
}

func ComputeVotosBu(b EntidadeBoletimUrna, cargos []CargoConstitucional) map[string]map[string]int {
	votosPorCargo := make(map[string]map[string]int)
	for _, cargo := range cargos {
		votosPorCargo[cargo.String()] = map[string]int{}
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

	return votosPorCargo
}

func ValidateVotos(buEntries []BuEntry) error {
	for _, entry := range buEntries {
		b, err := entry.ReadBu()
		if err != nil {
			log.Println("could not read", entry.path)
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
