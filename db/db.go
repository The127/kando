package db

import (
	"database/sql"
	"fmt"
	"kando-backend/config"
	"kando-backend/log"

	_ "github.com/lib/pq"
)

func ConnectToDatabase() *sql.DB {
	log.Logger.Infof("Connecting to database %s via %s:%d",
		config.C.Database.Database,
		config.C.Database.Host,
		config.C.Database.Port)

	connection, err := connectToDatabase(
		config.C.Database.Host,
		config.C.Database.Port,
		config.C.Database.User,
		config.C.Database.Password,
		config.C.Database.Database,
		config.C.Database.SslMode)
	if err != nil {
		log.Logger.Fatalf("Failed to connect to the database: %v", err)
	}

	return connection
}

func ConnectToTestDatabase(database string) (*sql.DB, error) {
	return connectToDatabase(
		"localhost",
		5784,
		"user",
		"password",
		database,
		"disable")
}

func connectToDatabase(host string,
	port int,
	user string,
	password string,
	database string,
	sslMode string) (*sql.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=%s",
		host,
		port,
		user,
		password,
		sslMode)

	if database != "" {
		connectionString += fmt.Sprintf(" dbname=%s", database)
	}

	return sql.Open("postgres", connectionString)
}
