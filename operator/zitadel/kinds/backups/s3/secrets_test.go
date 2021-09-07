package s3

import (
	"testing"

	"github.com/caos/orbos/pkg/secret"
	"github.com/stretchr/testify/assert"
)

func TestBucket_getSecretsFull(t *testing.T) {
	secrets, existing := getSecretsMap(&desired)
	assert.Equal(t, desired.Spec.AccessKeyID, secrets["accesskeyid"])
	assert.Equal(t, desired.Spec.ExistingAccessKeyID, existing["accesskeyid"])
	assert.Equal(t, desired.Spec.SecretAccessKey, secrets["secretaccesskey"])
	assert.Equal(t, desired.Spec.ExistingSecretAccessKey, existing["secretaccesskey"])
}

func TestBucket_getSecretsEmpty(t *testing.T) {
	secrets, existing := getSecretsMap(&desiredWithoutSecret)
	assert.Equal(t, &secret.Secret{}, secrets["accesskeyid"])
	assert.Equal(t, &secret.Existing{}, existing["accesskeyid"])
	assert.Equal(t, &secret.Secret{}, secrets["secretaccesskey"])
	assert.Equal(t, &secret.Existing{}, existing["secretaccesskey"])
}

func TestBucket_getSecretsNil(t *testing.T) {
	secrets, existing := getSecretsMap(&desiredNil)
	assert.Equal(t, &secret.Secret{}, secrets["accesskeyid"])
	assert.Equal(t, &secret.Existing{}, existing["accesskeyid"])
	assert.Equal(t, &secret.Secret{}, secrets["secretaccesskey"])
	assert.Equal(t, &secret.Existing{}, existing["secretaccesskey"])
}
