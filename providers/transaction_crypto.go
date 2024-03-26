package providers

import (
	"crypto-challenge/entities"
)

type TransactionCryptoProvider interface {
	Encrypt(*entities.Transaction) error
	Decrypt(*entities.Transaction) error
}

type StandardTransactionCryptoProvider struct {
	cp CryptoProvider
}

func NewStandardTransactionCryptoProvider(cp CryptoProvider) *StandardTransactionCryptoProvider {
	return &StandardTransactionCryptoProvider{cp}
}

func (tcp *StandardTransactionCryptoProvider) Encrypt(toEncrypt *entities.Transaction) error {
	encryptedUserDocument, err := tcp.cp.Encrypt([]byte(toEncrypt.UserDocument))
	if err != nil {
		return err
	}

	encryptedCreditCardToken, err := tcp.cp.Encrypt([]byte(toEncrypt.CreditCardToken))
	if err != nil {
		return err
	}

	toEncrypt.UserDocument = encryptedUserDocument
	toEncrypt.CreditCardToken = encryptedCreditCardToken

	return nil
}

func (tcp *StandardTransactionCryptoProvider) Decrypt(toDecrypt *entities.Transaction) error {
	decryptedUserDocument, err := tcp.cp.Decrypt(toDecrypt.UserDocument)
	if err != nil {
		return err
	}

	decryptedCreditCardToken, err := tcp.cp.Decrypt(toDecrypt.CreditCardToken)
	if err != nil {
		return err
	}

	toDecrypt.UserDocument = string(decryptedUserDocument)
	toDecrypt.CreditCardToken = string(decryptedCreditCardToken)

	return nil
}
