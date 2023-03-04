package ue

import (
	"fmt"
	"strconv"
)

func MunicipioByFile(filename string) string {
	id, err := strconv.Atoi(filename[7:12])
	if err != nil {
		fmt.Println(err)
		return filename
	}

	m, err := MunicipioFromId(id)
	if err != nil {
		return filename
	}

	return m.String()
}

func ZonaByFile(filename string) string {
	id, err := strconv.Atoi(filename[12:16])
	if err != nil {
		return filename
	}

	return strconv.Itoa(id)
}

func SecaoByFile(filename string) string {
	id, err := strconv.Atoi(filename[16:20])
	if err != nil {
		return filename
	}

	return strconv.Itoa(id)
}
