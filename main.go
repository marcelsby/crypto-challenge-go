package main

import (
	"crypto-challenge/database"
	"crypto-challenge/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	transactionRepository := database.NewTransactionInMemoryRepository()
	transactionHandler := handlers.NewTransactionHandler(transactionRepository)

	r.Use(middleware.Logger)

	r.Post("/transactions", transactionHandler.Create)
	r.Get("/transactions", transactionHandler.FindAll)

	err := http.ListenAndServe(":3000", r)

	if err != nil {
		panic(err)
	}
}
