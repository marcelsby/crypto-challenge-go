package handlers

import (
	"crypto-challenge/database"
	"crypto-challenge/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	repository *database.TransactionInMemoryRepository
}

func NewTransactionHandler(repository *database.TransactionInMemoryRepository) *TransactionHandler {
	return &TransactionHandler{repository: repository}
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var newTransaction models.Transaction

	err := json.NewDecoder(r.Body).Decode(&newTransaction)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	newTransaction.ID = uuid.NewString()

	h.repository.Create(&newTransaction)

	w.WriteHeader(http.StatusCreated)
}

func (h *TransactionHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	idToSearchBy := chi.URLParam(r, "id")

	searchedTransaction := h.repository.FindByID(idToSearchBy)

	w.Header().Add("Content-Type", "application/json")

	if searchedTransaction == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"error":      "Transaction not found with specified ID.",
			"searchedId": idToSearchBy,
		})
		return
	}

	json.NewEncoder(w).Encode(searchedTransaction)
}

func (h *TransactionHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	transactions := h.repository.FindAll()

	w.Header().Add("Content-Type", "application/json")

	json.NewEncoder(w).Encode(transactions)
}

func (h *TransactionHandler) Update(w http.ResponseWriter, r *http.Request) {
	idToUpdate := chi.URLParam(r, "id")

	searchedTransaction := h.repository.FindByID(idToUpdate)

	w.Header().Add("Content-Type", "application/json")

	if searchedTransaction == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"error":      "Transaction not found with specified ID.",
			"searchedId": idToUpdate,
		})
		return
	}

	var updatedTransaction models.Transaction

	json.NewDecoder(r.Body).Decode(&updatedTransaction)

	updatedTransaction.ID = searchedTransaction.ID

	h.repository.Update(searchedTransaction.ID, &updatedTransaction)
}

func (h *TransactionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	var idToBeDeleted = chi.URLParam(r, "id")

	var searchedTransaction = h.repository.FindByID(idToBeDeleted)

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
