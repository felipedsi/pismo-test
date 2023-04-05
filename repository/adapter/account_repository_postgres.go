package adapter

import (
	"database/sql"

	"github.com/felipedsi/pismo-test/model"
)

type AccountRepositoryPostgres struct {
	db *sql.DB
}

func NewAccountRepositoryPostgres(db *sql.DB) *AccountRepositoryPostgres {
	return &AccountRepositoryPostgres{
		db: db,
	}
}

func (a *AccountRepositoryPostgres) CreateAccount(account model.Account) (*model.Account, error) {
	query := "INSERT INTO accounts (document_number) VALUES ($1) RETURNING account_id"

	err := a.db.QueryRow(query, account.DocumentNumber).Scan(&account.AccountId)

	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (a *AccountRepositoryPostgres) FindAccount(accountId uint64) (*model.Account, error) {
	account := model.Account{}

	query := "SELECT account_id, document_number FROM accounts WHERE account_id=$1 LIMIT 1"

	result := a.db.QueryRow(query, accountId)

	err := result.Scan(&account.AccountId, &account.DocumentNumber)

	if err != nil {
		return nil, err
	}

	return &account, nil
}
