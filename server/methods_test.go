package main

import (
	"database/sql"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/now"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock

	methods         Methods
	transaction     *Transaction
	transactionPost *TransactionPost
	// transactionResponse *TransactionResponse
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	require.NoError(s.T(), err)

	s.methods = CreateDatabase(s.DB)
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func TestConvertStringAmountToFloat(t *testing.T) {
	amount := convertStringAmountToFloat("$123456.78")
	if amount != 123456.78 {
		t.Errorf("Amount was incorrect, got: %f, want: %f.", amount, 123456.78)
	}
	amount = convertStringAmountToFloat("$$456.99")
	if amount != 456.99 {
		t.Errorf("Amount was incorrect, got: %f, want: %f.", amount, 456.99)
	}
}

func (s *Suite) TestQueryTransactions() {

	var (
		customer   = "123"
		datestring = "2021-02-12T08:10:56Z"
	)

	convertedTime, err := time.Parse(time.RFC3339, datestring)
	if err != nil {
		fmt.Println(err)
	}

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(1) FROM "transactions" WHERE customer_id = $1 AND time >= $2`)).
		WithArgs(customer, convertedTime).
		WillReturnRows(
			sqlmock.NewRows([]string{"count"}).AddRow(1))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT sum(amount) as sum FROM "transactions" WHERE customer_id = $1 AND time >= $2`)).
		WithArgs(customer, convertedTime).
		WillReturnRows(
			sqlmock.NewRows([]string{"sum"}).AddRow(1234.56))

	sum, count := s.methods.queryTransactions(customer, convertedTime)

	assert.Equal(s.T(), sum, float64(1234.56))
	assert.Equal(s.T(), count, int64(1))
}

