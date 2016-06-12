package gcredstash

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
)

func Digest(message []byte, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	return mac.Sum(nil)
}

func ValidateHMAC(message []byte, digest []byte, key []byte) bool {
	expected := Digest(message, key)
	return hmac.Equal(digest, expected)
}

func Crypt(contents []byte, key []byte) []byte {
	block, err := aes.NewCipher(key)

	if err != nil {
		panic(err)
	}

	text := make([]byte, len(contents))

	iv := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(text, contents)

	return text
}
