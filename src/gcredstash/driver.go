package gcredstash

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
	"strings"
)

type Driver struct {
	Ddb dynamodbiface.DynamoDBAPI
	Kms kmsiface.KMSAPI
}

func (driver *Driver) GetMaterialWithoutVersion(name string, table string) (map[string]*dynamodb.AttributeValue, error) {
	params := &dynamodb.QueryInput{
		TableName:                aws.String(table),
		Limit:                    aws.Int64(1),
		ConsistentRead:           aws.Bool(true),
		ScanIndexForward:         aws.Bool(false),
		KeyConditionExpression:   aws.String("#name = :name"),
		ExpressionAttributeNames: map[string]*string{"#name": aws.String("name")},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":name": {S: aws.String(name)},
		},
	}

	resp, err := driver.Ddb.Query(params)

	if err != nil {
		return nil, err
	}

	if *resp.Count == 0 {
		return nil, fmt.Errorf("Item {'name': '%s'} couldn't be found.", name)
	}

	return resp.Items[0], nil
}

func (driver *Driver) GetMaterialWithVersion(name string, version string, table string) (map[string]*dynamodb.AttributeValue, error) {
	params := &dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key: map[string]*dynamodb.AttributeValue{
			"name":    {S: aws.String(name)},
			"version": {S: aws.String(version)},
		},
	}

	resp, err := driver.Ddb.GetItem(params)

	if err != nil {
		return nil, err
	}

	if resp.Item == nil {
		return nil, fmt.Errorf("Item {'name': '%s'} couldn't be found.", name)
	}

	return resp.Item, nil
}

func (driver *Driver) DecryptMaterial(name string, material map[string]*dynamodb.AttributeValue, context map[string]string) (string, error) {
	data := B64Decode(*material["key"].S)
	dataKey, hmacKey, err := KmsDecrypt(driver.Kms, data, context)

	if err != nil {
		if strings.Contains(err.Error(), "InvalidCiphertextException") {
			if len(context) < 1 {
				return "", fmt.Errorf("%s: Could not decrypt hmac key with KMS. The credential may require that an encryption context be provided to decrypt it.", name)
			} else {
				return "", fmt.Errorf("%s: Could not decrypt hmac key with KMS. The encryption context provided may not match the one used when the credential was stored.", name)
			}
		} else {
			return "", err
		}
	}

	contents := B64Decode(*material["contents"].S)
	var hmac []byte
	if (*material["hmac"]).B != nil {
		hmac_hex := (*material["hmac"]).B
		hmac = HexDecode(string(hmac_hex))
	} else if (*material["hmac"]).S != nil {
		hmac = HexDecode(*(*material["hmac"]).S)
	}

	if !ValidateHMAC(contents, hmac, hmacKey) {
		return "", fmt.Errorf("Computed HMAC on %s does not match stored HMAC", name)
	}

	decrypted := Crypt(contents, dataKey)

	return string(decrypted), nil
}

func (driver *Driver) GetHighestVersion(name string, table string) (int, error) {
	params := &dynamodb.QueryInput{
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
	}

	resp, err := driver.Ddb.Query(params)

	if err != nil {
		return -1, err
	}

	if *resp.Count == 0 {
		return 0, nil
	}

	version := *resp.Items[0]["version"].S
	versionNum := Atoi(version)

	return versionNum, nil
}

func (driver *Driver) PutItem(name string, version string, key []byte, contents []byte, hmac []byte, table string) error {
	b64key := B64Encode(key)
	b64contents := B64Encode(contents)
	hexHmac := HexEncode(hmac)

	params := &dynamodb.PutItemInput{
		TableName: aws.String(table),
		Item: map[string]*dynamodb.AttributeValue{
			"name":     {S: aws.String(name)},
			"version":  {S: aws.String(version)},
			"key":      {S: aws.String(b64key)},
			"contents": {S: aws.String(b64contents)},
			"hmac":     {S: aws.String(hexHmac)},
		},
		ConditionExpression:      aws.String("attribute_not_exists(#name)"),
		ExpressionAttributeNames: map[string]*string{"#name": aws.String("name")},
	}

	_, err := driver.Ddb.PutItem(params)

	if err != nil {
		return err
	}

	return nil
}

