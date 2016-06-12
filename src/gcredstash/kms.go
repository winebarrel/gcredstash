package gcredstash

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
)

func KmsDecrypt(svc kmsiface.KMSAPI, blob []byte, context map[string]string) ([]byte, []byte, error) {
	params := &kms.DecryptInput{
		CiphertextBlob: blob,
	}

	if len(context) > 0 {
		ctx := map[string]*string{}

		for key, value := range context {
			ctx[key] = aws.String(value)
		}

		params.EncryptionContext = ctx
	}

	resp, err := svc.Decrypt(params)

	if err != nil {
		return nil, nil, err
	}

	dataKey := resp.Plaintext[:32]
	hmacKey := resp.Plaintext[32:]

	return dataKey, hmacKey, nil
}

func KmsGenerateDataKey(svc kmsiface.KMSAPI, keyId string, context map[string]string) ([]byte, []byte, []byte, error) {
	params := &kms.GenerateDataKeyInput{
		KeyId:         aws.String(keyId),
		NumberOfBytes: aws.Int64(64),
	}

	if len(context) > 0 {
		ctx := map[string]*string{}

		for key, value := range context {
			ctx[key] = aws.String(value)
		}

		params.EncryptionContext = ctx
	}

	resp, err := svc.GenerateDataKey(params)

	if err != nil {
		return nil, nil, nil, err
	}

	dataKey := resp.Plaintext[:32]
	hmacKey := resp.Plaintext[32:]
	wrappedKey := resp.CiphertextBlob

	return dataKey, hmacKey, wrappedKey, nil
}
