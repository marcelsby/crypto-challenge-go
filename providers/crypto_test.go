package providers

import "testing"

var secretKey = "c7b81104b9fc8b05ff85995f6d34d5b18cfbb0cff21ff2ceab154a3bcfae3aba"

func TestEncryptAndDecryptString(t *testing.T) {
	// given
	underTest := NewCryptoProvider(secretKey)
	testString := "lorem ipsum"

	// when
	ciphertext, err := underTest.Encrypt([]byte(testString))
	if err != nil {
		t.Error(err)
	}

	decrypted, err := underTest.Decrypt(ciphertext)
	if err != nil {
		t.Error(err)
	}

	// then
	if string(decrypted) != testString {
		t.Errorf("got: %s, expected: %s", decrypted, testString)
	}
}