func (driver *Driver) GetDeleteTargetWithoutVersion(name string, table string) (map[*string]*string, error) {
	items := map[*string]*string{}

	params := &dynamodb.QueryInput{
		TableName:                aws.String(table),
		ConsistentRead:           aws.Bool(true),
		KeyConditionExpression:   aws.String("#name = :name"),
		ExpressionAttributeNames: map[string]*string{"#name": aws.String("name")},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":name": {S: aws.String(name)},
		},
	}

	resp, err := driver.Ddb.Query(params)

	if err != nil {
		return nil, err
	}

	if *resp.Count == 0 {
		return nil, fmt.Errorf("Item {'name': '%s'} couldn't be found.", name)
	}

	for _, i := range resp.Items {
		items[i["name"].S] = i["version"].S
	}

	return items, nil
}

func (driver *Driver) GetDeleteTargetWithVersion(name string, version string, table string) (map[*string]*string, error) {
	params := &dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key: map[string]*dynamodb.AttributeValue{
			"name":    {S: aws.String(name)},
			"version": {S: aws.String(version)},
		},
	}

	resp, err := driver.Ddb.GetItem(params)

	if err != nil {
		return nil, err
	}

	if resp.Item == nil {
		versionNum := Atoi(version)
		return nil, fmt.Errorf("Item {'name': '%s', 'version': %d} couldn't be found.", name, versionNum)
	}

	items := map[*string]*string{}
	items[resp.Item["name"].S] = resp.Item["version"].S

	return items, nil
}

func (driver *Driver) DeleteItem(name string, version string, table string) error {
	svc := driver.Ddb

	params := &dynamodb.DeleteItemInput{
		TableName: aws.String(table),
		Key: map[string]*dynamodb.AttributeValue{
			"name":    {S: aws.String(name)},
			"version": {S: aws.String(version)},
		},
	}

	_, err := svc.DeleteItem(params)

	if err != nil {
		return err
	}

	return nil
}

func (driver *Driver) DeleteSecrets(name string, version string, table string) error {
	var items map[*string]*string
	var err error

	if version == "" {
		items, err = driver.GetDeleteTargetWithoutVersion(name, table)
	} else {
		items, err = driver.GetDeleteTargetWithVersion(name, version, table)
	}

	if err != nil {
		return err
	}

	for name, version := range items {
		err := driver.DeleteItem(*name, *version, table)

		if err != nil {
			return err
		}

		versionNum := Atoi(*version)
		fmt.Printf("Deleting %s -- version %d\n", *name, versionNum)
	}

	return nil
}

func (driver *Driver) PutSecret(name string, secret string, version string, kmsKey string, table string, context map[string]string) error {
	dataKey, hmacKey, wrappedKey, err := KmsGenerateDataKey(driver.Kms, kmsKey, context)

	if err != nil {
		return fmt.Errorf("Could not generate key using KMS key(%s): %s", kmsKey, err.Error())
	}

	cipherText := Crypt([]byte(secret), dataKey)
	hmac := Digest(cipherText, hmacKey)

	err = driver.PutItem(name, version, wrappedKey, cipherText, hmac, table)

	if err != nil {
		if strings.Contains(err.Error(), "ConditionalCheckFailedException") {
			latestVersion, err := driver.GetHighestVersion(name, table)

			if err != nil {
				return err
			}

			return fmt.Errorf(
				"%s version %d is already in the credential store. Use the -v flag to specify a new version",
				name,
				latestVersion)
		} else {
			return err
		}
	}

	return nil
}

func (driver *Driver) GetSecret(name string, version string, table string, context map[string]string) (string, error) {
	var material map[string]*dynamodb.AttributeValue
	var err error

	if version == "" {
		material, err = driver.GetMaterialWithoutVersion(name, table)
	} else {
		material, err = driver.GetMaterialWithVersion(name, version, table)
	}

	if err != nil {
		return "", err
	}

	value, err := driver.DecryptMaterial(name, material, context)

	if err != nil {
		return "", err
	}

	return value, nil
}

func (driver *Driver) ListSecrets(table string) (map[*string]*string, error) {
	svc := driver.Ddb

	params := &dynamodb.ScanInput{
		TableName:                aws.String(table),
		ProjectionExpression:     aws.String("#name,version"),
		ExpressionAttributeNames: map[string]*string{"#name": aws.String("name")},
	}

	resp, err := svc.Scan(params)

	if err != nil {
		return nil, err
	}

	items := map[*string]*string{}

	for _, i := range resp.Items {
		items[i["name"].S] = i["version"].S
	}

	return items, nil
}
