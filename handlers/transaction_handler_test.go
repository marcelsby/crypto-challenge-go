package handlers_test

import (
	"crypto-challenge/entities"
	"crypto-challenge/handlers"
	"crypto-challenge/mocks/crypto-challenge/database/repositories"
	"crypto-challenge/mocks/crypto-challenge/providers"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const InvalidJSONResponsePayload = "invalid JSON response payload."

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

	ts.cryptoProviderMock.EXPECT().Encrypt(mock.AnythingOfType("*entities.Transaction")).Return(nil).Once()
	ts.repositoryMock.EXPECT().Create(mock.AnythingOfType("*entities.Transaction")).Return(nil).Once()

	// when
	res := makeRequest(ts.router, http.MethodPost, "/transactions", strings.NewReader(validNewTransactionJSON))

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

	// when
	res := makeRequest(ts.router, http.MethodPost, "/transactions", strings.NewReader(invalidNewTransactionJSON))

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

	ts.cryptoProviderMock.EXPECT().Encrypt(mock.AnythingOfType("*entities.Transaction")).
		Return(errorOnMethod("Encrypt"))

	// when
	res := makeRequest(ts.router, http.MethodPost, "/transactions", strings.NewReader(validNewTransactionJSON))

	// then
	ts.Require().Equal(http.StatusInternalServerError, res.Code)
	ts.Require().Equal("application/json", res.Header().Get("Content-Type"))
	ts.Require().NotEmpty(res.Body.Bytes())
	requireValidJSON(ts.T(), res.Body.Bytes(), "invalid error response JSON payload.",
		res.Body.String())
}

func (ts *TransactionHandlerTestSuite) TestCreate_WithErrorOnCreate() {
	// given
	validNewTransactionJSON, err := generateRandomTransactionJSON(false, true)
	if err != nil {
		ts.T().Fatal(err)
	}

	ts.cryptoProviderMock.EXPECT().Encrypt(mock.AnythingOfType("*entities.Transaction")).Return(nil)
	ts.repositoryMock.EXPECT().Create(mock.AnythingOfType("*entities.Transaction")).Return(errorOnMethod("create"))

	// when
	res := makeRequest(ts.router, http.MethodPost, "/transactions", strings.NewReader(validNewTransactionJSON))

	// then
	ts.Require().Equal(http.StatusInternalServerError, res.Code)
	ts.Require().Equal("application/json", res.Header().Get("Content-Type"))

	ts.Require().NotEmpty(res.Body.Bytes())
	requireValidJSON(ts.T(), res.Body.Bytes(), "invalid error response JSON payload.",
		res.Body.String())
}

func (ts *TransactionHandlerTestSuite) TestFindByID() {
	// given
	expectedTransaction := generateRandomTransaction(true)

	ts.repositoryMock.EXPECT().FindByID(expectedTransaction.ID).Return(expectedTransaction, nil).Once()
	ts.cryptoProviderMock.EXPECT().Decrypt(mock.AnythingOfType("*entities.Transaction")).Return(nil).Once()

	// when
	res := makeRequest(ts.router, http.MethodGet,
		fmt.Sprintf("/transactions/%s", expectedTransaction.ID), nil)

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

	ts.repositoryMock.EXPECT().FindByID(randomID).Return(nil, errorOnMethod("FindByID"))

	// when
	res := makeRequest(ts.router, http.MethodGet, fmt.Sprintf("/transactions/%s", randomID), nil)

	// then
	ts.Require().Equal(http.StatusInternalServerError, res.Code)
	ts.Require().Equal("application/json", res.Header().Get("Content-Type"))

	requireValidJSON(ts.T(), res.Body.Bytes(), "invalid error response JSON payload.",
		res.Body.String())
}

func (ts *TransactionHandlerTestSuite) TestFindByID_WhenNotFound() {
	// given
	randomID := uuid.NewString()

	ts.repositoryMock.EXPECT().FindByID(randomID).Return(nil, nil).Once()

	// when
	res := makeRequest(ts.router, http.MethodGet, fmt.Sprintf("/transactions/%s", randomID), nil)

	// then
	ts.Require().Equal(http.StatusNotFound, res.Code)
	ts.Require().Equal("application/json", res.Header().Get("Content-Type"))

	requireValidJSON(ts.T(), res.Body.Bytes(), InvalidJSONResponsePayload, res.Body.String())
}

