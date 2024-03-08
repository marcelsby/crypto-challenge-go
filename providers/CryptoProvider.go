// TODO: Importar a chave de uma fonte externa
// TODO: Aprimorar o tratamento de erros para encaixar com a API REST
package providers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

type CryptoProvider struct {
	secretKey []byte
}

func NewCryptoProvider() *CryptoProvider {
	secretKey, err := hex.DecodeString("442c71674bbc3fcb5b9eed338f63521a1e6b1c352e87768377bd3bfa86048404")
	if err != nil {
		panic(err.Error())
	}

	return &CryptoProvider{secretKey: secretKey}
}

func (cp *CryptoProvider) Encrypt(toEncrypt []byte) string {
	block, err := aes.NewCipher(cp.secretKey)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	ciphertext := aesgcm.Seal(nil, nonce, toEncrypt, nil)

	return fmt.Sprintf("%x-%x", nonce, ciphertext)
}

func (cp *CryptoProvider) Decrypt(toDecrypt string) []byte {
	nonceWithCiphertextSplitted := strings.Split(toDecrypt, "-")

	nonce, _ := hex.DecodeString(nonceWithCiphertextSplitted[0])
	ciphertext, _ := hex.DecodeString(nonceWithCiphertextSplitted[1])

	block, err := aes.NewCipher(cp.secretKey)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	decrypted, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return decrypted
}
