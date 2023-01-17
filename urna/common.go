package urna

import (
	"encoding/asn1"
	"reflect"
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
