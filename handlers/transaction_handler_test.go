package handlers_test

import (
	"crypto-challenge/entities"
	"crypto-challenge/handlers"
	"crypto-challenge/mocks/crypto-challenge/database/repositories"
	"crypto-challenge/mocks/crypto-challenge/providers"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TransactionHandlerTestSuite struct {
	suite.Suite
	router             *chi.Mux
	repositoryMock     *repositories.MockTransactionRepository
	cryptoProviderMock *providers.MockTransactionCryptoProvider
}

func (ts *TransactionHandlerTestSuite) SetupTest() {
	ts.router = chi.NewRouter()

	ts.repositoryMock = repositories.NewMockTransactionRepository(ts.T())
	ts.cryptoProviderMock = providers.NewMockTransactionCryptoProvider(ts.T())

	ts.router.Mount("/", handlers.NewTransactionRouter(ts.repositoryMock, ts.cryptoProviderMock))
}

func (ts *TransactionHandlerTestSuite) TestCreate() {
	// given
	validNewTransactionJSON, err := generateRandomTransactionJSON(false, true)
	if err != nil {
		ts.T().Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/transactions", strings.NewReader(validNewTransactionJSON))
	res := httptest.NewRecorder()

	ts.cryptoProviderMock.EXPECT().Encrypt(mock.AnythingOfType("*entities.Transaction")).Return(nil).Once()
	ts.repositoryMock.EXPECT().Create(mock.AnythingOfType("*entities.Transaction")).Return(nil).Once()

	// when
	ts.router.ServeHTTP(res, req)

	// then
	ts.Require().Equal(http.StatusCreated, res.Code)
	ts.Require().Empty(res.Body.Bytes())
}

func (ts *TransactionHandlerTestSuite) TestCreate_WithInvalidRequestBody() {
	// given
	invalidNewTransactionJSON, err := generateRandomTransactionJSON(false, false)
	if err != nil {
		ts.T().Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/transactions", strings.NewReader(invalidNewTransactionJSON))
	res := httptest.NewRecorder()

	// when
	ts.router.ServeHTTP(res, req)

	// then
	ts.Require().Equal(http.StatusUnprocessableEntity, res.Code)
	ts.Assert().Empty(res.Body.Bytes())
}

func (ts *TransactionHandlerTestSuite) TestCreate_WithErrorOnEncryption() {
	// given
	validNewTransactionJSON, err := generateRandomTransactionJSON(false, true)
	if err != nil {
		ts.T().Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/transactions", strings.NewReader(validNewTransactionJSON))
	res := httptest.NewRecorder()

	ts.cryptoProviderMock.EXPECT().Encrypt(mock.AnythingOfType("*entities.Transaction")).
		Return(errors.New("error on encryption"))

	// when
	ts.router.ServeHTTP(res, req)

	// then
	ts.Require().Equal(http.StatusInternalServerError, res.Code)
	ts.Require().Equal("application/json", res.Header().Get("Content-Type"))
	ts.Require().NotEmpty(res.Body.Bytes())
	ts.Require().True(json.Valid(res.Body.Bytes()), "invalid JSON response. Received:", res.Body.String())
}

func (ts *TransactionHandlerTestSuite) TestCreate_WithErrorOnCreate() {
	// given
	validNewTransactionJSON, err := generateRandomTransactionJSON(false, true)
	if err != nil {
		ts.T().Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/transactions", strings.NewReader(validNewTransactionJSON))
	res := httptest.NewRecorder()

	ts.cryptoProviderMock.EXPECT().Encrypt(mock.AnythingOfType("*entities.Transaction")).Return(nil)
	ts.repositoryMock.EXPECT().Create(mock.AnythingOfType("*entities.Transaction")).Return(errors.New("error on create"))

	// when
	ts.router.ServeHTTP(res, req)

	// then
	ts.Require().Equal(http.StatusInternalServerError, res.Code)
	ts.Require().Equal("application/json", res.Header().Get("Content-Type"))

	ts.Require().NotEmpty(res.Body.Bytes())
	ts.Assert().True(json.Valid(res.Body.Bytes()), "invalid JSON response. Received:", res.Body.String())
}

func (ts *TransactionHandlerTestSuite) TestFindByID() {
	// given
	expectedTransaction := generateRandomTransaction(true)

	ts.repositoryMock.EXPECT().FindByID(expectedTransaction.ID).Return(expectedTransaction, nil).Once()
	ts.cryptoProviderMock.EXPECT().Decrypt(mock.AnythingOfType("*entities.Transaction")).Return(nil).Once()

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/transactions/%s", expectedTransaction.ID), nil)
	res := httptest.NewRecorder()

	// when
	ts.router.ServeHTTP(res, req)

	// then
	ts.Require().Equal(http.StatusOK, res.Code)
	ts.Require().Equal("application/json", res.Header().Get("Content-Type"))

	var actualTransaction *entities.Transaction
	err := json.Unmarshal(res.Body.Bytes(), &actualTransaction)
	if err != nil {
		ts.T().Fatal(err)
	}

	ts.Require().Equal(expectedTransaction, actualTransaction)
}

func (ts *TransactionHandlerTestSuite) TestFindByID_WithErrorOnFindByID() {
	// given
	randomID := uuid.NewString()

	ts.repositoryMock.EXPECT().FindByID(randomID).Return(nil, errors.New("error on FindByID"))

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/transactions/%s", randomID), nil)
	res := httptest.NewRecorder()

	// when
	ts.router.ServeHTTP(res, req)

	// then
	ts.Require().Equal(http.StatusInternalServerError, res.Code)
	ts.Require().Equal("application/json", res.Header().Get("Content-Type"))

	isValidResponseJSON := json.Valid(res.Body.Bytes())
	ts.Require().True(isValidResponseJSON, "invalid error response JSON payload.", res.Body.String())
}

func (ts *TransactionHandlerTestSuite) TestFindByID_WhenNotFound() {
	// given
	randomID := uuid.NewString()

	ts.repositoryMock.EXPECT().FindByID(randomID).Return(nil, nil).Once()

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/transactions/%s", randomID), nil)
	res := httptest.NewRecorder()

	// when
	ts.router.ServeHTTP(res, req)

	// then
	ts.Require().Equal(http.StatusNotFound, res.Code)
	ts.Require().Equal("application/json", res.Header().Get("Content-Type"))

	isValidResponseJSON := json.Valid(res.Body.Bytes())
	ts.Require().True(isValidResponseJSON, "invalid error response JSON payload.", res.Body.String())
}

func (ts *TransactionHandlerTestSuite) TestFindByID_WithErrorOnDecrypt() {
	// given
	randomID := uuid.NewString()

	ts.repositoryMock.EXPECT().FindByID(randomID).Return(&entities.Transaction{}, nil)
	ts.cryptoProviderMock.EXPECT().Decrypt(mock.AnythingOfType("*entities.Transaction")).
		Return(errors.New("error on Decrypt"))

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/transactions/%s", randomID), nil)
	res := httptest.NewRecorder()

	// when
	ts.router.ServeHTTP(res, req)

	// then
	ts.Require().Equal(http.StatusInternalServerError, res.Code)
	ts.Require().Equal("application/json", res.Header().Get("Content-Type"))

	isValidResponseJSON := json.Valid(res.Body.Bytes())
	ts.Require().True(isValidResponseJSON, "invalid error response JSON payload.", res.Body.String())
}

func (ts *TransactionHandlerTestSuite) TestFindAll() {
	// given
	expectedTransactions := []*entities.Transaction{
		generateRandomTransaction(true),
		generateRandomTransaction(true),
	}

	ts.repositoryMock.EXPECT().FindAll().Return(expectedTransactions, nil)
	ts.cryptoProviderMock.EXPECT().Decrypt(mock.AnythingOfType("*entities.Transaction")).
		Return(nil).Times(2)

	req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
	res := httptest.NewRecorder()

	// when
	ts.router.ServeHTTP(res, req)

	// then
	ts.Require().Equal(http.StatusOK, res.Code)
	ts.Require().Equal("application/json", res.Header().Get("Content-Type"))

	var actualTransactions []*entities.Transaction

	err := json.Unmarshal(res.Body.Bytes(), &actualTransactions)
	if err != nil {
		ts.T().Fatal(err)
	}

	ts.Require().Equal(expectedTransactions, actualTransactions)
}

func (ts *TransactionHandlerTestSuite) TestFindAll_WithErrorOnFindAll() {
	// given
	ts.repositoryMock.EXPECT().FindAll().Return(nil, errors.New("error on FindAll"))

	req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
	res := httptest.NewRecorder()

	// when
	ts.router.ServeHTTP(res, req)

	// then
	ts.Require().Equal(http.StatusInternalServerError, res.Code)
	ts.Require().Equal("application/json", res.Header().Get("Content-Type"))

	isValidResponseJSON := json.Valid(res.Body.Bytes())
	ts.Require().True(isValidResponseJSON)
}

func (ts *TransactionHandlerTestSuite) TestFindAll_WithErrorOnDecrypt() {
	// given
	transactions := []*entities.Transaction{
		generateRandomTransaction(true),
	}

	ts.repositoryMock.EXPECT().FindAll().Return(transactions, nil)
	ts.cryptoProviderMock.EXPECT().Decrypt(transactions[0]).Return(errors.New("error on Decrypt"))

	req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
	res := httptest.NewRecorder()

	// when
	ts.router.ServeHTTP(res, req)

	// then
	ts.Require().Equal(http.StatusInternalServerError, res.Code)
	ts.Require().Equal("application/json", res.Header().Get("Content-Type"))

	isValidResponseJSON := json.Valid(res.Body.Bytes())
	ts.Require().True(isValidResponseJSON)
}

func TestTransactionHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionHandlerTestSuite))
}

func generateRandomTransactionJSON(withID, validJSON bool) (string, error) {
	t := generateRandomTransaction(withID)

	tBytesJSON, err := json.Marshal(t)
	tStrJSON := string(tBytesJSON)

	if !validJSON {
		tStrJSON = strings.ReplaceAll(tStrJSON, `"`, "")
	}

	return tStrJSON, err
}

func generateRandomTransaction(withID bool) *entities.Transaction {
	fakeUserDocuments := []string{"50277613433", "19318615400", "43872034856", "25694674300", "56214093854",
		"01927386406", "89673401520", "73619405823", "40198237610", "58327490120"}
	randomUserDocument := fakeUserDocuments[rand.Intn(len(fakeUserDocuments))]

	// Get random 3 digits number to simulate the credit card CVV code
	randomCreditCardToken := fmt.Sprint(rand.Intn(900) + 100)

	// Get a float64 with 2 floating point precision and max of 4 integer digits - DECIMAL (6,2)
	integerPart := rand.Intn(10000)
	decimalPart := rand.Intn(100)
	randomValue := float64(integerPart) + float64(decimalPart)/100

	var id string
	if withID {
		id = uuid.NewString()
	}

	return &entities.Transaction{
		ID:              id,
		UserDocument:    randomUserDocument,
		CreditCardToken: randomCreditCardToken,
		Value:           randomValue,
	}
}