func (s *Suite) TestValidateTransaction_ValidTransaction() {
	var localDiffUTC time.Duration = envTimeOffset
	now.WeekStartDay = time.Monday
	currentDay := now.BeginningOfDay().UTC().Add(time.Hour * localDiffUTC)
	currentWeek := now.BeginningOfWeek().UTC().Add(time.Hour * localDiffUTC)

	var (
		id       = "99911"
		customer = "123"
		amount1  = "$1000.00"
		amount2  = float64(1000.00)
		time1    = currentDay
		time2    = currentWeek
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(1) FROM "transactions" WHERE customer_id = $1 AND time >= $2`)).
		WithArgs(customer, time1).
		WillReturnRows(
			sqlmock.NewRows([]string{"count"}).AddRow(1))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT sum(amount) as sum FROM "transactions" WHERE customer_id = $1 AND time >= $2`)).
		WithArgs(customer, time1).
		WillReturnRows(
			sqlmock.NewRows([]string{"sum"}).AddRow(1234.56))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(1) FROM "transactions" WHERE customer_id = $1 AND time >= $2`)).
		WithArgs(customer, time2).
		WillReturnRows(
			sqlmock.NewRows([]string{"count"}).AddRow(1))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT sum(amount) as sum FROM "transactions" WHERE customer_id = $1 AND time >= $2`)).
		WithArgs(customer, time2).
		WillReturnRows(
			sqlmock.NewRows([]string{"sum"}).AddRow(1234.56))

	s.mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO "transactions" ("id","customer_id","amount","time") VALUES ($1,$2,$3,$4)`)).
		WithArgs(id, customer, amount2, time1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	trans := TransactionPost{ID: id, CustomerID: customer, Amount: amount1, Time: time1}
	res := s.methods.validateTransaction(trans)

	assert.Equal(s.T(), res, true)
}

func (s *Suite) TestValidateTransaction_FailedDailyMax() {
	var localDiffUTC time.Duration = envTimeOffset
	now.WeekStartDay = time.Monday
	currentDay := now.BeginningOfDay().UTC().Add(time.Hour * localDiffUTC)
	currentWeek := now.BeginningOfWeek().UTC().Add(time.Hour * localDiffUTC)

	var (
		id       = "99911"
		customer = "123"
		amount1  = "$1000.00"
		time1    = currentDay
		time2    = currentWeek
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(1) FROM "transactions" WHERE customer_id = $1 AND time >= $2`)).
		WithArgs(customer, time1).
		WillReturnRows(
			sqlmock.NewRows([]string{"count"}).AddRow(3))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT sum(amount) as sum FROM "transactions" WHERE customer_id = $1 AND time >= $2`)).
		WithArgs(customer, time1).
		WillReturnRows(
			sqlmock.NewRows([]string{"sum"}).AddRow(1234.56))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(1) FROM "transactions" WHERE customer_id = $1 AND time >= $2`)).
		WithArgs(customer, time2).
		WillReturnRows(
			sqlmock.NewRows([]string{"count"}).AddRow(1))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT sum(amount) as sum FROM "transactions" WHERE customer_id = $1 AND time >= $2`)).
		WithArgs(customer, time2).
		WillReturnRows(
			sqlmock.NewRows([]string{"sum"}).AddRow(1234.56))

	trans := TransactionPost{ID: id, CustomerID: customer, Amount: amount1, Time: time1}
	res := s.methods.validateTransaction(trans)

	assert.Equal(s.T(), res, false)
}

func (s *Suite) TestValidateTransaction_FailedDailyAmount() {
	var localDiffUTC time.Duration = envTimeOffset
	now.WeekStartDay = time.Monday
	currentDay := now.BeginningOfDay().UTC().Add(time.Hour * localDiffUTC)
	currentWeek := now.BeginningOfWeek().UTC().Add(time.Hour * localDiffUTC)

	var (
		id       = "99911"
		customer = "123"
		amount1  = "$1000.00"
		time1    = currentDay
		time2    = currentWeek
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(1) FROM "transactions" WHERE customer_id = $1 AND time >= $2`)).
		WithArgs(customer, time1).
		WillReturnRows(
			sqlmock.NewRows([]string{"count"}).AddRow(2))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT sum(amount) as sum FROM "transactions" WHERE customer_id = $1 AND time >= $2`)).
		WithArgs(customer, time1).
		WillReturnRows(
			sqlmock.NewRows([]string{"sum"}).AddRow(4234.56))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(1) FROM "transactions" WHERE customer_id = $1 AND time >= $2`)).
		WithArgs(customer, time2).
		WillReturnRows(
			sqlmock.NewRows([]string{"count"}).AddRow(1))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT sum(amount) as sum FROM "transactions" WHERE customer_id = $1 AND time >= $2`)).
		WithArgs(customer, time2).
		WillReturnRows(
			sqlmock.NewRows([]string{"sum"}).AddRow(1234.56))

	trans := TransactionPost{ID: id, CustomerID: customer, Amount: amount1, Time: time1}
	res := s.methods.validateTransaction(trans)

	assert.Equal(s.T(), res, false)
}

func (s *Suite) TestValidateTransaction_FailedWeeklyAmount() {
	var localDiffUTC time.Duration = envTimeOffset
	now.WeekStartDay = time.Monday
	currentDay := now.BeginningOfDay().UTC().Add(time.Hour * localDiffUTC)
	currentWeek := now.BeginningOfWeek().UTC().Add(time.Hour * localDiffUTC)

	var (
		id       = "99911"
		customer = "123"
		amount1  = "$1000.00"
		time1    = currentDay
		time2    = currentWeek
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(1) FROM "transactions" WHERE customer_id = $1 AND time >= $2`)).
		WithArgs(customer, time1).
		WillReturnRows(
			sqlmock.NewRows([]string{"count"}).AddRow(2))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT sum(amount) as sum FROM "transactions" WHERE customer_id = $1 AND time >= $2`)).
		WithArgs(customer, time1).
		WillReturnRows(
			sqlmock.NewRows([]string{"sum"}).AddRow(1234.56))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(1) FROM "transactions" WHERE customer_id = $1 AND time >= $2`)).
		WithArgs(customer, time2).
		WillReturnRows(
			sqlmock.NewRows([]string{"count"}).AddRow(1))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT sum(amount) as sum FROM "transactions" WHERE customer_id = $1 AND time >= $2`)).
		WithArgs(customer, time2).
		WillReturnRows(
			sqlmock.NewRows([]string{"sum"}).AddRow(19234.56))

	trans := TransactionPost{ID: id, CustomerID: customer, Amount: amount1, Time: time1}
	res := s.methods.validateTransaction(trans)

	assert.Equal(s.T(), res, false)
}
