package main

import (
	"sync"
	"time"
)

type Account struct {
	ID      int64   `json:"id"`
	Owner   string  `json:"owner"`
	Balance float64 `json:"balance"`
}

type Transaction struct {
	ID        int64     `json:"id"`
	AccountID int64     `json:"account_id"`
	Type      string    `json:"type"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

var (
	accounts         = make(map[int64]*Account)
	transactions     = make(map[int64][]Transaction)
	accountIDSeq     int64
	transactionIDSeq int64
	mutex            sync.Mutex
)
