package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// TestCreateBankTransaction tests the creation of bank transactions (deposit and withdrawal)
func TestCreateBankTransaction(t *testing.T) {
	// Create an account first
	account := &Account{ID: 1, Owner: "John Doe", Balance: 100.0}
	accounts[1] = account

	// Test deposit transaction
	reqBody := `{"type": "deposit", "amount": 50.0}`
	req, err := http.NewRequest("POST", "/accounts/1/transactions", bytes.NewBufferString(reqBody))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/accounts/{id}/transactions", createTransaction)
	router.ServeHTTP(rr, req)

	// Assert that the response status is 201 Created
	assert.Equal(t, http.StatusCreated, rr.Code)

	var transaction Transaction
	err = json.NewDecoder(rr.Body).Decode(&transaction)
	assert.NoError(t, err)

	// Validate the transaction details
	assert.Equal(t, "deposit", transaction.Type)
	assert.Equal(t, 50.0, transaction.Amount)
	assert.Equal(t, 150.0, accounts[1].Balance)

	// Test withdrawal transaction
	reqBody = `{"type": "withdrawal", "amount": 30.0}`
	req, err = http.NewRequest("POST", "/accounts/1/transactions", bytes.NewBufferString(reqBody))
	assert.NoError(t, err)

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Assert that the response status is 201 Created
	assert.Equal(t, http.StatusCreated, rr.Code)

	err = json.NewDecoder(rr.Body).Decode(&transaction)
	assert.NoError(t, err)

	// Validate the transaction details
	assert.Equal(t, "withdrawal", transaction.Type)
	assert.Equal(t, 30.0, transaction.Amount)
	assert.Equal(t, 120.0, accounts[1].Balance)
}

// TestGetBankTransactions tests retrieving bank transactions for a specific account.
func TestGetBankTransactions(t *testing.T) {
	// Create an account and some transactions first
	account := &Account{ID: 1, Owner: "John Doe", Balance: 100.0}
	accounts[1] = account
	transactions[1] = []Transaction{
		{ID: 1, AccountID: 1, Type: "deposit", Amount: 50.0, Timestamp: time.Now()},
		{ID: 2, AccountID: 1, Type: "withdrawal", Amount: 30.0, Timestamp: time.Now()},
	}

	req, err := http.NewRequest("GET", "/accounts/1/transactions", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/accounts/{id}/transactions", getTransactions)
	router.ServeHTTP(rr, req)

	// Assert that the response status is 200 OK
	assert.Equal(t, http.StatusOK, rr.Code)

	var accountTransactions []Transaction
	err = json.NewDecoder(rr.Body).Decode(&accountTransactions)
	assert.NoError(t, err)

	// Validate the number of transactions and their types
	assert.Len(t, accountTransactions, 2)
	assert.Equal(t, "deposit", accountTransactions[0].Type)
	assert.Equal(t, "withdrawal", accountTransactions[1].Type)
}

// TestTransferBankFunds tests transferring funds between two bank accounts.
func TestTransferBankFunds(t *testing.T) {
	// Reset accounts and transactions before the test
	accounts = make(map[int64]*Account)
	transactions = make(map[int64][]Transaction)

	// Create two accounts first for testing fund
	fromAccount := &Account{ID: 1, Owner: "John Doe", Balance: 100.0}
	toAccount := &Account{ID: 2, Owner: "Jane Doe", Balance: 50.0}
	accounts[1] = fromAccount
	accounts[2] = toAccount

	reqBody := `{"from_account_id": 1, "to_account_id": 2, "amount": 30.0}`
	req, err := http.NewRequest("POST", "/transfer", bytes.NewBufferString(reqBody))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/transfer", transferFunds)
	router.ServeHTTP(rr, req)

	// Assert that the response status is 200 OK
	assert.Equal(t, http.StatusOK, rr.Code)

	// Validate the updated balances after transfer
	assert.Equal(t, 70.0, accounts[1].Balance)
	assert.Equal(t, 80.0, accounts[2].Balance)

	// Check transactions for both accounts
	assert.Len(t, transactions[1], 1)
	assert.Equal(t, "transfer_out", transactions[1][0].Type)
	assert.Equal(t, 30.0, transactions[1][0].Amount)

	assert.Len(t, transactions[2], 1)
	assert.Equal(t, "transfer_in", transactions[2][0].Type)
	assert.Equal(t, 30.0, transactions[2][0].Amount)
}
