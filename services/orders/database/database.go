package database

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/url"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mrstecklo/micropet/services/orders/orders"
)

var ErrNotFound = errors.New("not found")

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

func (db Database) GetOrder(id int) (orders.Order, error) {
	var order orders.Order
	err := db.db.QueryRow("SELECT id, title FROM orders WHERE id = $1", id).
		Scan(&order.ID, &order.Title)
	if err == sql.ErrNoRows {
		return order, ErrNotFound
	}
	return order, err
}

func (db Database) CreateOrder(title string) (int, error) {
	var id int
	err := db.db.QueryRow("INSERT INTO orders (title) VALUES ($1) RETURNING id", title).
		Scan(&id)
	return id, err
}

func (db Database) Clear() error {
	_, err := db.db.Exec("DELETE FROM orders")
	return err
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
