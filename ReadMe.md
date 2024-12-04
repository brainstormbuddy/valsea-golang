# Bank API

This project is a simple RESTful API for managing bank accounts and transactions using Go and the Gorilla Mux router.

## Features

- Create accounts
- Retrieve account details
- List all accounts
- Create transactions (deposits and withdrawals)
- Retrieve transactions for an account
- Transfer funds between accounts

## Requirements

- Go 1.23.4 or higher
- Gorilla Mux package

## Installation

Install the required dependencies:

   ```bash
   go mod tidy
   ```

## Running the Application

To run the application, execute the following command:

```bash
go run main.go models.go account.go transaction.go
```

The server will start and listen on port 8080. You should see the following message in the console:

`Server is running on port 8080`

## Running Test Cases

To run the test cases, execute the following command:

   ```bash
   go test -v
   ```
