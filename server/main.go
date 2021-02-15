package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Main function
func main() {
	dbCreds, _ := loadConfiguration()
	db, dberr := gorm.Open(postgres.Open(dbCreds), &gorm.Config{})
	database := &DB{db: db}

	if dberr != nil {
		fmt.Println("Connection failed")
	}

	// Auto Migrate
	db.AutoMigrate(&Transaction{})

	// Init router
	r := mux.NewRouter()

	// Route handles & endpoints
	r.HandleFunc("/transaction", database.addTransaction).Methods("POST", "OPTIONS")

	// Read input.txt file
	database.readInput()

	// Start server
	log.Fatal(http.ListenAndServe(":8000", r))
}
