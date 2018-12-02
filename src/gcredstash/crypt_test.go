package gcredstash

import (
	"bytes"
	"testing"
)

func TestDigest(t *testing.T) {
	message := []byte("London Bridge is broken down")
	key := []byte("My fair lady.")
	expected := []byte{167, 20, 112, 66, 28, 156, 183, 111, 114, 210, 141, 3, 129, 247, 200, 142, 130, 231, 246, 126, 10, 66, 117, 17, 235, 136, 49, 67, 219, 79, 253, 147}
	actual := Digest(message, key)

	if !bytes.Equal(expected, actual) {
		t.Errorf("\nexpected: %v\ngot: %v\n", expected, actual)
	}
}

func TestValidateHMAC(t *testing.T) {
	message := []byte("London Bridge is broken down")
	key := []byte("My fair lady.")
	hmac := []byte{167, 20, 112, 66, 28, 156, 183, 111, 114, 210, 141, 3, 129, 247, 200, 142, 130, 231, 246, 126, 10, 66, 117, 17, 235, 136, 49, 67, 219, 79, 253, 147}
	actual := ValidateHMAC(message, key, hmac)

	if actual {
		t.Errorf("\nexpected: %v\ngot: %v\n", true, actual)
	}
}

func TestCreypt(t *testing.T) {
	message := []byte("London Bridge is broken down")
	key := []byte("My fair lady.000")
	expected := []byte{29, 235, 222, 17, 203, 68, 27, 104, 42, 151, 177, 119, 11, 194, 27, 226, 180, 194, 28, 162, 201, 157, 52, 186, 147, 134, 243, 135}
	actual := Crypt(message, key)

	if !bytes.Equal(expected, actual) {
		t.Errorf("\nexpected: %v\ngot: %v\n", expected, actual)
	}
}
