package repositories

import (
	"crypto-challenge/entities"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type TransactionInMemoryUnitTestSuite struct {
	suite.Suite
	underTest *TransactionInMemoryRepository
}

func (ts *TransactionInMemoryUnitTestSuite) SetupTest() {
	ts.underTest = NewTransactionInMemoryRepository()
}

func (ts *TransactionInMemoryUnitTestSuite) TestCreate() {
	// given
	expected := createTransaction()

	// when
	err := ts.underTest.Create(expected)
	ts.Nil(err)

	// then
	actual := ts.underTest.transactions[0]

	ts.Equal(*expected, *actual)

	ts.NotSame(expected, actual, "The stored Transaction must be a copy of the passed Transaction, and not a pointer to it.")
}

func (ts *TransactionInMemoryUnitTestSuite) TestFindByID() {
	// given
	expected := createTransaction()
	ts.underTest.transactions = append(ts.underTest.transactions, expected)

	// when
	actual, err := ts.underTest.FindByID(expected.ID)
	ts.Nil(err)

	// then
	ts.Equal(expected, actual)

	ts.NotSame(expected, actual, "The returned Transaction must be a copy of the stored Transaction, and not a pointer to it.")
}

func (ts *TransactionInMemoryUnitTestSuite) TestFindAll() {
	// given
	transaction1 := createTransaction()
	transaction2 := createTransaction()

	ts.underTest.transactions = append(ts.underTest.transactions, transaction1, transaction2)

	expected := []*entities.Transaction{transaction1, transaction2}

	// when
	actual, err := ts.underTest.FindAll()
	ts.Nil(err)

	// then
	ts.Len(actual, 2)

	ts.ElementsMatch(actual, expected)

	ts.NotSame(actual[0], transaction1, "The returned Transaction must be a copy of the stored Transaction, and not a pointer to it.")
	ts.NotSame(actual[1], transaction2, "The returned Transaction must be a copy of the stored Transaction, and not a pointer to it.")
}

func (ts *TransactionInMemoryUnitTestSuite) TestUpdateByID() {
	// given
	transaction := createTransaction()
	ts.underTest.transactions = append(ts.underTest.transactions, transaction)

	updatedTransaction := &entities.Transaction{
		ID:              transaction.ID,
		CreditCardToken: "783",
		UserDocument:    "39474452024",
		Value:           5999.89,
	}

	// when
	err := ts.underTest.UpdateByID(updatedTransaction)
	ts.Nil(err)

	// then
	ts.Equal(ts.underTest.transactions[0], updatedTransaction)

	ts.NotSame(ts.underTest.transactions[0], updatedTransaction,
		"The stored updated Transaction must be a copy that was made from the parameter 'updatedTransaction' passed.")
}

func (ts *TransactionInMemoryUnitTestSuite) TestDeleteByID() {
	// given
	transaction := createTransaction()
	ts.underTest.transactions = append(ts.underTest.transactions, transaction)

	// when
	err := ts.underTest.DeleteByID(transaction.ID)
	ts.Nil(err)

	// then
	ts.Len(ts.underTest.transactions, 0)
}

func (ts *TransactionInMemoryUnitTestSuite) TestDeleteByID_With3Transactions() {
	// given
	id1, id2, id3 := uuid.NewString(), uuid.NewString(), uuid.NewString()

	transactions := []*entities.Transaction{
		createTransactionWithID(id1),
		createTransactionWithID(id2),
		createTransactionWithID(id3),
	}

	ts.underTest.transactions = transactions

	// when
	err := ts.underTest.DeleteByID(id2)
	ts.Nil(err)

	// then
	ts.Len(ts.underTest.transactions, 2)

	ts.Equal(ts.underTest.transactions[0], transactions[0])
	ts.Equal(ts.underTest.transactions[1], transactions[2])
}

func TestTransactionInMemoryUnitTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionInMemoryUnitTestSuite))
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
