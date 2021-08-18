package coaptest

import (
	"fmt"
	"testing"
)

func TestContainsOptions(t *testing.T) {
	expected := []Option{{Number: 11, Value: []byte("bs")}}
	got := []Option{{Number: 11, Delta: 11, Len: 2, DeltaExtended: []byte{}, LenExtended: []byte{}, Value: []byte{98, 115}}, {Number: 12, Delta: 1, Len: 0, DeltaExtended: []byte{}, LenExtended: []byte{}, Value: []byte{}}}

	if !containsOptions(expected, got) {
		t.Errorf("Error")
	}
}

func TestParseCoap(t *testing.T) {
	expected := CoapMsg{Version: 1, Type: 0, TokenLen: 4, Code: 2, MsgID: []byte{111, 111}, Token: []byte{222, 222, 222, 222}, Options: []Option{{Number: 11, Delta: 11, Len: 2, Value: []byte("bs")}}, Payload: []byte("hej")}

	bytes := []byte{68, 2, 111, 111, 222, 222, 222, 222, 178, 98, 115, 255, 104, 101, 106}

	got := ParseCoap(bytes, len(bytes))

	fmt.Printf("%+v\n\n", expected)
	fmt.Printf("%+v\n\n", got)

	err := CompareCoap(expected, got)
	if err != nil {
		t.Errorf(err.Error())
	}
}
