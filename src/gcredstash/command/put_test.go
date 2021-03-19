package command

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/golang/mock/gomock"
	"github.com/kgaughan/gcredstash/src/gcredstash"
	. "github.com/kgaughan/gcredstash/src/gcredstash/command"
	"github.com/kgaughan/gcredstash/src/gcredstash/testutils"
	"github.com/kgaughan/gcredstash/src/mockaws"
	"testing"
)

func TestPutCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mddb := mockaws.NewMockDynamoDBAPI(ctrl)
	mkms := mockaws.NewMockKMSAPI(ctrl)

	table := "credential-store"
	secret := "100"
	name := "test.key"
	version := "0000000000000000002"
	newVersion := "0000000000000000003"
	kmsKey := "alias/credstash"

	item := map[string]string{
		"contents": "twnH",
		"hmac":     "01cc6772cf2c889c8c0dae1f0ec3d7659e21103d56cd3436039cf29d18759958",
		"key":      "CiDY1vsR456LEdoL3+0p+PrTCleoqi/sutbDfJZNiUSpphLLAQEBAQB42Nb7EeOeixHaC9/tKfj60wpXqKov7LrWw3yWTYlEqaYAAACiMIGfBgkqhkiG9w0BBwaggZEwgY4CAQAwgYgGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMeq7h5wZtkuXM8PpxAgEQgFusrxgmwCbvRObKTdbH2yvma5kNrgx3bF3ghmu7pjq6ZhPao8gZJAG2YdwwTvdbjr/wck++u0W8utaP6r07Pe8M8+oUGwWxit9X6UzxfOR6Q4eoW8g2hRUncOgF",
		"name":     name,
		"version":  version,
	}

	mddb.EXPECT().Query(&dynamodb.QueryInput{
		TableName:                aws.String(table),
		Limit:                    aws.Int64(1),
		ConsistentRead:           aws.Bool(true),
		ScanIndexForward:         aws.Bool(false),
		KeyConditionExpression:   aws.String("#name = :name"),
		ExpressionAttributeNames: map[string]*string{"#name": aws.String("name")},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":name": {S: aws.String(name)},
		},
		ProjectionExpression: aws.String("version"),
	}).Return(&dynamodb.QueryOutput{
		Count: aws.Int64(1),
		Items: []map[string]*dynamodb.AttributeValue{testutils.MapToItem(item)},
	}, nil)

	item["version"] = newVersion

	mkms.EXPECT().GenerateDataKey(&kms.GenerateDataKeyInput{
		KeyId:         aws.String(kmsKey),
		NumberOfBytes: aws.Int64(64),
	}).Return(&kms.GenerateDataKeyOutput{
		CiphertextBlob: []byte{10, 32, 216, 214, 251, 17, 227, 158, 139, 17, 218, 11, 223, 237, 41, 248, 250, 211, 10, 87, 168, 170, 47, 236, 186, 214, 195, 124, 150, 77, 137, 68, 169, 166, 18, 203, 1, 1, 1, 1, 0, 120, 216, 214, 251, 17, 227, 158, 139, 17, 218, 11, 223, 237, 41, 248, 250, 211, 10, 87, 168, 170, 47, 236, 186, 214, 195, 124, 150, 77, 137, 68, 169, 166, 0, 0, 0, 162, 48, 129, 159, 6, 9, 42, 134, 72, 134, 247, 13, 1, 7, 6, 160, 129, 145, 48, 129, 142, 2, 1, 0, 48, 129, 136, 6, 9, 42, 134, 72, 134, 247, 13, 1, 7, 1, 48, 30, 6, 9, 96, 134, 72, 1, 101, 3, 4, 1, 46, 48, 17, 4, 12, 122, 174, 225, 231, 6, 109, 146, 229, 204, 240, 250, 113, 2, 1, 16, 128, 91, 172, 175, 24, 38, 192, 38, 239, 68, 230, 202, 77, 214, 199, 219, 43, 230, 107, 153, 13, 174, 12, 119, 108, 93, 224, 134, 107, 187, 166, 58, 186, 102, 19, 218, 163, 200, 25, 36, 1, 182, 97, 220, 48, 78, 247, 91, 142, 191, 240, 114, 79, 190, 187, 69, 188, 186, 214, 143, 234, 189, 59, 61, 239, 12, 243, 234, 20, 27, 5, 177, 138, 223, 87, 233, 76, 241, 124, 228, 122, 67, 135, 168, 91, 200, 54, 133, 21, 39, 112, 232, 5},
		Plaintext:      []byte{145, 99, 240, 141, 84, 162, 135, 185, 20, 181, 81, 249, 15, 215, 56, 150, 222, 94, 65, 27, 27, 196, 165, 220, 49, 90, 199, 244, 14, 165, 188, 116, 135, 60, 104, 13, 136, 145, 109, 232, 87, 153, 237, 234, 174, 87, 7, 124, 131, 121, 67, 68, 239, 184, 174, 16, 197, 129, 97, 139, 146, 144, 89, 5},
	}, nil)

	mddb.EXPECT().PutItem(&dynamodb.PutItemInput{
		TableName:                aws.String(table),
		Item:                     testutils.MapToItem(item),
		ConditionExpression:      aws.String("attribute_not_exists(#name)"),
		ExpressionAttributeNames: map[string]*string{"#name": aws.String("name")},
	}).Return(nil, nil)

	cmd := &PutCommand{
		Meta: Meta{
			Table:  table,
			KmsKey: kmsKey,
			Driver: &gcredstash.Driver{Ddb: mddb, Kms: mkms},
		},
	}

	args := []string{name, secret, "-a"}
	err := cmd.RunImpl(args)

	if err != nil {
		t.Errorf("\nexpected: %v\ngot: %v\n", nil, err)
	}
}
