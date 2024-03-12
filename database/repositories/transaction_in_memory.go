package repositories

import (
	"crypto-challenge/entities"
)

type TransactionInMemoryRepository struct {
	transactions []*entities.Transaction
}

func NewTransactionInMemoryRepository() *TransactionInMemoryRepository {
	return &TransactionInMemoryRepository{[]*entities.Transaction{}}
}

func (p *TransactionInMemoryRepository) Create(transaction *entities.Transaction) error {
	copy := *transaction
	p.transactions = append(p.transactions, &copy)

	return nil
}

func (p *TransactionInMemoryRepository) FindByID(idToSearch string) (*entities.Transaction, error) {
	for _, transaction := range p.transactions {
		if transaction.ID == idToSearch {
			transactionFound := *transaction
			return &transactionFound, nil
		}
	}

	return nil, nil
}

func (p *TransactionInMemoryRepository) FindAll() ([]*entities.Transaction, error) {
	foundTransactions := make([]*entities.Transaction, len(p.transactions))

	for index, transaction := range p.transactions {
		foundTransactions[index] = &entities.Transaction{
			ID:              transaction.ID,
			UserDocument:    transaction.UserDocument,
			CreditCardToken: transaction.CreditCardToken,
			Value:           transaction.Value,
		}
	}

	return foundTransactions, nil
}

func (p *TransactionInMemoryRepository) UpdateByID(idToUpdate string, updatedTransaction *entities.Transaction) error {
	for index, transaction := range p.transactions {
		if transaction.ID == idToUpdate {
			copy := *updatedTransaction
			p.transactions[index] = &copy
		}
	}

	return nil
}

func (p *TransactionInMemoryRepository) DeleteByID(idToDelete string) error {
	for index, transaction := range p.transactions {
		if transaction.ID == idToDelete {
			p.transactions = append(p.transactions[:index], p.transactions[index+1:]...)
			break
		}
	}

	return nil
}
