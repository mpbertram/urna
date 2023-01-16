package urna

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Municipio struct {
	Id int
}

var cache map[int]string = make(map[int]string)

func (m Municipio) String() string {
	v, ok := cache[m.Id]
	if ok {
		return v
	}

	f, err := os.Open("resource/municipios.csv")
	if err != nil {
		log.Println(err)
		return ""
	}
	defer f.Close()

	r := csv.NewReader(f)
	for {
		record, err := r.Read()
		if err != nil {
			log.Println(err)
			return ""
		}

		id, err := strconv.Atoi(record[0])
		if err != nil {
			log.Println(err)
			return ""
		}

		if id == m.Id {
			s := fmt.Sprintf("%s (%s)", record[1], record[2])
			cache[m.Id] = s
			return s
		}
	}
}
