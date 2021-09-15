package adaptors

import (
	"reflect"

	"github.com/vmihailenco/msgpack/v5"
)

func encode(value interface{}) ([]byte, error) {
	return msgpack.Marshal(value)
}

func decode(rawData []byte, reference interface{}) error {
	return msgpack.Unmarshal(rawData, &reference)
}

func DeepCopy(source interface{}, destination interface{}) {
	if source == nil || destination == nil {
		return
	}
	x := reflect.ValueOf(source)
	if x.Kind() == reflect.Ptr {
		starX := x.Elem()
		y := reflect.New(starX.Type())
		starY := y.Elem()
		starY.Set(starX)
		reflect.ValueOf(destination).Elem().Set(y.Elem())
	} else {
		reflect.ValueOf(destination).Elem().Set(x)
	}
}
