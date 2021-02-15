# Go Project Break Down #

-------------------------------------------------------------------------------------------------------

# Requirements

GO 						- https://golang.org/dl/
Postgres DB 			- https://www.postgresql.org/download/

-------------------------------------------------------------------------------------------------------

# Installation

go get -u gorm.io/driver/postgres
go get -u github.com/gorilla/mux
go get -u github.com/jinzhu/now
go get -u gorm.io/gorm
go get -u github.com/DATA-DOG/go-sqlmock
go get -u github.com/stretchr/testify/assert
go get -u github.com/stretchr/testify/require
go get -u github.com/stretchr/testify/suite

-------------------------------------------------------------------------------------------------------

# Project Structure

confiurations.go 	- Config for PostgreSQL credentials and TimeOffset of local to UTC time
models.go 			- Models
methods.go 			- Methods
methods_test.go 	- Testing of Methods
main.go 			- Main
input.txt 			- Input file
output.txt 			- Output file

-------------------------------------------------------------------------------------------------------

# Usage

1. Configure values for "dbCredentials" and "timeOffset" in configurations.go
2. run command: go build
3. Modify inputs.txt
4. run executable: server.exe
5. output.txt file gets generated
6. (BONUS) POST API available at localhost:8000/transaction
	- Client available in client directory


-------------------------------------------------------------------------------------------------------

# Assumptions

Since the transactions are in real time. The input date should be the current date.


