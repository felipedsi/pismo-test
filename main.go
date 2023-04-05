package main

import (
	"os"

	_ "github.com/lib/pq"

	"database/sql"
	"log"
	"net/http"

	"github.com/felipedsi/pismo-test/handler"
	"github.com/felipedsi/pismo-test/repository/adapter"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	connStr := os.Getenv("POSTGRESQL_URL")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	accountRepositoryPostgres := adapter.NewAccountRepositoryPostgres(db)
	transactionRepositoryPostgres := adapter.NewTransactionRepositoryPostgres(db)

	accountHandler := handler.NewAccountHandler(accountRepositoryPostgres)
	transactionHandler := handler.NewTransactionHandler(transactionRepositoryPostgres)

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Post("/accounts", accountHandler.CreateAccount)
	r.Get("/accounts/{accountId}", accountHandler.GetAccount)
	r.Post("/transactions", transactionHandler.CreateTransaction)

	http.ListenAndServe(":3000", r)
}
