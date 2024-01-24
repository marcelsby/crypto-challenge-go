package handlers

import (
	"crypto-challenge/database"
	"crypto-challenge/models"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type TransactionHandler struct {
	repository *database.TransactionInMemoryRepository
}

func NewTransactionHandler(repository *database.TransactionInMemoryRepository) *TransactionHandler {
	return &TransactionHandler{repository: repository}
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	newPerson := models.Transaction{
		ID:              uuid.NewString(),
		UserDocument:    "31377680045",
		CreditCardToken: "5466681299600307",
		Value:           299.8,
	}

	h.repository.Create(&newPerson)

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newPerson)
	w.WriteHeader(http.StatusCreated)
}

// func (h *PersonHandler) FindByID(w http.ResponseWriter, r *http.Request) *models.Person {
// 	result = h.repository.FindByID(r.)

// 	if result == nil {
// 		w.WriteHeader(404)
// 		return
// 	}

// 	w.Write()
// }

func (h *TransactionHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	transactions := h.repository.FindAll()

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}
