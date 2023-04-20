package tests

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/suite"
	"kando-backend/db"
	"kando-backend/ioc"
	"kando-backend/services"
	"strings"
	"testing"
)

var isTemplateDbInitialized bool = false

func SetupTestDatabase(t *testing.T) *sql.DB {
	connection, err := db.ConnectToTestDatabase("")
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

	return connection
}

func setupTemplateDatabase() {
	if isTemplateDbInitialized {
		return
	}
	isTemplateDbInitialized = true

	connection, err := db.ConnectToTestDatabase("")
	if err != nil {
		panic(err)
	}

	_, err = connection.Exec(`drop database if exists "test_db_template";`)
	if err != nil {
		panic(err)
	}

	_, err = connection.Exec(`create database "test_db_template" with owner "user";`)
	if err != nil {
		panic(err)
	}

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

func TestContext(dbConn *sql.DB) context.Context {
	dpb := ioc.NewDependencyProviderBuilder()

	ioc.AddSingleton(dpb, func(dp *ioc.DependencyProvider) *sql.DB {
		return dbConn
	})

	ioc.AddScoped(dpb, func(dp *ioc.DependencyProvider) *services.RequestContextService {
		return services.NewRequestContextService(dp)
	})
	ioc.AddCloseHandler[*services.RequestContextService](dpb, func(rcs *services.RequestContextService) error {
		return rcs.Close()
	})

	dp := dpb.Build()

	return context.WithValue(context.TODO(), "scope", dp)
}
