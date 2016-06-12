package command

import (
	"fmt"
	"gcredstash"
	. "gcredstash/command"
	"gcredstash/testutils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/golang/mock/gomock"
	"mockaws"
	"testing"
)

func TestListCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mddb := mockaws.NewMockDynamoDBAPI(ctrl)
	mkms := mockaws.NewMockKMSAPI(ctrl)

	table := "credential-store"
	name := "test.key"
	version := "0000000000000000002"

	item := map[string]string{
		"contents": "eBtO1lgLxIe6Yw==",
		"hmac":     "b23a3efafd4795e50ca87afd7d764f263e9ae456499a8d40eece70a63ed5da27",
		"key":      "CiDY1vsR456LEdoL3+0p+PrTCleoqi/sutbDfJZNiUSpphLLAQEBAQB42Nb7EeOeixHaC9/tKfj60wpXqKov7LrWw3yWTYlEqaYAAACiMIGfBgkqhkiG9w0BBwaggZEwgY4CAQAwgYgGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMy/Oc2pOJsR0y9nbhAgEQgFsHECqku7QZiRjLmmeGyhcsgWdWvi7Op3luJu4soi5sP0pqcsjTrBJqOXHLazgyBS9wb6deP8zpXa/41WT0ZpNY9at4gw7+XRtbz8f4Rlh8WnyFnK5RZ7i0mOlD",
		"name":     name,
		"version":  version,
	}

	mddb.EXPECT().Scan(&dynamodb.ScanInput{
		TableName:                aws.String(table),
		ProjectionExpression:     aws.String("#name,version"),
		ExpressionAttributeNames: map[string]*string{"#name": aws.String("name")},
	}).Return(&dynamodb.ScanOutput{
		Items: []map[string]*dynamodb.AttributeValue{testutils.MapToItem(item)},
	}, nil)

	cmd := &ListCommand{
		Meta: Meta{
			Table:  table,
			KmsKey: "alias/credstash",
			Driver: &gcredstash.Driver{Ddb: mddb, Kms: mkms},
		},
	}

	args := []string{}
	out, err := cmd.RunImpl(args)
	expected := fmt.Sprintf("%s -- version: %d", name, gcredstash.Atoi(version))

	if err != nil {
		t.Errorf("\nexpected: %v\ngot: %v\n", nil, err)
	}

	if expected != out {
		t.Errorf("\nexpected: %v\ngot: %v\n", expected, out)
	}
}
