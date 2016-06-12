package gcredstash

import (
	"encoding/hex"
)

func HexDecode(encoded string) []byte {
	decoded, err := hex.DecodeString(encoded)

	if err != nil {
		panic(err)
	}

	return decoded
}

func HexDecodeStr(encoded string) string {
	decoded := HexDecode(encoded)
	return string(decoded)
}

func HexEncode(decoded []byte) string {
	return hex.EncodeToString(decoded)
}

func HexEncodeStr(decoded string) string {
	return HexEncode([]byte(decoded))
}
