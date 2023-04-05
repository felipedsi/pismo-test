package model

import "net/http"

type Account struct {
	AccountId      uint64 `json:"account_id,omitempty"`
	DocumentNumber uint64 `json:"document_number"`
}

func (a Account) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
