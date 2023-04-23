package tests

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/stretchr/testify/suite"
	"kando-backend/db"
	"kando-backend/log"
	"strings"
	"sync"
	"testing"
)

var isTemplateDbInitialized = false
var createTemplateDbLock = &sync.Mutex{}

func SetupTestDatabase(t *testing.T) *sql.DB {
	connection, err := db.ConnectToTestDatabase("postgres")
	if err != nil {
		panic(err)
	}

	dbName := strings.ReplaceAll(t.Name(), "/", "__")

	_, err = connection.Exec(`drop database if exists "` + dbName + `";`)
	if err != nil {
		panic(err)
	}

	_, err = connection.Exec(`create database "` + dbName + `" with template "test_db_template" owner "user";`)
	if err != nil {
		panic(err)
	}

	err = connection.Close()
	if err != nil {
		panic(err)
	}

	connection, err = db.ConnectToTestDatabase(dbName)
	if err != nil {
		panic(err)
	}

	return connection
}

func setupTemplateDatabase() {
	createTemplateDbLock.Lock()
	if isTemplateDbInitialized {
		return
	}
	isTemplateDbInitialized = true
	createTemplateDbLock.Unlock()

	connection, err := db.ConnectToTestDatabase("postgres")
	if err != nil {
		panic(err)
	}

	log.Logger.Info("Dropping test_db_template")
	_, err = connection.Exec(`drop database if exists "test_db_template";`)
	if err != nil {
		panic(err)
	}

	log.Logger.Info("Creating test_db_template")
	_, err = connection.Exec(`create database "test_db_template" with owner "user";`)
	if err != nil {
		panic(err)
	}

	err = connection.Close()
	if err != nil {
		panic(err)
	}

	log.Logger.Info("Connecting to test_db_template")
	connection, err = db.ConnectToTestDatabase("test_db_template")
	if err != nil {
		panic(err)
	}

	log.Logger.Info("Migrating test_db_template")
	db.MigrateDatabase(connection)

	err = connection.Close()
	if err != nil {
		panic(err)
	}
}

type DbTestSuite struct {
	suite.Suite
	dbConn *sql.DB
}

func (s *DbTestSuite) DbConn() *sql.DB {
	return s.dbConn
}

func (s *DbTestSuite) SetupSuite() {
	setupTemplateDatabase()
}

func (s *DbTestSuite) SetupTest() {
	log.Logger.Debugf("Setting up test database for %s", s.T().Name())
	s.dbConn = SetupTestDatabase(s.T())
}

func (s *DbTestSuite) TearDownTest() {
	if s.dbConn == nil {
		return
	}
	err := s.dbConn.Close()
	if err != nil {
		panic(err)
	}
}

func (s *DbTestSuite) VerifyRow(tableName string, columns map[string]any) {
	var columnNames []string
	for key := range columns {
		columnNames = append(columnNames, key)
	}

	selectBuilder := sqlbuilder.Select(columnNames...).From(tableName)

	for key, value := range columns {
		selectBuilder.Where(selectBuilder.Equal(key, value))
	}

	sqlString, args := selectBuilder.Build()

	log.Logger.Debugf("Verifying row in %s with SQL: %s", tableName, sqlString)

	var columnValues []interface{}
	for range columns {
		columnValues = append(columnValues, new(interface{}))
	}

	err := s.dbConn.QueryRow(sqlString, args...).Scan(columnValues...)
	if err != nil {
		panic(err)
	}

	for i, columnValue := range columnValues {
		s.Equal(columns[columnNames[i]], columnValue)
	}
}

func (s *DbTestSuite) InsertRow(tableName string, columns map[string]any) uuid.UUID {
	insertBuilder := sqlbuilder.InsertInto(tableName)

	for key, value := range columns {
		insertBuilder.Cols(key).
			Values(value)
	}

	insertBuilder.SQL("returning id")

	sqlString, args := insertBuilder.Build()

	log.Logger.Debugf("Inserting row into %s with SQL: %s", tableName, sqlString)

	var id uuid.UUID
	err := s.dbConn.QueryRow(sqlString, args...).Scan(&id)
	if err != nil {
		panic(err)
	}

	return id
}
