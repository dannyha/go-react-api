package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/jinzhu/now"
	"gorm.io/gorm"
)

// Generic Error handler
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Post Transaction API
func (database *DB) addTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		return
	}

	var post TransactionPost
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("error with post")
	}
	result := database.validateTransaction(post)
	response := TransactionResponse{ID: post.ID, CustomerID: post.CustomerID, Accepted: result}
	json.NewEncoder(w).Encode(response)
	fmt.Println(response)
}

// Post Transaction API
func (database *DB) readInput() {
	input, err := os.Open("./input.txt")
	check(err)
	output, err := os.Create("./output.txt")
	check(err)
	writer := bufio.NewWriter(output)

	read := bufio.NewScanner(input)
	for read.Scan() {
		var post TransactionPost
		if err := json.Unmarshal(read.Bytes(), &post); err != nil {
			check(err)
		}
		result := database.validateTransaction(post)
		response := TransactionResponse{ID: post.ID, CustomerID: post.CustomerID, Accepted: result}
		jsonString, _ := json.Marshal(response)
		writer.WriteString(string(jsonString) + "\n")
	}
	if read.Err() != nil {
		// handle scan error
		fmt.Println("handle scan error")
	}
	writer.Flush()
}

// Query Transaction for based on customer and time and returns sum and count from database
func (database *DB) queryTransactions(customer string, time time.Time) (float64, int64) {
	var sum float64
	var count int64
	if result := database.db.Table("transactions").Where("customer_id = ?", customer).Where("time >= ?", time).Count(&count).Select("sum(amount) as sum").Scan(&sum); result.Error != nil {
		fmt.Println("Transactions for the current parameter are unavailable")
	}
	return sum, count
}

// Create Transaction in database
func (database *DB) createTransaction(record Transaction) bool {
	if result := database.db.Create(&record); result.Error != nil {
		fmt.Println("Unable to create transaction")
		return false
	}
	fmt.Println("Transaction created successfully")
	return true
}

//Validates the transaction against the rules
func (database *DB) validateTransaction(trans TransactionPost) bool {
	_, offset := loadConfiguration()
	var localDiffUTC time.Duration = offset //Based off Toronto time
	var dailyMaxSum float64 = 5000.00
	var dailyMaxCount int64 = 3
	var weeklyMaxSum float64 = 20000.00
	var weeklyMaxCount int64 = 21

	amount := convertStringAmountToFloat(trans.Amount)
	current := Transaction{ID: trans.ID, CustomerID: trans.CustomerID, Amount: amount, Time: trans.Time}
	now.WeekStartDay = time.Monday
	currentDay := now.BeginningOfDay().UTC().Add(time.Hour * localDiffUTC)
	currentWeek := now.BeginningOfWeek().UTC().Add(time.Hour * localDiffUTC)
	currentDaySum, currentDayCount := database.queryTransactions(trans.CustomerID, currentDay)
	currentWeekSum, currentWeekCount := database.queryTransactions(trans.CustomerID, currentWeek)

	if currentDayCount < dailyMaxCount &&
		currentWeekCount < weeklyMaxCount &&
		(currentDaySum+amount) <= dailyMaxSum &&
		(currentWeekSum+amount) <= weeklyMaxSum {
		return database.createTransaction(current)
	}
	return false
}

//Convert string to float64
func convertStringAmountToFloat(a string) float64 {
	amount := a
	re, err := regexp.Compile(`[$,]`)
	check(err)
	amount = re.ReplaceAllString(amount, "")
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		fmt.Println("Invalid Float")
	}
	return amountFloat
}

// CreateDatabase - Used for the Database
func CreateDatabase(dataB *gorm.DB) Methods {
	return &DB{
		db: dataB,
	}
}
