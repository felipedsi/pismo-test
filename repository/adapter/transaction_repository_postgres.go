package adapter

import (
	"database/sql"
	"log"

	"github.com/felipedsi/pismo-test/model"
)

type TransactionRepositoryPostgres struct {
	db *sql.DB
}

func NewTransactionRepositoryPostgres(db *sql.DB) *TransactionRepositoryPostgres {
	return &TransactionRepositoryPostgres{
		db: db,
	}
}

func (t *TransactionRepositoryPostgres) CreateTransaction(transaction model.Transaction) (*model.Transaction, error) {
	query := "INSERT INTO transactions (account_id, operation_type_id, amount) VALUES ($1, $2, $3) RETURNING transaction_id"

	err := t.db.QueryRow(
		query,
		transaction.AccountId,
		transaction.OperationTypeId,
		transaction.Amount).Scan(&transaction.TransactionId)

	if err != nil {
		log.Printf("TransactionRepositoryPostgres#CreateTransaction: Database query (%s) failed: %s", query, err)

		return nil, err
	}

	return &transaction, nil
}
