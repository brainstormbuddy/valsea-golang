package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// TestCreateBankAccount tests the creation of a new bank account.
func TestCreateBankAccount(t *testing.T) {
	// Prepare the request body with account details
	reqBody := `{"owner": "John Doe", "initial_balance": 100.0}`
	req, err := http.NewRequest("POST", "/accounts", bytes.NewBufferString(reqBody))
	assert.NoError(t, err)

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createAccount)
	handler.ServeHTTP(rr, req)

	// Assert that the response status code is 201 Created
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Decode the response body into an Account struct
	var account Account
	err = json.NewDecoder(rr.Body).Decode(&account)
	assert.NoError(t, err)
	// Assert that the account details are as expected
	assert.Equal(t, "John Doe", account.Owner)
	assert.Equal(t, 100.0, account.Balance)
}

// TestGetBankAccount tests retrieving a bank account by ID.
func TestGetBankAccount(t *testing.T) {
	// Create an account first
	account := &Account{ID: 1, Owner: "John Doe", Balance: 100.0}
	accounts[1] = account // Store the account in the accounts map

	// Prepare a GET request to retrieve the account
	req, err := http.NewRequest("GET", "/accounts/1", nil)
	assert.NoError(t, err)

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/accounts/{id}", getAccount)
	router.ServeHTTP(rr, req)

	// Assert that the response status code is 200 OK
	assert.Equal(t, http.StatusOK, rr.Code)

	// Decode the response body into an Account struct
	var retrievedAccount Account
	err = json.NewDecoder(rr.Body).Decode(&retrievedAccount)
	assert.NoError(t, err)

	// Assert that the retrieved account matches the original account
	assert.Equal(t, account.ID, retrievedAccount.ID)
	assert.Equal(t, account.Owner, retrievedAccount.Owner)
	assert.Equal(t, account.Balance, retrievedAccount.Balance)
}

// TestListBankAccounts tests listing all bank accounts.
func TestListBankAccounts(t *testing.T) {
	// Create some accounts first
	accounts[1] = &Account{ID: 1, Owner: "John Doe", Balance: 100.0}
	accounts[2] = &Account{ID: 2, Owner: "Jane Doe", Balance: 200.0}

	// Prepare a GET request to list all accounts
	req, err := http.NewRequest("GET", "/accounts", nil)
	assert.NoError(t, err)

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(listAccounts)
	handler.ServeHTTP(rr, req)

	// Assert that the response status code is 200 OK
	assert.Equal(t, http.StatusOK, rr.Code)

	var accountList []*Account
	err = json.NewDecoder(rr.Body).Decode(&accountList)
	assert.NoError(t, err)

	// Assert that the number of accounts returned is correct
	assert.Len(t, accountList, 2)

	// Assert that the account owners are as expected
	assert.Equal(t, "John Doe", accountList[0].Owner)
	assert.Equal(t, "Jane Doe", accountList[1].Owner)
}
