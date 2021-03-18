package gcredstash

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/golang/mock/gomock"
	. "github.com/kgaughan/gcredstash/src/gcredstash"
	"github.com/kgaughan/gcredstash/src/mockaws"
	"testing"
)

func TestCreateTable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mddb := mockaws.NewMockDynamoDBAPI(ctrl)
	mkms := mockaws.NewMockKMSAPI(ctrl)
	table := "credential-store"

	mddb.EXPECT().CreateTable(&dynamodb.CreateTableInput{
		TableName: aws.String(table),
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("name"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("version"),
				KeyType:       aws.String("RANGE"),
			},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("name"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("version"),
				AttributeType: aws.String("S"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
	}).Return(nil, nil)

	driver := &Driver{
		Ddb: mddb,
		Kms: mkms,
	}

	err := driver.CreateTable(table)

	if err != nil {
		t.Errorf("\nexpected: %v\ngot: %v\n", nil, err)
	}
}

func TestWaitUntilTableExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mddb := mockaws.NewMockDynamoDBAPI(ctrl)
	mkms := mockaws.NewMockKMSAPI(ctrl)
	table := "credential-store"

	mddb.EXPECT().DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(table),
	}).Return(&dynamodb.DescribeTableOutput{
		Table: &dynamodb.TableDescription{
			TableStatus: aws.String("ACTIVE"),
		},
	}, nil)

	driver := &Driver{
		Ddb: mddb,
		Kms: mkms,
	}

	err := driver.WaitUntilTableExists(table)

	if err != nil {
		t.Errorf("\nexpected: %v\ngot: %v\n", nil, err)
	}
}

func TestIsTableExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mddb := mockaws.NewMockDynamoDBAPI(ctrl)
	mkms := mockaws.NewMockKMSAPI(ctrl)
	table := "credential-store"

	driver := &Driver{
		Ddb: mddb,
		Kms: mkms,
	}

	mddb.EXPECT().ListTablesPages(
		&dynamodb.ListTablesInput{},
		gomock.Any(),
	).Return(nil)

	isExist, err := driver.IsTableExists(table)

	if isExist {
		t.Errorf("\nexpected: %v\ngot: %v\n", false, isExist)
	}

	if err != nil {
		t.Errorf("\nexpected: %v\ngot: %v\n", nil, err)
	}
}

func TestCreateDdbTable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mddb := mockaws.NewMockDynamoDBAPI(ctrl)
	mkms := mockaws.NewMockKMSAPI(ctrl)
	table := "credential-store"

	mddb.EXPECT().ListTablesPages(
		&dynamodb.ListTablesInput{},
		gomock.Any(),
	).Return(nil)

	mddb.EXPECT().CreateTable(&dynamodb.CreateTableInput{
		TableName: aws.String(table),
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("name"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("version"),
				KeyType:       aws.String("RANGE"),
			},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("name"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("version"),
				AttributeType: aws.String("S"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
	}).Return(nil, nil)

	mddb.EXPECT().DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(table),
	}).Return(&dynamodb.DescribeTableOutput{
		Table: &dynamodb.TableDescription{
			TableStatus: aws.String("ACTIVE"),
		},
	}, nil)

	driver := &Driver{
		Ddb: mddb,
		Kms: mkms,
	}

	err := driver.CreateDdbTable(table)

	if err != nil {
		t.Errorf("\nexpected: %v\ngot: %v\n", nil, err)
	}
}
