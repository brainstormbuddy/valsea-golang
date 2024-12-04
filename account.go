package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Create a new account
func createAccount(w http.ResponseWriter, r *http.Request) {
	// Define the request structure for account details
	var req struct {
		Owner          string  `json:"owner"`
		InitialBalance float64 `json:"initial_balance"`
	}

	// Decode the JSON request body into the req struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	accountIDSeq++
	account := &Account{
		ID:      accountIDSeq,       // Set the account ID
		Owner:   req.Owner,          // Set the account owner
		Balance: req.InitialBalance, // Set the initial balance
	}
	accounts[account.ID] = account

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

// Retrieve account details
func getAccount(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64) // Parse the account ID from the URL

	mutex.Lock()
	account, exists := accounts[id] // Check if the account exists
	mutex.Unlock()

	if !exists {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(account)
}

// List all accounts
func listAccounts(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	var accountList []*Account // Create a slice to hold the accounts
	for _, account := range accounts {
		accountList = append(accountList, account) // Append each account to the list
	}
	json.NewEncoder(w).Encode(accountList)
}
