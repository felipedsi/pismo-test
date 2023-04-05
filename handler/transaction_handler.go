package handler

import (
	"net/http"
	"strings"

	"github.com/felipedsi/pismo-test/model"
	"github.com/felipedsi/pismo-test/repository"
	"github.com/go-chi/render"
)

type TransactionHandler struct {
	repository repository.TransactionRepository
}

func NewTransactionHandler(repository repository.TransactionRepository) *TransactionHandler {
	return &TransactionHandler{
		repository: repository,
	}
}

func (c *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	payload := &TransactionPayload{}

	err := render.Bind(r, payload)

	if err != nil {
		render.Render(w, r, errorInvalidRequest(err, "The account_id and operation_type_id must be valid positive integers. The amount must be a valid decimal."))
		return
	}

	payloadErrors := validatePayload(payload)

	if len(payloadErrors) > 0 {
		render.Render(w, r, errorInvalidRequest(err, strings.Join(payloadErrors, " ")))
		return
	}

	transaction, err := c.repository.CreateTransaction(model.Transaction{
		AccountId:       payload.AccountId,
		OperationTypeId: payload.OperationTypeId,
		Amount:          payload.Amount,
	})

	if err != nil {
		render.Render(w, r, errorInvalidRequest(err, "The provided account does not exist."))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, transaction)
}

func validatePayload(payload *TransactionPayload) []string {
	var errors []string

	if payload.AccountId <= 0 {
		errors = append(errors, "The account_id must be a valid positive integer.")
	}

	if !model.ValidateOperationType(payload.OperationTypeId) {
		errors = append(errors, "The operation_type_id must be one of the following valid values: 1, 2, 3, 4")
	}

	if !model.ValidateOperationTypeAmount(payload.OperationTypeId, payload.Amount) {
		errors = append(errors, "Purchases and withdraw operations must have a negative amount. Payment operations must have a positive amount.")
	}

	return errors
}

type TransactionPayload struct {
	AccountId       uint64  `json:"account_id"`
	OperationTypeId uint32  `json:"operation_type_id"`
	Amount          float32 `json:"amount"`
}

func (t *TransactionPayload) Bind(r *http.Request) error {
	return nil
}

func (t *TransactionPayload) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
