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

	ts.Require().True(ts.cryptoProviderMock.AssertCalled(ts.T(), "Encrypt", mock.AnythingOfType("*entities.Transaction")))
	ts.Require().True(ts.cryptoProviderMock.AssertNumberOfCalls(ts.T(), "Encrypt", 1))

	ts.Require().True(ts.repositoryMock.AssertCalled(ts.T(), "Create", mock.AnythingOfType("*entities.Transaction")))
	ts.Require().True(ts.repositoryMock.AssertNumberOfCalls(ts.T(), "Create", 1))
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
	ts.Require().True(ts.cryptoProviderMock.AssertNotCalled(ts.T(), "Encrypt"))
	ts.Require().True(ts.repositoryMock.AssertNotCalled(ts.T(), "Create"))
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
		Return(errors.New("error on encryption")).Once()

	// when
	ts.router.ServeHTTP(res, req)

	// then
	ts.Require().Equal(http.StatusInternalServerError, res.Code)
	ts.Require().Equal("application/json", res.Header().Get("Content-Type"))
	ts.Require().NotEmpty(res.Body.Bytes())
	ts.Require().True(json.Valid(res.Body.Bytes()), "invalid JSON response. Received:", res.Body.String())

	ts.Assert().True(ts.repositoryMock.AssertNotCalled(ts.T(), "Create"))
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

	ts.Require().True(ts.cryptoProviderMock.AssertCalled(ts.T(), "Encrypt", mock.AnythingOfType("*entities.Transaction")))
	ts.Require().True(ts.cryptoProviderMock.AssertNumberOfCalls(ts.T(), "Encrypt", 1))

	ts.Require().NotEmpty(res.Body.Bytes())
	ts.Assert().True(json.Valid(res.Body.Bytes()), "invalid JSON response. Received:", res.Body.String())
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

func generateRandomTransaction(withID bool) entities.Transaction {
	fakeUserDocuments := []string{"50277613433", "19318615400", "43872034856", "25694674300", "56214093854", "01927386406", "89673401520", "73619405823", "40198237610", "58327490120"}
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

	return entities.Transaction{
		ID:              id,
		UserDocument:    randomUserDocument,
		CreditCardToken: randomCreditCardToken,
		Value:           randomValue,
	}
}
