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
	router         *chi.Mux
	repository     repositories.TransactionRepository
	cryptoProvider providers.CryptoProvider
}

func MountTransactionHandler(router *chi.Mux, repository repositories.TransactionRepository, cryptoProvider providers.CryptoProvider) {
	handler := &TransactionHandler{router, repository, cryptoProvider}
	mountHandlers(handler)
}

func (h *TransactionHandler) Create() {
	h.router.Post("/transactions", func(w http.ResponseWriter, r *http.Request) {
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
	})
}

func (h *TransactionHandler) FindByID() {
	h.router.Get("/transactions/{id}", func(w http.ResponseWriter, r *http.Request) {
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
	})
}

func (h *TransactionHandler) FindAll() {
	h.router.Get("/transactions", func(w http.ResponseWriter, r *http.Request) {
		transactions, err := h.repository.FindAll()
		if err != nil {
			h.setupInternalServerErrorResponse(w)
			return
		}

		w.Header().Add("Content-Type", "application/json")

		for _, transaction := range transactions {
			err := h.decryptTransaction(transaction)
			if err != nil {
				h.setupInternalServerErrorResponse(w)
				return
			}
		}

		json.NewEncoder(w).Encode(transactions)
	})
}

func (h *TransactionHandler) UpdateByID() {
	h.router.Put("/transactions/{id}", func(w http.ResponseWriter, r *http.Request) {
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

		err = h.repository.UpdateByID(&updatedTransaction)
		if err != nil {
			h.setupInternalServerErrorResponse(w)
			return
		}
	})
}

func (h *TransactionHandler) DeleteByID() {
	h.router.Delete("/transactions/{id}", func(w http.ResponseWriter, r *http.Request) {
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

		err := h.repository.DeleteByID(idToBeDeleted)
		if err != nil {
			h.setupInternalServerErrorResponse(w)
			return
		}
	})
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
