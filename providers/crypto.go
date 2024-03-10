// TODO: Importar a chave de uma fonte externa
package providers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"strings"
)

type CryptoProvider struct {
	key []byte
}

func NewCryptoProvider(secretKey string) *CryptoProvider {
	key, _ := hex.DecodeString(secretKey)

	return &CryptoProvider{key}
}

func (cp *CryptoProvider) Encrypt(toEncrypt []byte) (*string, error) {
	block, err := aes.NewCipher(cp.key)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Println(err)
		return nil, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, toEncrypt, nil)

	nonceAndCiphertext := fmt.Sprintf("%x-%x", nonce, ciphertext)

	return &nonceAndCiphertext, nil
}

func (cp *CryptoProvider) Decrypt(toDecrypt string) (*[]byte, error) {
	nonceWithCiphertextSplitted := strings.Split(toDecrypt, "-")

	nonce, _ := hex.DecodeString(nonceWithCiphertextSplitted[0])
	ciphertext, _ := hex.DecodeString(nonceWithCiphertextSplitted[1])

	block, err := aes.NewCipher(cp.key)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	decrypted, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &decrypted, nil
}
