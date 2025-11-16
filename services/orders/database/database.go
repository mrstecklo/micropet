package database

import (
	"database/sql"
	"log/slog"
	"net/url"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Database struct {
	db     *sql.DB
	logger *slog.Logger
}

func (db Database) Close() {
	db.logger.Info("closing database")
	err := db.db.Close()
	if err != nil {
		db.logger.Error("failed to close database", "error", err.Error())
	}
}

func NewDatabase(dsn string, logger *slog.Logger) (Database, error) {
	safeDsn, err := hideUrlAuthority(dsn)
	if err != nil {
		logger.Error("failed to parse database URL", "error", err.Error())
		return Database{}, nil
	}
	logger.Info("opening database", "url", safeDsn)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Error("failed to open database", "url", safeDsn)
		return Database{}, err
	}
	return Database{
		db:     db,
		logger: logger,
	}, nil
}

func hideUrlAuthority(input string) (string, error) {
	parsed, err := url.Parse(input)
	if err != nil {
		return "", err
	}
	parsed.User = url.UserPassword("xxx", "xxx")
	return parsed.String(), nil
}
