package handlers

import (
	"crypto-challenge/database/repositories"
	"crypto-challenge/entities"
	"crypto-challenge/providers"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	repository                repositories.TransactionRepository
	transactionCryptoProvider providers.TransactionCryptoProvider
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var newTransaction entities.Transaction

	err := json.NewDecoder(r.Body).Decode(&newTransaction)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	newTransaction.ID = uuid.NewString()

	err = h.transactionCryptoProvider.Encrypt(&newTransaction)
	if err != nil {
		setupInternalServerErrorResponse(w)
		return
	}

	err = h.repository.Create(&newTransaction)
	if err != nil {
		setupInternalServerErrorResponse(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *TransactionHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	idToSearchBy := chi.URLParam(r, "id")

	searchedTransaction, err := h.repository.FindByID(idToSearchBy)
	if err != nil {
		setupInternalServerErrorResponse(w)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	if searchedTransaction == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"error":      "Transaction not found with specified ID.",
			"searchedId": idToSearchBy,
		})
		return
	}

	err = h.transactionCryptoProvider.Decrypt(searchedTransaction)
	if err != nil {
		setupInternalServerErrorResponse(w)
		return
	}

	json.NewEncoder(w).Encode(searchedTransaction)
}

func (h *TransactionHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	transactions, err := h.repository.FindAll()
	if err != nil {
		setupInternalServerErrorResponse(w)
		return
	}

	for _, transaction := range transactions {
		err := h.transactionCryptoProvider.Decrypt(transaction)
		if err != nil {
			setupInternalServerErrorResponse(w)
			return
		}
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

func (h *TransactionHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	idToUpdate := chi.URLParam(r, "id")

	searchedTransaction, err := h.repository.FindByID(idToUpdate)
	if err != nil {
		setupInternalServerErrorResponse(w)
		return
	}

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

	err = json.NewDecoder(r.Body).Decode(&updatedTransaction)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	updatedTransaction.ID = searchedTransaction.ID

	err = h.transactionCryptoProvider.Encrypt(&updatedTransaction)
	if err != nil {
		setupInternalServerErrorResponse(w)
		return
	}

	err = h.repository.UpdateByID(&updatedTransaction)
	if err != nil {
		setupInternalServerErrorResponse(w)
		return
	}
}

func (h *TransactionHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	var idToBeDeleted = chi.URLParam(r, "id")

	searchedTransaction, err := h.repository.FindByID(idToBeDeleted)
	if err != nil {
		setupInternalServerErrorResponse(w)
		return
	}

	if searchedTransaction == nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"error":      "Transaction not found with specified ID.",
			"searchedId": idToBeDeleted,
		})
		return
	}

	err = h.repository.DeleteByID(idToBeDeleted)
	if err != nil {
		setupInternalServerErrorResponse(w)
		return
	}
}

func NewTransactionRouter(repository repositories.TransactionRepository, transactionCryptoProvider providers.TransactionCryptoProvider) *chi.Mux {
	r := chi.NewRouter()

	handler := &TransactionHandler{repository, transactionCryptoProvider}

	r.Route("/transactions", func(r chi.Router) {
		r.Post("/", handler.Create)
		r.Get("/", handler.FindAll)
		r.Get("/{id}", handler.FindByID)
		r.Put("/{id}", handler.UpdateByID)
		r.Delete("/{id}", handler.DeleteByID)
	})

	return r
}

func setupInternalServerErrorResponse(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]any{
		"error": "An error occurred with the server while processing the request, please try again later.",
	})
}
