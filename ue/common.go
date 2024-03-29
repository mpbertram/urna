package ue

import (
	"archive/zip"
	"bytes"
	"github.com/google/certificate-transparency-go/asn1"
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

var zipCache map[string]*zip.ReadCloser = make(map[string]*zip.ReadCloser)

type ZipProcessCtx struct {
	ZipFilename string // name of the `*.zip` file
	Filename    string // name of the file inside the `*.zip` file
}

func ProcessAllZipRaw(dir string, process func(*zip.File) bool) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range dirEntries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".zip") {
			ProcessZipRaw((strings.Join([]string{dir, e.Name()}, "/")), process)
		}
	}
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

func ProcessZipRaw(path string, process func(*zip.File) bool) {
	r, err := openZipReader(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range r.File {
		done := process(f)
		if done {
			break
		}
	}
}

func ProcessZip(path string, process any) {
	r, err := openZipReader(path)
	if err != nil {
		log.Fatal(err)
	}

	entityType := reflect.TypeOf(process).In(0)
	entity := reflect.New(entityType)

	extensionMethod := entity.MethodByName("Extension")
	if !extensionMethod.IsValid() {
		log.Fatal("type of process function argument does not define Extension()")
	}
	extension := extensionMethod.Call([]reflect.Value{})[0]

	for _, f := range r.File {
		if strings.HasSuffix(f.Name, extension.String()) {
			rc, err := f.Open()
			if err != nil {
				log.Println("could not open file inside zip:", err)
				return
			}

			var buf bytes.Buffer
			io.Copy(&buf, rc)

			if err != nil {
				log.Println(err)
				return
			}

			functionType := reflect.TypeOf(process)
			entityType := functionType.In(0)
			entity := reflect.New(entityType)
			_, err = asn1.Unmarshal(buf.Bytes(), entity.Interface())
			if err != nil {
				log.Println(err)
			}

			if functionType.NumIn() > 1 {
				ctx := functionType.In(1)
				if ctx == reflect.TypeOf(ZipProcessCtx{}) {
					reflect.ValueOf(process).Call(
						[]reflect.Value{
							entity.Elem(),
							reflect.ValueOf(ZipProcessCtx{path, f.Name}),
						},
					)
				}
			} else {
				reflect.ValueOf(process).Call([]reflect.Value{entity.Elem()})
			}

			if err != nil {
				log.Println(err)
			}

			rc.Close()
		}
	}
}

func openZipReader(path string) (*zip.ReadCloser, error) {
	r, ok := zipCache[path]
	if !ok {
		var err error
		r, err = zip.OpenReader(path)
		if err != nil {
			log.Println("could not open:", err)
			return nil, err
		}

		zipCache[path] = r

		if len(zipCache) > 10 {
			for f, r := range zipCache {
				if f == path {
					continue
				}

				err := r.Close()
				if err != nil {
					log.Println("could not close:", err)
				}
				delete(zipCache, f)

				break
			}
		}
	}

	return r, nil
}
