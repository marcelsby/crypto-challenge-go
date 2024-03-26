package repositories_test

import (
	"crypto-challenge/config"
	"crypto-challenge/database/repositories"
	"crypto-challenge/entities"
	"crypto-challenge/testhelpers"
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type TransactionMySqlIntTestSuite struct {
	suite.Suite
	terminateMySqlContainer *func()
	db                      *sql.DB
	underTest               *repositories.TransactionMySqlRepository
}

func (ts *TransactionMySqlIntTestSuite) SetupSuite() {
	dotenvFilePath, err := filepath.Abs(filepath.Join("..", "..", ".env"))
	if err != nil {
		ts.T().Fatal(err)
	}

	cfg := config.GetAppConfig(dotenvFilePath)

	migrationsFolderPath, err := filepath.Abs(filepath.Join("..", "..", ".docker", "sql"))
	if err != nil {
		ts.T().Fatal(err)
	}

	mySqlC, terminateMySqlC, ctxMySqlC := testhelpers.SetupMySqlContainer(cfg, migrationsFolderPath)

	ts.terminateMySqlContainer = terminateMySqlC

	db := testhelpers.GetMySqlContainerDB(ts.T(), mySqlC, ctxMySqlC, cfg)

	ts.db = db

	ts.underTest = repositories.NewTransactionMySqlRepository(db)
}

func (ts *TransactionMySqlIntTestSuite) TearDownSuite() {
	(*ts.terminateMySqlContainer)()
	ts.db.Close()
}

func (ts *TransactionMySqlIntTestSuite) SetupTest() {
	_, err := ts.db.Exec("DELETE FROM transactions")
	ts.Nil(err)
}

func (ts *TransactionMySqlIntTestSuite) TestCreate() {
	//given
	expected := createTransaction()

	//when
	err := ts.underTest.Create(&expected)
	ts.Nil(err)

	//then
	actual := entities.Transaction{}
	err = ts.db.QueryRow("SELECT id, user_document, credit_card_token, `value` FROM transactions WHERE id = ?", expected.ID).Scan(
		&actual.ID,
		&actual.UserDocument,
		&actual.CreditCardToken,
		&actual.Value,
	)
	ts.Nil(err)

	ts.Equal(expected, actual)
}

func (ts *TransactionMySqlIntTestSuite) TestFindByID() {
	//given
	expected := createTransaction()
	_, err := ts.db.Exec("INSERT INTO transactions (id, user_document, credit_card_token, `value`) VALUES (?, ?, ?, ?)",
		expected.ID, expected.UserDocument, expected.CreditCardToken, expected.Value)
	ts.Nil(err)

	//when
	actual, err := ts.underTest.FindByID(expected.ID)
	ts.Nil(err)

	//then
	ts.Equal(expected, *actual)
}

func (ts *TransactionMySqlIntTestSuite) TestFindByID_WhenNotFound() {
	//given
	idToSearch := uuid.NewString()

	//when
	actual, err := ts.underTest.FindByID(idToSearch)
	ts.Nil(err)

	//then
	ts.Nil(actual)
}

func (ts *TransactionMySqlIntTestSuite) TestFindAll() {
	//given
	expected1, expected2 := createTransaction(), createTransaction()

	_, err := ts.db.Exec("INSERT INTO transactions (id, user_document, credit_card_token, `value`) VALUES (?, ?, ?, ?), (?, ?, ?, ?)",
		expected1.ID, expected1.UserDocument, expected1.CreditCardToken, expected1.Value,
		expected2.ID, expected2.UserDocument, expected2.CreditCardToken, expected2.Value)
	ts.Nil(err)

	//when
	actual, err := ts.underTest.FindAll()
	ts.Nil(err)

	//then
	ts.Len(actual, 2)

	expectedIDs := []string{expected1.ID, expected2.ID}

	for _, foundTransaction := range actual {
		ts.Contains(expectedIDs, foundTransaction.ID)
	}
}

func (ts *TransactionMySqlIntTestSuite) TestFindAll_WhenEmpty() {
	//when
	actual, err := ts.underTest.FindAll()
	ts.Nil(err)

	//then
	ts.Empty(actual)
}

func (ts *TransactionMySqlIntTestSuite) TestUpdateByID() {
	//given
	newTransaction := createTransaction()

	_, err := ts.db.Exec("INSERT INTO transactions (id, user_document, credit_card_token, `value`) VALUES (?, ?, ?, ?)",
		newTransaction.ID, newTransaction.UserDocument, newTransaction.CreditCardToken, newTransaction.Value)
	ts.Nil(err)

	expected := entities.Transaction{
		ID:              newTransaction.ID,
		Value:           299.99,
		UserDocument:    "27184927",
		CreditCardToken: "663",
	}

	//when
	err = ts.underTest.UpdateByID(&expected)
	ts.Nil(err)

	//then
	actual := entities.Transaction{}
	err = ts.db.QueryRow("SELECT id, user_document, credit_card_token, `value` FROM transactions WHERE id = ?", newTransaction.ID).Scan(
		&actual.ID,
		&actual.UserDocument,
		&actual.CreditCardToken,
		&actual.Value,
	)
	ts.Nil(err)

	ts.Equal(expected, actual)
}

func (ts *TransactionMySqlIntTestSuite) TestDeleteByID() {
	//given
	newTransaction := createTransaction()

	_, err := ts.db.Exec("INSERT INTO transactions (id, user_document, credit_card_token, `value`) VALUES (?, ?, ?, ?)",
		newTransaction.ID, newTransaction.UserDocument, newTransaction.CreditCardToken, newTransaction.Value)
	ts.Nil(err)

	//when
	err = ts.underTest.DeleteByID(newTransaction.ID)
	ts.Nil(err)

	//then
	var actual int

	err = ts.db.QueryRow("SELECT count(*) FROM transactions").Scan(&actual)
	ts.Nil(err)

	ts.Zero(actual)
}

func TestTransactionMySqlIntTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionMySqlIntTestSuite))
}

func createTransaction() entities.Transaction {
	return entities.Transaction{
		ID:              uuid.NewString(),
		UserDocument:    "12345",
		CreditCardToken: "755",
		Value:           9999.99,
	}
}
