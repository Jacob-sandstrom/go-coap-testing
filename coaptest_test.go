package coaptest

import (
	"testing"
)

func TestContainsOptions(t *testing.T) {
	expected := []Option{{Number: 11, Value: []byte("bs")}}
	got := []Option{{Number: 11, Delta: 11, Len: 2, DeltaExtended: []byte{}, LenExtended: []byte{}, Value: []byte{98, 115}}, {Number: 12, Delta: 1, Len: 0, DeltaExtended: []byte{}, LenExtended: []byte{}, Value: []byte{}}}

	if !containsOptions(expected, got) {
		t.Errorf("Error")
	}
}
