package main

import "time"

func loadConfiguration() (string, time.Duration) {
	dbCredentials := "host=localhost user=postgres dbname=koho port=5432 sslmode=disable TimeZone=America/Toronto"
	var timeOffset time.Duration = -5 //Offset of local time to UTC, this value is based off TORONTO time
	return dbCredentials, timeOffset
}
