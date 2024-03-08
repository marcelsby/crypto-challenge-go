package database

import (
	"crypto-challenge/models"
)

type TransactionInMemoryRepository struct {
	transactions []*models.Transaction
}

func NewTransactionInMemoryRepository() *TransactionInMemoryRepository {
	return &TransactionInMemoryRepository{[]*models.Transaction{}}
}

func (p *TransactionInMemoryRepository) Create(transaction *models.Transaction) {
	p.transactions = append(p.transactions, transaction)
}

func (p *TransactionInMemoryRepository) FindByID(id string) *models.Transaction {
	for _, transaction := range p.transactions {
		if transaction.ID == id {
			transactionFound := *transaction
			return &transactionFound
		}
	}

	return nil
}

func (p *TransactionInMemoryRepository) FindAll() []*models.Transaction {
	foundTransactions := make([]*models.Transaction, len(p.transactions))

	for index, transaction := range p.transactions {
		foundTransactions[index] = &models.Transaction{
			ID:              transaction.ID,
			UserDocument:    transaction.UserDocument,
			CreditCardToken: transaction.CreditCardToken,
			Value:           transaction.Value,
		}
	}

	return foundTransactions
}

func (p *TransactionInMemoryRepository) Update(id string, updatedTransaction *models.Transaction) {
	for index, transaction := range p.transactions {
		if transaction.ID == id {
			p.transactions[index] = updatedTransaction
		}
	}
}

func (p *TransactionInMemoryRepository) DeleteByID(id string) {
	for index, transaction := range p.transactions {
		if transaction.ID == id {
			p.transactions = append(p.transactions[:index], p.transactions[index+1:]...)
			break
		}
	}
}
