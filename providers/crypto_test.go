package providers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var secretKey = "c7b81104b9fc8b05ff85995f6d34d5b18cfbb0cff21ff2ceab154a3bcfae3aba"

func TestEncryptAndDecryptString(t *testing.T) {
	// given
	underTest := NewCryptoProvider(secretKey)
	expected := "lorem ipsum"

	// when
	ciphertext, err := underTest.Encrypt([]byte(expected))
	require.Nil(t, err)

	actual, err := underTest.Decrypt(ciphertext)
	require.Nil(t, err)

	// then
	assert.Equal(t, expected, string(actual))
}
