package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Create a transaction
func createTransaction(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64) // Parse the account ID from the URL parameters

	mutex.Lock()
	account, exists := accounts[id]
	mutex.Unlock()

	if !exists {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	// Define the request structure for transaction details
	var req struct {
		Type   string  `json:"type"`
		Amount float64 `json:"amount"`
	}

	// Decode the JSON request body into the req struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	// Check for sufficient funds in case of withdrawal
	if req.Type == "withdrawal" && account.Balance < req.Amount {
		http.Error(w, "Insufficient funds", http.StatusBadRequest)
		return
	}

	// Create a new transaction
	transactionIDSeq++
	transaction := Transaction{
		ID:        transactionIDSeq, // Unique identifier for the transaction
		AccountID: account.ID,       // ID of the account associated with the transaction
		Type:      req.Type,         // Type of transaction (e.g., deposit, withdrawal, transfer)
		Amount:    req.Amount,       // Amount of money involved in the transaction
		Timestamp: time.Now(),       // Record the current timestamp
	}

	// Update the account balance based on the transaction type
	if req.Type == "deposit" {
		account.Balance += req.Amount
	} else if req.Type == "withdrawal" {
		account.Balance -= req.Amount
	} else {
		http.Error(w, "Invalid transaction type", http.StatusBadRequest)
		return
	}

	// Append the transaction to the account's transaction history
	transactions[account.ID] = append(transactions[account.ID], transaction)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}

// Retrieve transactions for an account
func getTransactions(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64) // Parse the account ID from the URL parameters

	mutex.Lock()
	accountTransactions, exists := transactions[id]
	mutex.Unlock()

	if !exists {
		http.Error(w, "No transactions found for this account", http.StatusNotFound)
		return
	}

	// Return the list of transactions as JSON
	json.NewEncoder(w).Encode(accountTransactions)
}

// Transfer funds between accounts
func transferFunds(w http.ResponseWriter, r *http.Request) {
	// Define the request structure for fund transfer details
	var req struct {
		FromAccountID int64   `json:"from_account_id"`
		ToAccountID   int64   `json:"to_account_id"`
		Amount        float64 `json:"amount"`
	}

	// Decode the request body into the req struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	fromAccount, fromExists := accounts[req.FromAccountID]
	toAccount, toExists := accounts[req.ToAccountID]
	mutex.Unlock()

	if !fromExists || !toExists {
		http.Error(w, "Invalid account IDs", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	// Check for sufficient funds in the from account
	if fromAccount.Balance < req.Amount {
		http.Error(w, "Insufficient funds", http.StatusBadRequest)
		return
	}

	// Perform the fund transfer
	fromAccount.Balance -= req.Amount // Deduct from the sender's account
	toAccount.Balance += req.Amount   // Add to the receiver's account

	// Create a transaction record for the transfer out
	transactionIDSeq++
	transaction := Transaction{
		ID:        transactionIDSeq,
		AccountID: fromAccount.ID,
		Type:      "transfer_out", // Indicate this is a transfer out
		Amount:    req.Amount,
		Timestamp: time.Now(),
	}
	transactions[fromAccount.ID] = append(transactions[fromAccount.ID], transaction)

	// Create a transaction record for the transfer in
	transactionIDSeq++
	transaction = Transaction{
		ID:        transactionIDSeq,
		AccountID: toAccount.ID,
		Type:      "transfer_in", // Indicate this is a transfer in
		Amount:    req.Amount,
		Timestamp: time.Now(),
	}
	transactions[toAccount.ID] = append(transactions[toAccount.ID], transaction)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Transfer successful")
}
