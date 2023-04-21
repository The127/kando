package db

import (
	"database/sql"
	"github.com/xlab/handysort"
	"io/fs"
	"kando-backend/log"
	"kando-backend/migrations"
	"path/filepath"
	"sort"
	"strings"
)

func MigrateDatabase(dbConnection *sql.DB) {
	log.Logger.Info("Applying migrations...")

	ensureMigrationTableExists(dbConnection)
	applyMissingMigrations(dbConnection)
}

func ensureMigrationTableExists(dbConnection *sql.DB) {
	_, err := dbConnection.Exec("create table if not exists \"__migrations\" (\"name\" text not null, \"timestamp\" timestamptz not null default now());")
	if err != nil {
		log.Logger.Fatalf("Failed to create migrations table: %v", err)
	}
}

func applyMissingMigrations(dbConnection *sql.DB) {
	migrationFilePaths := getMigrationFilePaths()
	appliedMigrations := getAppliedMigrationNames(dbConnection)

	for _, appliedMigration := range appliedMigrations {
		log.Logger.Debugf("Migration %s already applied", appliedMigration)
	}

	var missingMigrationFilePaths []string

	for _, migrationFilePath := range migrationFilePaths {
		fileName := filepath.Base(migrationFilePath)
		migrationName := fileName[:len(fileName)-len(".sql")]

		if !isMigrationAlreadyApplied(migrationName, appliedMigrations) {
			missingMigrationFilePaths = append(missingMigrationFilePaths, migrationFilePath)
		}
	}

	log.Logger.Debugf("Found %d migrations to apply...", len(missingMigrationFilePaths))

	sort.Sort(handysort.Strings(missingMigrationFilePaths))

	for _, migrationFilePath := range missingMigrationFilePaths {
		migrationName := filepath.Base(migrationFilePath)
		migrationName = migrationName[:len(migrationName)-len(".sql")]

		fileBytes, err := fs.ReadFile(migrations.MigrationFs, migrationFilePath)
		if err != nil {
			log.Logger.Fatalf("Could not read migration file %s: %v", migrationFilePath, err)
		}
		migration := string(fileBytes)

		executeMigration(migrationName, migration, dbConnection)
	}
}

func executeMigration(migrationName string, migration string, dbConnection *sql.DB) {
	log.Logger.Debugf("Executing migration %s...", migrationName)

	tx, err := dbConnection.Begin()
	if err != nil {
		log.Logger.Fatalf("Failed to create transaction: %v", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec("insert into \"__migrations\" (\"name\") values($1);", migrationName)
	if err != nil {
		log.Logger.Fatalf("Failed to insert migration table entry for %s: %v", migrationName, err)
	}

	_, err = tx.Exec(migration)
	if err != nil {
		log.Logger.Fatalf("Failed to execute migration %s: %v", migrationName, err)
	}

	err = tx.Commit()
	if err != nil {
		log.Logger.Fatalf("Failed to commit transaction: %v", err)
	}
}

func isMigrationAlreadyApplied(migrationName string, appliedMigrations []string) bool {
	for _, migration := range appliedMigrations {
		if migration == migrationName {
			return true
		}
	}

	return false
}

func getMigrationFilePaths() []string {
	var migrationFiles []string

	err := fs.WalkDir(migrations.MigrationFs, ".",
		func(path string, d fs.DirEntry, err error) error {
			if strings.HasSuffix(d.Name(), ".sql") {
				migrationFiles = append(migrationFiles, path)
			}
			return nil
		})

	if err != nil {
		log.Logger.Fatalf("Failed to find migrations: %v", err)
	}

	return migrationFiles
}

func getAppliedMigrationNames(dbConnection *sql.DB) []string {
	rows, err := dbConnection.Query("select \"name\" from \"__migrations\";")
	if err != nil {
		log.Logger.Fatalf("Failed to query applied migrations: %v", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Logger.Fatalf("Failed to close rows: %v", err)
		}
	}(rows)

	var appliedMigrations []string

	for rows.Next() {
		var appliedMigration string

		err = rows.Scan(&appliedMigration)
		if err != nil {
			log.Logger.Fatalf("Failed to scan applied migrations: %v", err)
		}

		appliedMigrations = append(appliedMigrations, appliedMigration)
	}

	return appliedMigrations
}
