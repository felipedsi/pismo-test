package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/felipedsi/pismo-test/model"
	"github.com/felipedsi/pismo-test/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type AccountHandler struct {
	repository repository.AccountRepository
}

func NewAccountHandler(repository repository.AccountRepository) *AccountHandler {
	return &AccountHandler{
		repository: repository,
	}
}

func (c *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	payload := &AccountPayload{}

	err := render.Bind(r, payload)

	if (err != nil) || (payload.DocumentNumber <= 0) {
		render.Render(w, r, errorInvalidRequest(err, "The document_number must be a valid positive integer."))
		return
	}

	account, err := c.repository.CreateAccount(model.Account{
		DocumentNumber: payload.DocumentNumber,
	})

	if err != nil {
		render.Render(w, r, errorInvalidRequest(err, "An error occurred when creating the account."))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, account)
}

func (c *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	accountIdParam := chi.URLParam(r, "accountId")
	accountId, err := strconv.ParseUint(accountIdParam, 10, 64)

	if (err != nil) || (accountId <= 0) {
		render.Render(w, r, errorInvalidRequest(err, "The account_id must be a valid positive integer."))
		return
	}

	account, err := c.repository.FindAccount(accountId)

	if err != nil {
		if err == sql.ErrNoRows {
			render.Render(w, r, errorNotFound(err, "No account found for the provided account ID."))
			return
		}

		render.Render(w, r, errorInvalidRequest(err, "An error occurred when fetching the account from the database."))
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, account)
}

type AccountPayload struct {
	AccountId      uint64 `json:"account_id,omitempty"`
	DocumentNumber uint64 `json:"document_number"`
}

func (a *AccountPayload) Bind(r *http.Request) error {
	return nil
}

func (a *AccountPayload) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
