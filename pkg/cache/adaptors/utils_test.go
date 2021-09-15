package adaptors

import (
	"reflect"
	"testing"
)

func Test_encode_decode(t *testing.T) {
	type Person struct {
		Name      string
		Age       uint8
		Languages []string
	}

	originalData := Person{
		Name:      "Mahdi",
		Age:       23,
		Languages: []string{"java", "python", "go", "c++"},
	}

	b, err := encode(&originalData)
	if err != nil {
		t.Error("Error on encoding data", err)
	}

	var newPerson Person
	err = decode(b, &newPerson)
	if err != nil {
		t.Error("Error on decoding data", err)
	}

	if eq := reflect.DeepEqual(newPerson, originalData); !eq {
		t.Error("Expected", originalData, "Got", newPerson)
	}
}
