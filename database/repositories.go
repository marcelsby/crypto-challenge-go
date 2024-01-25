package database

import "crypto-challenge/models"

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
			return transaction
		}
	}

	return nil
}

func (p *TransactionInMemoryRepository) FindAll() []*models.Transaction {
	return p.transactions
}

func (p *TransactionInMemoryRepository) Update(id string, updatedTransaction *models.Transaction) {
	for index, transaction := range p.transactions {
		if transaction.ID == id {
			p.transactions[index] = updatedTransaction
		}
	}
}
