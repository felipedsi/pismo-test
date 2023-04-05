package repository

import "github.com/felipedsi/pismo-test/model"

type TransactionRepository interface {
	CreateTransaction(model.Transaction) (*model.Transaction, error)
}
