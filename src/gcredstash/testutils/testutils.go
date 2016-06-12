package testutils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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
