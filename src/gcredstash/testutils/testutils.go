package testutils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"io/ioutil"
	"os"
)

func MapToItem(m map[string]string) map[string]*dynamodb.AttributeValue {
	item := map[string]*dynamodb.AttributeValue{}

	for key, value := range m {
		item[key] = &dynamodb.AttributeValue{S: aws.String(value)}
	}

	return item
}

func ItemToMap(item map[string]*dynamodb.AttributeValue) map[string]string {
	m := map[string]string{}

	for key, value := range item {
		m[key] = *value.S
	}

	return m
}

func TempFile(content string, f func(*os.File)) {
	tmpfile, err := ioutil.TempFile("", "gcredstash")

	if err != nil {
		panic(err)
	}

	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.WriteString(content)

	if err != nil {
		panic(err)
	}

	err = tmpfile.Sync()

	if err != nil {
		panic(err)
	}

	f(tmpfile)

	err = tmpfile.Close()

	if err != nil {
		panic(err)
	}
}
