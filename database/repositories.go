package database

import "crypto-challenge/models"

type TransactionInMemoryRepository struct {
	people []*models.Transaction
}

func NewTransactionInMemoryRepository() *TransactionInMemoryRepository {
	return &TransactionInMemoryRepository{[]*models.Transaction{}}
}

func (p *TransactionInMemoryRepository) Create(person *models.Transaction) {
	p.people = append(p.people, person)
}

func (p *TransactionInMemoryRepository) FindByID(id string) *models.Transaction {
	for _, person := range p.people {
		if person.ID == id {
			return person
		}
	}

	return nil
}

func (p *TransactionInMemoryRepository) FindAll() []*models.Transaction {
	return p.people
}
