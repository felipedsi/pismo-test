package repository

import "github.com/felipedsi/pismo-test/model"

type AccountRepository interface {
	CreateAccount(account model.Account) (*model.Account, error)
	FindAccount(accountId uint64) (*model.Account, error)
}
