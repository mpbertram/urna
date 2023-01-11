package urna

import (
	"encoding/asn1"
	"reflect"
)

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