func (ts *TransactionHandlerTestSuite) TestFindByID_WithErrorOnDecrypt() {
	// given
	randomID := uuid.NewString()

	ts.repositoryMock.EXPECT().FindByID(randomID).Return(&entities.Transaction{}, nil)
	ts.cryptoProviderMock.EXPECT().Decrypt(mock.AnythingOfType("*entities.Transaction")).
		Return(errorOnMethod("Decrypt"))

	// when
	res := makeRequest(ts.router, http.MethodGet, fmt.Sprintf("/transactions/%s", randomID), nil)

	// then
	ts.Require().Equal(http.StatusInternalServerError, res.Code)
	ts.Require().Equal("application/json", res.Header().Get("Content-Type"))

	requireValidJSON(ts.T(), res.Body.Bytes(), InvalidJSONResponsePayload, res.Body.String())
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

	// when
	res := makeRequest(ts.router, http.MethodGet, "/transactions", nil)

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
	ts.repositoryMock.EXPECT().FindAll().Return(nil, errorOnMethod("FindAll"))

	// when
	res := makeRequest(ts.router, http.MethodGet, "/transactions", nil)

	// then
	ts.Require().Equal(http.StatusInternalServerError, res.Code)
	ts.Require().Equal("application/json", res.Header().Get("Content-Type"))

	requireValidJSON(ts.T(), res.Body.Bytes(), InvalidJSONResponsePayload, res.Body.String())
}

func (ts *TransactionHandlerTestSuite) TestFindAll_WithErrorOnDecrypt() {
	// given
	transactions := []*entities.Transaction{
		generateRandomTransaction(true),
	}

	ts.repositoryMock.EXPECT().FindAll().Return(transactions, nil)
	ts.cryptoProviderMock.EXPECT().Decrypt(transactions[0]).Return(errorOnMethod("Decrypt"))

	// when
	res := makeRequest(ts.router, http.MethodGet, "/transactions", nil)

	// then
	ts.Require().Equal(http.StatusInternalServerError, res.Code)
	ts.Require().Equal("application/json", res.Header().Get("Content-Type"))

	requireValidJSON(ts.T(), res.Body.Bytes(), InvalidJSONResponsePayload, res.Body.String())
}

func (ts *TransactionHandlerTestSuite) TestUpdateByID() {
	// given
	randomID := uuid.NewString()

	updatedTransaction, err := generateRandomTransactionJSON(false, true)
	if err != nil {
		ts.T().Fatal(err)
	}

	ts.repositoryMock.EXPECT().FindByID(randomID).Return(&entities.Transaction{}, nil)
	ts.cryptoProviderMock.EXPECT().Encrypt(mock.AnythingOfType("*entities.Transaction")).
		Return(nil)
	ts.repositoryMock.EXPECT().UpdateByID(mock.AnythingOfType("*entities.Transaction")).
		Return(nil)

	// when
	res := makeRequest(ts.router, http.MethodPut, fmt.Sprintf("/transactions/%s", randomID),
		strings.NewReader(updatedTransaction))

	// then
	ts.Require().Equal(http.StatusOK, res.Code)
	ts.Require().Equal("application/json", res.Header().Get("Content-Type"))

	ts.Require().Empty(res.Body.Bytes())
}

func (ts *TransactionHandlerTestSuite) TestUpdateByID_WhenNotFound() {
	// given
	randomID := "abc"

	updatedTransaction, err := generateRandomTransactionJSON(false, true)
	if err != nil {
		ts.T().Fatal(err)
	}

	ts.repositoryMock.EXPECT().FindByID(randomID).Return(nil, nil)

	// when
	res := makeRequest(ts.router, http.MethodPut, fmt.Sprintf("/transactions/%s", randomID),
		strings.NewReader(updatedTransaction))

	// then
	ts.Require().Equal(http.StatusNotFound, res.Code)
	ts.Require().Equal("application/json", res.Header().Get("Content-Type"))

	requireValidJSON(ts.T(), res.Body.Bytes(), InvalidJSONResponsePayload, res.Body.String())
}

func (ts *TransactionHandlerTestSuite) TestUpdateByID_WithErrorOnFindByID() {
	// given
	randomID := uuid.NewString()

	updatedTransaction, err := generateRandomTransactionJSON(false, true)
	if err != nil {
		ts.T().Fatal(err)
	}

	ts.repositoryMock.EXPECT().FindByID(randomID).Return(nil, errorOnMethod("FindByID"))

	// when
	res := makeRequest(ts.router, http.MethodPut, fmt.Sprintf("/transactions/%s", randomID),
		strings.NewReader(updatedTransaction))

	// then
	ts.Require().Equal(http.StatusInternalServerError, res.Code)
	ts.Require().Equal("application/json", res.Header().Get("Content-Type"))

	requireValidJSON(ts.T(), res.Body.Bytes(), InvalidJSONResponsePayload, res.Body.String())
}

func (ts *TransactionHandlerTestSuite) TestUpdateByID_WithInvalidBody() {
	// given
	randomID := uuid.NewString()

	updatedTransaction, err := generateRandomTransactionJSON(false, false)
	if err != nil {
		ts.T().Fatal(err)
	}

	ts.repositoryMock.EXPECT().FindByID(randomID).Return(&entities.Transaction{}, nil)

	// when
	res := makeRequest(ts.router, http.MethodPut, fmt.Sprintf("/transactions/%s", randomID),
		strings.NewReader(updatedTransaction))

	// then
	ts.Require().Equal(http.StatusUnprocessableEntity, res.Code)
	ts.Require().Equal("application/json", res.Header().Get("Content-Type"))

	ts.Require().Empty(res.Body.Bytes())
}

