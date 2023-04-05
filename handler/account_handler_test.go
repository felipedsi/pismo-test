package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/felipedsi/pismo-test/model"
)

type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) CreateAccount(account model.Account) (*model.Account, error) {
	args := m.Called(account)
	return args.Get(0).(*model.Account), args.Error(1)
}

func (m *MockAccountRepository) FindAccount(accountId uint64) (*model.Account, error) {
	args := m.Called(accountId)
	return args.Get(0).(*model.Account), args.Error(1)
}

func TestGetAccount(t *testing.T) {
	mockRepo := new(MockAccountRepository)
	accountHandler := NewAccountHandler(mockRepo)

	account := &model.Account{
		AccountId:      123,
		DocumentNumber: 456,
	}

	mockRepo.On("FindAccount", account.AccountId).Return(account, nil)

	req := httptest.NewRequest("GET", "/accounts/123", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("accountId", "123")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	accountHandler.GetAccount(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	response := &AccountPayload{}
	err := json.Unmarshal(w.Body.Bytes(), response)
	assert.NoError(t, err)
	assert.Equal(t, account.AccountId, response.AccountId)
	assert.Equal(t, account.DocumentNumber, response.DocumentNumber)

	mockRepo.AssertExpectations(t)
}

func TestGetAccountFailsWhenInvalidRequest(t *testing.T) {
	var scenarios = []struct {
		accountId          string
		expectedError      error
		expectedStatusCode int
	}{
		{
			"-1",
			errors.New("Error!"),
			http.StatusBadRequest,
		},
		{
			"999",
			sql.ErrNoRows,
			http.StatusNotFound,
		},
		{
			"1",
			errors.New("Database error!"),
			http.StatusBadRequest,
		},
	}

	for _, scenario := range scenarios {
		mockRepo := new(MockAccountRepository)
		handler := &AccountHandler{repository: mockRepo}

		req, _ := http.NewRequest("GET", fmt.Sprintf("/accounts/%s", scenario.accountId), nil)

		rctx := chi.NewRouteContext()

		rctx.URLParams.Add("accountId", scenario.accountId)

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()

		accountId, _ := strconv.ParseUint(scenario.accountId, 10, 64)

		mockRepo.On("FindAccount", accountId).Return(&model.Account{}, scenario.expectedError)

		handler.GetAccount(w, req)

		if w.Code != scenario.expectedStatusCode {
			t.Errorf("Expected status code %d but got %d", scenario.expectedStatusCode, w.Code)
		}
	}
}

func TestCreateAccount(t *testing.T) {
	mockRepo := new(MockAccountRepository)

	payload := `{"document_number": 123456789}`
	req, _ := http.NewRequest("POST", "/accounts", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	expectedAccount := &model.Account{AccountId: 1, DocumentNumber: 123456789}
	mockRepo.On("CreateAccount", mock.AnythingOfType("model.Account")).Return(expectedAccount, nil)

	handler := &AccountHandler{repository: mockRepo}
	handler.CreateAccount(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d but got %d", http.StatusCreated, w.Code)
	}

	expectedResponse := `{"account_id":1,"document_number":123456789}`
	actualResponse := w.Body.String()

	expectedResponseJson := map[string]string{}
	actualResponseJson := map[string]string{}

	json.Unmarshal([]byte(expectedResponse), &expectedResponseJson)
	json.Unmarshal([]byte(actualResponse), &actualResponseJson)

	if !reflect.DeepEqual(expectedResponseJson, actualResponseJson) {
		t.Errorf("Expected response body %s but got %s", expectedResponse, w.Body.String())
	}
}

func TestCreateAccountWhenAccountCreationFails(t *testing.T) {
	mockRepo := new(MockAccountRepository)

	payload := `{"document_number": 123456789}`
	req, _ := http.NewRequest("POST", "/accounts", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockRepo.On("CreateAccount", mock.AnythingOfType("model.Account")).Return(&model.Account{}, errors.New("Error!"))

	handler := &AccountHandler{repository: mockRepo}
	handler.CreateAccount(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d but got %d", http.StatusCreated, w.Code)
	}

	expectedResponse := `{"status": "Invalid request","error":"An error occurred when creating the account."}`
	actualResponse := w.Body.String()

	expectedResponseJson := map[string]string{}
	actualResponseJson := map[string]string{}

	json.Unmarshal([]byte(expectedResponse), &expectedResponseJson)
	json.Unmarshal([]byte(actualResponse), &actualResponseJson)

	if !reflect.DeepEqual(expectedResponseJson, actualResponseJson) {
		t.Errorf("Expected response body %s but got %s", expectedResponse, w.Body.String())
	}
}

func TestCreateAccountFailsWhenInvalidRequest(t *testing.T) {
	var scenarios = []struct {
		payload            string
		expectedResponse   string
		expectedStatusCode int
	}{
		{
			`{"document_number": "invalid"}`,
			`{"status":"Invalid request","error":"The document_number must be a valid positive integer."}`,
			http.StatusBadRequest,
		},
		{
			`{"document_number": ""}`,
			`{"status":"Invalid request","error":"The document_number must be a valid positive integer."}`,
			http.StatusBadRequest,
		},
		{
			`{"document_number": null}`,
			`{"status":"Invalid request","error":"The document_number must be a valid positive integer."}`,
			http.StatusBadRequest,
		},
		{
			`{"document_number": -1}`,
			`{"status":"Invalid request","error":"The document_number must be a valid positive integer."}`,
			http.StatusBadRequest,
		},
	}

	for _, scenario := range scenarios {
		mockRepo := new(MockAccountRepository)

		req, _ := http.NewRequest("POST", "/accounts", strings.NewReader(scenario.payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler := &AccountHandler{repository: mockRepo}
		handler.CreateAccount(w, req)

		if w.Code != scenario.expectedStatusCode {
			t.Errorf("Expected status code %d but got %d", scenario.expectedStatusCode, w.Code)
		}

		expectedResponse := scenario.expectedResponse
		actualResponse := w.Body.String()

		expectedResponseJson := map[string]string{}
		actualResponseJson := map[string]string{}

		json.Unmarshal([]byte(expectedResponse), &expectedResponseJson)
		json.Unmarshal([]byte(actualResponse), &actualResponseJson)

		if !reflect.DeepEqual(expectedResponseJson, actualResponseJson) {
			t.Errorf("Expected response body %s but got %s", expectedResponse, w.Body.String())
		}
	}
}

func TestNewAccountHandler(t *testing.T) {
	repository := &MockAccountRepository{}
	handler := NewAccountHandler(repository)

	if handler.repository != repository {
		t.Errorf("The repository field for the handler wasn't assigned. Expect %s but got %s", repository, handler.repository)
	}
}
