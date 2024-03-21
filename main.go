package main

import (
	"crypto-challenge/config"
	"crypto-challenge/database/repositories"
	"crypto-challenge/handlers"
	"crypto-challenge/providers"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg := config.GetAppConfig(".env")

	r := chi.NewRouter()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err.Error())
	}

	transactionRepository := repositories.NewTransactionMySqlRepository(db)
	cryptoProvider := providers.NewAesGcm256CryptoProvider(cfg.Cryptography.SecretKey)
	transactionHandler := handlers.NewTransactionHandler(transactionRepository, cryptoProvider)

	r.Use(middleware.Logger)

	r.Post("/transactions", transactionHandler.Create)
	r.Get("/transactions", transactionHandler.FindAll)
	r.Get("/transactions/{id}", transactionHandler.FindByID)
	r.Put("/transactions/{id}", transactionHandler.Update)
	r.Delete("/transactions/{id}", transactionHandler.Delete)

	err = http.ListenAndServe(":3000", r)
	if err != nil {
		panic(err)
	}
}
