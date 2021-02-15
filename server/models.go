package main

import (
	"net/http"
	"time"

	"gorm.io/gorm"
)

// DB - Database struct (Model)
type DB struct {
	db *gorm.DB
}

// Transaction struct (Model)
type Transaction struct {
	ID         string    `json:"id"`
	CustomerID string    `json:"customer_id"`
	Amount     float64   `json:"load_amount" gorm:"type:numeric(15,2)"`
	Time       time.Time `json:"time"`
}

var transaction []Transaction

// TransactionPost - Transaction Post struct (Model)
type TransactionPost struct {
	ID         string    `json:"id"`
	CustomerID string    `json:"customer_id"`
	Amount     string    `json:"load_amount,omitempty" gorm:"type:numeric(15,2)"`
	Time       time.Time `json:"time"`
}

// TransactionResponse - Transaction Response struct (Model)
type TransactionResponse struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
	Accepted   bool   `json:"accepted"`
}

// Methods (interface )
type Methods interface {
	addTransaction(w http.ResponseWriter, r *http.Request)
	queryTransactions(customer string, time time.Time) (float64, int64)
	createTransaction(record Transaction) bool
	validateTransaction(trans TransactionPost) bool
}
