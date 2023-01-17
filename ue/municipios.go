package ue

import (
	"embed"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
)

//go:embed resource/municipios.csv
var f embed.FS

type Municipio struct {
	Id   int
	Nome string
	Uf   string
}

var cache map[int]Municipio = make(map[int]Municipio)

func (m Municipio) String() string {
	return fmt.Sprintf("%s (%s)", m.Nome, m.Uf)
}

func MunicipioFromId(id int) (Municipio, error) {
	m, ok := cache[id]
	if ok {
		return m, nil
	}

	municipios, err := f.Open("resource/municipios.csv")
	if err != nil {
		log.Println(err)
		return Municipio{id, "?", "?"}, errors.New("could not process csv")
	}
	defer municipios.Close()

	r := csv.NewReader(municipios)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			return Municipio{id, "?", "?"}, errors.New("could not process csv")
		}

		i, err := strconv.Atoi(record[0])
		if err != nil {
			log.Println(err)
			return Municipio{id, "?", "?"}, errors.New("could not process csv")
		}

		if i == id {
			mun := Municipio{id, record[1], record[2]}
			cache[id] = mun
			return mun, nil
		}
	}

	return Municipio{id, "?", "?"}, fmt.Errorf("could not find for id=%d", id)
}
