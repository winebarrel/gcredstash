package gcredstash

import (
	"testing"
)

func TestB64DecodeStr(t *testing.T) {
	expected := "London Bridge is broken down"
	actual := B64DecodeStr("TG9uZG9uIEJyaWRnZSBpcyBicm9rZW4gZG93bg==")

	if expected != actual {
		t.Errorf("\nexpected: %v\ngot: %v\n", expected, actual)
	}
}

func TestB64EncodeStr(t *testing.T) {
	expected := "TG9uZG9uIEJyaWRnZSBpcyBicm9rZW4gZG93bg=="
	actual := B64EncodeStr("London Bridge is broken down")

	if expected != actual {
		t.Errorf("\nexpected: %v\ngot: %v\n", expected, actual)
	}
}
