package gcredstash

import (
	"testing"
)

func TestHexDecodeStr(t *testing.T) {
	expected := "London Bridge is broken down"
	actual := HexDecodeStr("4c6f6e646f6e204272696467652069732062726f6b656e20646f776e")

	if expected != actual {
		t.Errorf("\nexpected: %v\ngot: %v\n", expected, actual)
	}
}

func TestHexEncodeStr(t *testing.T) {
	expected := "4c6f6e646f6e204272696467652069732062726f6b656e20646f776e"
	actual := HexEncodeStr("London Bridge is broken down")

	if expected != actual {
		t.Errorf("\nexpected: %v\ngot: %v\n", expected, actual)
	}
}
