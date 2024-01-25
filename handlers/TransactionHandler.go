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
	var newPerson models.Transaction

	err := json.NewDecoder(r.Body).Decode(&newPerson)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	newPerson.ID = uuid.NewString()

	h.repository.Create(&newPerson)

	w.WriteHeader(http.StatusCreated)
}

func (h *TransactionHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	idToSearchBy := chi.URLParam(r, "id")

	result := h.repository.FindByID(idToSearchBy)

	if result == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *TransactionHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	transactions := h.repository.FindAll()

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}
