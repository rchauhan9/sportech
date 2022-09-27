package database

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// Logger Implementation of golang-migrate logger.
type Logger struct {
	logger log.Logger
}

func (l *Logger) Printf(format string, v ...interface{}) {
	level.Info(l.logger).Log("msg", fmt.Sprintf(format, v...))
}

func (l *Logger) Verbose() bool {
	return true
}

type Migrator struct {
	migrator    *migrate.Migrate
	sourceURL   string
	databaseURL string
}

func NewMigrator(databaseURL string, migrationPath string, logger log.Logger) (*Migrator, error) {
	var fileSource string
	if strings.HasPrefix(migrationPath, ".") {
		// Relative path
		fileSource = "file://"
	} else {
		// Absolute path
		fileSource = "file:///"
	}
	sourceURL := fmt.Sprintf("%s%s", fileSource, migrationPath)
	embeddedMigrator, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return nil, err
	}
	embeddedMigrator.Log = &Logger{logger: logger}
	return &Migrator{embeddedMigrator, sourceURL, databaseURL}, nil
}

func (m *Migrator) MigrateDb() error {
	if err := m.migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func (m *Migrator) UnmigrateDB() error {
	if err := m.migrator.Down(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func (m *Migrator) PurgeDB() error {
	if err := m.migrator.Drop(); err != nil {
		return err
	}
	// Refresh the underlying migrate instance after dropping the database; this resolves an issue
	// with restoring the schema migrations: https://github.com/golang-migrate/migrate/issues/226
	if err := m.refresh(); err != nil {
		return err
	}
	return nil
}

func (m *Migrator) Close() error {
	sourceErr, dbErr := m.migrator.Close()
	if sourceErr != nil {
		return sourceErr
	}
	if dbErr != nil {
		return dbErr
	}
	return nil
}

func (m *Migrator) refresh() error {
	var err error

	if err = m.Close(); err != nil {
		return err
	}

	logger := m.migrator.Log
	if m.migrator, err = migrate.New(m.sourceURL, m.databaseURL); err != nil {
		return err
	}
	m.migrator.Log = logger

	return nil
}
