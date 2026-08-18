// Package encryptutil provides client-side encryption/decryption using AWS KMS
// via the AWS Encryption SDK.
package encryptutil

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"

	mpl "github.com/aws/aws-cryptographic-material-providers-library/releases/go/mpl/awscryptographymaterialproviderssmithygenerated"
	mpltypes "github.com/aws/aws-cryptographic-material-providers-library/releases/go/mpl/awscryptographymaterialproviderssmithygeneratedtypes"
	encryptClient "github.com/aws/aws-encryption-sdk/releases/go/encryption-sdk/awscryptographyencryptionsdksmithygenerated"
	esdktypes "github.com/aws/aws-encryption-sdk/releases/go/encryption-sdk/awscryptographyencryptionsdksmithygeneratedtypes"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

// EncryptString encrypts plaintext using AWS KMS envelope encryption and returns
// the ciphertext as a base64-encoded string.
func EncryptString(ctx context.Context, kmsClient *kms.Client, kmsKeyID string, plaintext string) (string, error) {
	if kmsKeyID == "" {
		return "", errors.New("KMS key id is empty")
	}

	matProv, err := mpl.NewClient(mpltypes.MaterialProvidersConfig{})
	if err != nil {
		return "", err
	}

	keyring, err := matProv.CreateAwsKmsKeyring(ctx, mpltypes.CreateAwsKmsKeyringInput{
		KmsClient: kmsClient,
		KmsKeyId:  kmsKeyID,
	})
	if err != nil {
		return "", err
	}

	enc, err := encryptClient.NewClient(esdktypes.AwsEncryptionSdkConfig{})
	if err != nil {
		return "", err
	}

	out, err := enc.Encrypt(ctx, esdktypes.EncryptInput{
		Plaintext: []byte(plaintext),
		Keyring:   keyring,
	})
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(out.Ciphertext), nil
}

// DecryptString decrypts a base64-encoded ciphertext using AWS KMS and returns
// the original plaintext.
func DecryptString(ctx context.Context, kmsClient *kms.Client, kmsKeyID string, ciphertextB64 string) (string, error) {
	if kmsKeyID == "" {
		return "", errors.New("KMS key id is empty")
	}

	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return "", err
	}

	matProv, err := mpl.NewClient(mpltypes.MaterialProvidersConfig{})
	if err != nil {
		return "", err
	}

	keyring, err := matProv.CreateAwsKmsKeyring(ctx, mpltypes.CreateAwsKmsKeyringInput{
		KmsClient: kmsClient,
		KmsKeyId:  kmsKeyID,
	})
	if err != nil {
		return "", err
	}

	enc, err := encryptClient.NewClient(esdktypes.AwsEncryptionSdkConfig{})
	if err != nil {
		return "", err
	}

	out, err := enc.Decrypt(ctx, esdktypes.DecryptInput{
		Ciphertext: ciphertext,
		Keyring:    keyring,
	})
	if err != nil {
		return "", err
	}

	return string(out.Plaintext), nil
}

// DecodeJWT decodes a JWT token's payload (without verification) and returns the claims.
func DecodeJWT(token string) (map[string]any, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid JWT format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, err
	}

	return claims, nil
}

// DecodeBase64 decodes a standard base64-encoded string.
func DecodeBase64(input string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(input)
}

// DecodeBase64URL decodes a URL-safe base64-encoded string (no padding).
func DecodeBase64URL(input string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(input)
}
