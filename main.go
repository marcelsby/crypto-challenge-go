package main

import (
	"crypto-challenge/database"
	"crypto-challenge/handlers"
	"crypto-challenge/providers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	transactionRepository := database.NewTransactionInMemoryRepository()
	cryptoProvider := providers.NewCryptoProvider()
	transactionHandler := handlers.NewTransactionHandler(transactionRepository, cryptoProvider)

	r.Use(middleware.Logger)

	r.Post("/transactions", transactionHandler.Create)
	r.Get("/transactions", transactionHandler.FindAll)
	r.Get("/transactions/{id}", transactionHandler.FindByID)
	r.Put("/transactions/{id}", transactionHandler.Update)
	r.Delete("/transactions/{id}", transactionHandler.Delete)

	err := http.ListenAndServe(":3000", r)

	if err != nil {
		panic(err)
	}
}
