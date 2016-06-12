package gcredstash

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"time"
)

func (driver *Driver) IsTableExists(table string) (bool, error) {
	params := &dynamodb.ListTablesInput{}
	isExist := false

	err := driver.Ddb.ListTablesPages(params, func(page *dynamodb.ListTablesOutput, lastPage bool) bool {
		for _, tableName := range page.TableNames {
			if *tableName == table {
				isExist = true
				return false
			}
		}

		return true
	})

	if err != nil {
		return false, err
	}

	return isExist, nil
}

func (driver *Driver) CreateTable(table string) error {
	params := &dynamodb.CreateTableInput{
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
	}

	_, err := driver.Ddb.CreateTable(params)

	return err
}

func (driver *Driver) WaitUntilTableExists(table string) error {
	delay := 20 * time.Second
	maxAttempts := 25
	isCreated := false

	params := &dynamodb.DescribeTableInput{
		TableName: aws.String(table),
	}

	for i := 0; i < 25; i++ {
		resp, err := driver.Ddb.DescribeTable(params)

		if err != nil {
			return err
		}

		if *resp.Table.TableStatus == "ACTIVE" {
			isCreated = true
			break
		}

		time.Sleep(delay)
	}

	if !isCreated {
		return fmt.Errorf("exceeded %d wait attempts", maxAttempts)
	}

	return nil
}

func (driver *Driver) CreateDdbTable(table string) error {
	tableIsExist, err := driver.IsTableExists(table)

	if err != nil {
		return err
	}

	if tableIsExist {
		return fmt.Errorf("Credential Store table already exists: %s", table)
	}

	err = driver.CreateTable(table)

	if err != nil {
		return err
	}

	fmt.Println("Creating table...")
	fmt.Println("Waiting for table to be created...")

	err = driver.WaitUntilTableExists(table)

	if err != nil {
		return err
	}

	fmt.Println("Table has been created. Go read the README about how to create your KMS key")

	return nil
}
