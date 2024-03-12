package repositories

import (
	"crypto-challenge/entities"
	"testing"

	"github.com/google/uuid"
)

func TestCreate(t *testing.T) {
	// given
	underTest := NewTransactionInMemoryRepository()
	expected := createTransaction()

	// when
	err := underTest.Create(expected)
	if err != nil {
		t.Error(err)
	}

	// then
	got := underTest.transactions[0]

	if *got != *expected {
		t.Fatalf("got: %+v\n expected: %+v", got, expected)
	}

	if got == expected {
		t.Fatalf("the stored Transaction must be a copy of the passed Transaction, and not a pointer to it.")
	}
}

func TestFindByID(t *testing.T) {
	// given
	underTest := NewTransactionInMemoryRepository()
	transaction := createTransaction()
	underTest.transactions = append(underTest.transactions, transaction)
	transaction.ID = uuid.NewString()

	// when
	got, err := underTest.FindByID(transaction.ID)
	if err != nil {
		t.Error(err)
	}

	// then
	if got == transaction {
		t.Fatalf("the found Transaction must be a copy of the passed Transaction, and not a pointer to it.")
	}
}

func TestFindAll(t *testing.T) {
	// given
	repo := NewTransactionInMemoryRepository()

	transaction1 := createTransaction()

	transaction2 := createTransaction()
	transaction2.ID = uuid.NewString()

	repo.transactions = append(repo.transactions, transaction1, transaction2)

	// when
	got, err := repo.FindAll()
	if err != nil {
		t.Error(err)
	}

	// then
	if len(got) != 2 {
		t.Errorf("len(got) - expected: 2, got: %d", len(got))
	}

	if got[0] == transaction1 || got[1] == transaction2 {
		t.Fatalf("the returned slice must contain a pointer to a copy of each stored Transaction, and not a pointer to the stored Transaction.")
	}
}

func TestUpdateByID(t *testing.T) {
	// given
	underTest := NewTransactionInMemoryRepository()
	transaction := createTransaction()
	underTest.transactions = append(underTest.transactions, transaction)

	updatedTransaction := createTransaction()
	updatedTransaction.CreditCardToken = "783"
	updatedTransaction.UserDocument = "39474452024"
	updatedTransaction.Value = 5999.89

	// when
	err := underTest.UpdateByID(transaction.ID, updatedTransaction)
	if err != nil {
		t.Error(err)
	}

	// then
	if underTest.transactions[0] == updatedTransaction {
		t.Error("The stored updated Transaction must be a copy that was made from the parameter 'updatedTransaction' passed.")
	}

	if *underTest.transactions[0] != *updatedTransaction {
		t.Errorf("expected: %+v\n got: %+v", updatedTransaction, *underTest.transactions[0])
	}
}

func TestDeleteByID(t *testing.T) {
	// given
	underTest := NewTransactionInMemoryRepository()
	transaction := createTransaction()
	underTest.transactions = append(underTest.transactions, transaction)

	// when
	err := underTest.DeleteByID(transaction.ID)
	if err != nil {
		t.Error(err)
	}

	// then
	if len(underTest.transactions) != 0 {
		t.Errorf("Transaction not deleted, got underTest.transactions: %v", underTest.transactions)
	}
}

func TestDeleteByIDWith3Transactions(t *testing.T) {
	// given
	underTest := NewTransactionInMemoryRepository()

	id1, id2, id3 := uuid.NewString(), uuid.NewString(), uuid.NewString()

	transactions := []*entities.Transaction{
		createTransactionWithID(id1),
		createTransactionWithID(id2),
		createTransactionWithID(id3),
	}

	underTest.transactions = transactions

	// when
	err := underTest.DeleteByID(id2)
	if err != nil {
		t.Error(err)
	}

	// then
	if len(underTest.transactions) != 2 {
		t.Errorf("the deletion isn't working, got len(underTest.transactions): %d, expected: 2", len(underTest.transactions))
	}

	if underTest.transactions[0].ID != id1 && underTest.transactions[1].ID != id3 {
		t.Errorf("the first and second remaining Transactions IDs isn't matching. Expected: [0]: %s [1]: %s, got: [0]: %s, [1]: %s", id1, id2,
			underTest.transactions[0].ID, underTest.transactions[1].ID)
	}
}

func createTransactionWithID(id string) *entities.Transaction {
	return &entities.Transaction{
		ID:              id,
		UserDocument:    "46402249076",
		CreditCardToken: "337",
		Value:           250.33,
	}
}

func createTransaction() *entities.Transaction {
	return &entities.Transaction{
		ID:              uuid.NewString(),
		UserDocument:    "46402249076",
		CreditCardToken: "337",
		Value:           250.33,
	}
}