func (ts *TransactionHandlerTestSuite) TestUpdateByID_WithErrorOnEncrypt() {
	// given
	randomID := uuid.NewString()

	updatedTransaction, err := generateRandomTransactionJSON(false, true)
	if err != nil {
		ts.T().Fatal(err)
	}

	ts.repositoryMock.EXPECT().FindByID(randomID).Return(&entities.Transaction{}, nil)
	ts.cryptoProviderMock.EXPECT().Encrypt(mock.AnythingOfType("*entities.Transaction")).
		Return(errorOnMethod("Encrypt"))

	// when
	res := makeRequest(ts.router, http.MethodPut, fmt.Sprintf("/transactions/%s", randomID),
		strings.NewReader(updatedTransaction))

	// then
	ts.Require().Equal(http.StatusInternalServerError, res.Code)
	ts.Require().Equal("application/json", res.Header().Get("Content-Type"))

	requireValidJSON(ts.T(), res.Body.Bytes(), InvalidJSONResponsePayload, res.Body.String())
}

func (ts *TransactionHandlerTestSuite) TestUpdateByID_WithErrorOnUpdateByID() {
	// given
	randomID := uuid.NewString()

	updatedTransaction, err := generateRandomTransactionJSON(false, true)
	if err != nil {
		ts.T().Fatal(err)
	}

	ts.repositoryMock.EXPECT().FindByID(randomID).Return(&entities.Transaction{}, nil)
	ts.cryptoProviderMock.EXPECT().Encrypt(mock.AnythingOfType("*entities.Transaction")).Return(nil)
	ts.repositoryMock.EXPECT().UpdateByID(mock.AnythingOfType("*entities.Transaction")).
		Return(errorOnMethod("UpdateByID"))

	// when
	res := makeRequest(ts.router, http.MethodPut, fmt.Sprintf("/transactions/%s", randomID),
		strings.NewReader(updatedTransaction))

	// then
	ts.Require().Equal(http.StatusInternalServerError, res.Code)
	ts.Require().Equal("application/json", res.Header().Get("Content-Type"))

	requireValidJSON(ts.T(), res.Body.Bytes(), InvalidJSONResponsePayload, res.Body.String())
}

func (ts *TransactionHandlerTestSuite) TestDeleteByID() {
	// given
	randomID := uuid.NewString()

	ts.repositoryMock.EXPECT().FindByID(randomID).Return(&entities.Transaction{}, nil)
	ts.repositoryMock.EXPECT().DeleteByID(randomID).Return(nil)

	// when
	res := makeRequest(ts.router, http.MethodDelete, fmt.Sprintf("/transactions/%s", randomID), nil)

	// then
	ts.Require().Equal(http.StatusOK, res.Code)
	ts.Require().Empty(res.Body.Bytes())
}

func (ts *TransactionHandlerTestSuite) TestDeleteByID_WithErrorOnFindByID() {
	// given
	randomID := uuid.NewString()

	ts.repositoryMock.EXPECT().FindByID(randomID).Return(nil, errorOnMethod("DeleteByID"))

	// when
	res := makeRequest(ts.router, http.MethodDelete, fmt.Sprintf("/transactions/%s", randomID), nil)

	// then
	ts.Require().Equal(http.StatusInternalServerError, res.Code)
	requireValidJSON(ts.T(), res.Body.Bytes(), InvalidJSONResponsePayload, res.Body.String())
}

func (ts *TransactionHandlerTestSuite) TestDeleteByID_WhenNotFound() {
	// given
	randomID := uuid.NewString()

	ts.repositoryMock.EXPECT().FindByID(randomID).Return(nil, nil)

	// when
	res := makeRequest(ts.router, http.MethodDelete, fmt.Sprintf("/transactions/%s", randomID), nil)

	// then
	ts.Require().Equal(http.StatusNotFound, res.Code)
	requireValidJSON(ts.T(), res.Body.Bytes(), InvalidJSONResponsePayload, res.Body.String())
}

func (ts *TransactionHandlerTestSuite) TestDeleteByID_WithErrorOnDeleteByID() {
	// given
	randomID := uuid.NewString()

	ts.repositoryMock.EXPECT().FindByID(randomID).Return(&entities.Transaction{}, nil)
	ts.repositoryMock.EXPECT().DeleteByID(randomID).Return(errorOnMethod("DeleteByID"))

	// when
	res := makeRequest(ts.router, http.MethodDelete, fmt.Sprintf("/transactions/%s", randomID), nil)

	// then
	ts.Require().Equal(http.StatusInternalServerError, res.Code)
	requireValidJSON(ts.T(), res.Body.Bytes(), InvalidJSONResponsePayload, res.Body.String())
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

func requireValidJSON(t *testing.T, data []byte, msgAndArgs ...string) {
	require.True(t, json.Valid(data), msgAndArgs)
}

func errorOnMethod(method string) error {
	return fmt.Errorf("error on %s", method)
}

func makeRequest(router *chi.Mux, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, body)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	return rr
}
