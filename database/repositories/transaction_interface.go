package repositories

import "crypto-challenge/entities"

type TransactionRepository interface {
	Create(newTransaction *entities.Transaction) error
	FindByID(idToSearch string) (*entities.Transaction, error)
	FindAll() ([]*entities.Transaction, error)
	UpdateByID(updatedTransaction *entities.Transaction) error
	DeleteByID(idToDelete string) error
}
