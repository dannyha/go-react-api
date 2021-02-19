package main

import "time"

// Database credentials
var envDbCredentials string = "host=localhost user=postgres dbname=koho port=5432 sslmode=disable TimeZone=America/Toronto"

// Offset of local time to UTC, this value is based off TORONTO time
var envTimeOffset time.Duration = -5

// Localhost Port
var envLocalhostPort string = ":8000"
