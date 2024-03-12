package handlers

import (
	"crypto-challenge/database/repositories"
	"crypto-challenge/entities"
	"crypto-challenge/providers"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	repository     repositories.TransactionRepository
	cryptoProvider *providers.CryptoProvider
}

func NewTransactionHandler(repository repositories.TransactionRepository, cryptoProvider *providers.CryptoProvider) *TransactionHandler {
	return &TransactionHandler{repository, cryptoProvider}
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var newTransaction entities.Transaction

	err := json.NewDecoder(r.Body).Decode(&newTransaction)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	newTransaction.ID = uuid.NewString()

	err = h.encryptTransaction(&newTransaction)
	if err != nil {
		h.setupInternalServerErrorResponse(w)
		return
	}

	err = h.repository.Create(&newTransaction)
	if err != nil {
		h.setupInternalServerErrorResponse(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *TransactionHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	idToSearchBy := chi.URLParam(r, "id")

	searchedTransaction, _ := h.repository.FindByID(idToSearchBy)

	w.Header().Add("Content-Type", "application/json")

	if searchedTransaction == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"error":      "Transaction not found with specified ID.",
			"searchedId": idToSearchBy,
		})
		return
	}

	err := h.decryptTransaction(searchedTransaction)
	if err != nil {
		h.setupInternalServerErrorResponse(w)
		return
	}

	json.NewEncoder(w).Encode(searchedTransaction)
}

func (h *TransactionHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	transactions, _ := h.repository.FindAll()

	w.Header().Add("Content-Type", "application/json")

	for _, transaction := range transactions {
		err := h.decryptTransaction(transaction)
		if err != nil {
			h.setupInternalServerErrorResponse(w)
			return
		}
	}

	json.NewEncoder(w).Encode(transactions)
}

func (h *TransactionHandler) Update(w http.ResponseWriter, r *http.Request) {
	idToUpdate := chi.URLParam(r, "id")

	searchedTransaction, _ := h.repository.FindByID(idToUpdate)

	w.Header().Add("Content-Type", "application/json")

	if searchedTransaction == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"error":      "Transaction not found with specified ID.",
			"searchedId": idToUpdate,
		})
		return
	}

	var updatedTransaction entities.Transaction

	json.NewDecoder(r.Body).Decode(&updatedTransaction)

	updatedTransaction.ID = searchedTransaction.ID

	err := h.encryptTransaction(&updatedTransaction)
	if err != nil {
		h.setupInternalServerErrorResponse(w)
		return
	}

	h.repository.UpdateByID(searchedTransaction.ID, &updatedTransaction)
}

func (h *TransactionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	var idToBeDeleted = chi.URLParam(r, "id")

	searchedTransaction, _ := h.repository.FindByID(idToBeDeleted)

	if searchedTransaction == nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"error":      "Transaction not found with specified ID.",
			"searchedId": idToBeDeleted,
		})
		return
	}

	h.repository.DeleteByID(idToBeDeleted)
}

func (h *TransactionHandler) encryptTransaction(toEncrypt *entities.Transaction) error {
	encryptedUserDocument, err := h.cryptoProvider.Encrypt([]byte(toEncrypt.UserDocument))
	if err != nil {
		return err
	}

	encryptedCreditCardToken, err := h.cryptoProvider.Encrypt([]byte(toEncrypt.CreditCardToken))
	if err != nil {
		return err
	}

	toEncrypt.UserDocument = encryptedUserDocument
	toEncrypt.CreditCardToken = encryptedCreditCardToken

	return nil
}

func (h *TransactionHandler) decryptTransaction(toDecrypt *entities.Transaction) error {
	decryptedUserDocument, err := h.cryptoProvider.Decrypt(toDecrypt.UserDocument)
	if err != nil {
		return err
	}

	decryptedCreditCardToken, err := h.cryptoProvider.Decrypt(toDecrypt.CreditCardToken)
	if err != nil {
		return err
	}

	toDecrypt.UserDocument = string(decryptedUserDocument)
	toDecrypt.CreditCardToken = string(decryptedCreditCardToken)

	return nil
}

func (h *TransactionHandler) setupInternalServerErrorResponse(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]any{
		"error": "An error occurred with the server while processing the request, please try again later.",
	})
}
