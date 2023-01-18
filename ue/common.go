package ue

import (
	"archive/zip"
	"bytes"
	"encoding/asn1"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
)

func FillSlice(bytes []byte, form any) error {
	var i int
	rest := bytes
	for {
		val := reflect.New(reflect.TypeOf(form).Elem().Elem())
		var err error
		rest, err = asn1.Unmarshal(rest, val.Interface())
		if err != nil {
			return err
		}

		fv := reflect.ValueOf(form).Elem()
		fv.Set(reflect.Append(fv, val.Elem()))

		if len(rest) == 0 {
			break
		}

		i++
	}

	return nil
}

func FillSequence(bytes []byte, form any) error {
	var i int
	rest := bytes
	for {
		f := reflect.ValueOf(form).Elem().Field(i)
		val := reflect.New(f.Type())
		var err error
		rest, err = asn1.Unmarshal(rest, val.Interface())
		if err != nil {
			return err
		}

		f.Set(val.Elem())

		if len(rest) == 0 {
			break
		}

		i++
	}

	return nil
}

func ProcessAllZip(dir string, process any) {
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

func ProcessZip(path string, process any) {
	r, err := zip.OpenReader(path)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	entityType := reflect.TypeOf(process).In(0)
	entity := reflect.New(entityType)
	extension := entity.MethodByName("Extension").Call([]reflect.Value{})[0]

	for _, f := range r.File {
		if strings.HasSuffix(f.Name, extension.String()) {
			processZipFile(f, process)
		}
	}
}

func processZipFile(f *zip.File, process any) {
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

	entityType := reflect.TypeOf(process).In(0)
	entity := reflect.New(entityType)
	_, err = asn1.Unmarshal(buf.Bytes(), entity.Interface())
	if err != nil {
		log.Println(err)
	}

	reflect.ValueOf(process).Call([]reflect.Value{entity.Elem()})

	if err != nil {
		log.Println(err)
	}

	rc.Close()
}
