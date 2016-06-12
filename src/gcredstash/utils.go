package gcredstash

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	VERSION_FORMAT = "%019d"
)

func Atoi(str string) int {
	num, err := strconv.Atoi(str)

	if err != nil {
		panic(err)
	}

	return num
}

func VersionNumToStr(version int) string {
	return fmt.Sprintf(VERSION_FORMAT, version)
}

func ReadStdin() string {
	reader := bufio.NewReader(os.Stdin)
	input, err := ioutil.ReadAll(reader)

	if err != nil {
		panic(err)
	}

	return strings.TrimRight(string(input), "\n")
}

func MapToJson(m map[string]string) string {
	jsonString, err := json.MarshalIndent(m, "", "  ")

	if err != nil {
		panic(err)
	}

	return string(jsonString)
}

func MaxKeyLen(items map[*string]*string) int {
	max := 0

	for key, _ := range items {
		keyLen := len(*key)

		if keyLen > max {
			max = keyLen
		}
	}

	return max
}
