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

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err.Error())
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	transactionRepository := repositories.NewTransactionMySqlRepository(db)
	cryptoProvider := providers.NewAesGcm256CryptoProvider(cfg.Cryptography.SecretKey)

	handlers.MountTransactionHandler(r, transactionRepository, cryptoProvider)

	err = http.ListenAndServe(":3000", r)
	if err != nil {
		panic(err)
	}
}
