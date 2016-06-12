package gcredstash

import (
	"encoding/base64"
)

func B64Decode(encoded string) []byte {
	decoded, err := base64.StdEncoding.DecodeString(encoded)

	if err != nil {
		panic(err)
	}

	return decoded
}

func B64DecodeStr(encoded string) string {
	decoded := B64Decode(encoded)
	return string(decoded)
}

func B64Encode(decoded []byte) string {
	return base64.StdEncoding.EncodeToString(decoded)
}

func B64EncodeStr(decoded string) string {
	return B64Encode([]byte(decoded))
}
