package gcredstash

import (
	"bytes"
	. "github.com/winebarrel/gcredstash/src/gcredstash"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/golang/mock/gomock"
	"mockaws"
	"testing"
)

func TestKmsDecrypt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	blob := []byte("123")
	context := map[string]string{"foo": "bar"}
	expectedDataKey := []byte("12345678901234567890123456789012")
	expectedHmacKey := []byte("abc")

	mkms := mockaws.NewMockKMSAPI(ctrl)

	mkms.EXPECT().Decrypt(&kms.DecryptInput{
		CiphertextBlob:    blob,
		EncryptionContext: map[string]*string{"foo": aws.String("bar")},
	}).Return(&kms.DecryptOutput{
		Plaintext: append(expectedDataKey, expectedHmacKey...),
	}, nil)

	dataKey, hmacKey, err := KmsDecrypt(mkms, blob, context)

	if err != nil {
		t.Errorf("\nexpected: %v\ngot: %v\n", nil, err)
	}

	if !bytes.Equal(expectedDataKey, dataKey) {
		t.Errorf("\nexpected: %v\ngot: %v\n", expectedDataKey, dataKey)
	}

	if !bytes.Equal(expectedHmacKey, hmacKey) {
		t.Errorf("\nexpected: %v\ngot: %v\n", expectedHmacKey, hmacKey)
	}
}

func TestKmsGenerateDataKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	keyId := "alias/credstash"
	context := map[string]string{"foo": "bar"}

	expectedDataKey := []byte("12345678901234567890123456789012")
	expectedHmacKey := []byte("abc")
	expectedWrappedKey := []byte("blobData")

	mkms := mockaws.NewMockKMSAPI(ctrl)

	mkms.EXPECT().GenerateDataKey(&kms.GenerateDataKeyInput{
		KeyId:             aws.String(keyId),
		NumberOfBytes:     aws.Int64(64),
		EncryptionContext: map[string]*string{"foo": aws.String("bar")},
	}).Return(&kms.GenerateDataKeyOutput{
		Plaintext:      append(expectedDataKey, expectedHmacKey...),
		CiphertextBlob: expectedWrappedKey,
	}, nil)

	dataKey, hmacKey, wrappedKey, err := KmsGenerateDataKey(mkms, keyId, context)

	if err != nil {
		t.Errorf("\nexpected: %v\ngot: %v\n", nil, err)
	}

	if !bytes.Equal(expectedDataKey, dataKey) {
		t.Errorf("\nexpected: %v\ngot: %v\n", expectedDataKey, dataKey)
	}

	if !bytes.Equal(expectedHmacKey, hmacKey) {
		t.Errorf("\nexpected: %v\ngot: %v\n", expectedHmacKey, hmacKey)
	}

	if !bytes.Equal(expectedWrappedKey, wrappedKey) {
		t.Errorf("\nexpected: %v\ngot: %v\n", expectedWrappedKey, wrappedKey)
	}
}
