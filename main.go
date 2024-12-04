package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Route for creating a new account
	r.HandleFunc("/accounts", createAccount).Methods("POST")

	// Route for retrieving a specific account by ID
	r.HandleFunc("/accounts/{id}", getAccount).Methods("GET")

	// Route for listing all accounts
	r.HandleFunc("/accounts", listAccounts).Methods("GET")

	// Route for creating a transaction for a specific account
	r.HandleFunc("/accounts/{id}/transactions", createTransaction).Methods("POST")

	// Route for retrieving transactions for a specific account
	r.HandleFunc("/accounts/{id}/transactions", getTransactions).Methods("GET")

	// Route for transferring funds between accounts
	r.HandleFunc("/transfer", transferFunds).Methods("POST")

	// Start the server on port 8080 and log any fatal errors
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
