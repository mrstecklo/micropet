package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/mrstecklo/micropet/services/orders/database"
)

func main() {
	logger := createLogger()
	err := godotenv.Load()
	if err != nil {
		logger.Error("failed to load .env file", "error", err.Error())
		return
	}

	dsn := os.Getenv("DATABASE_URL")
	db, err := database.NewDatabase(dsn, logger)
	if err != nil {
		return
	}
	defer db.Close()
}

func createLogger() *slog.Logger {
	options := &slog.HandlerOptions{Level: slog.LevelDebug}
	handler := slog.NewTextHandler(os.Stdout, options)
	return slog.New(handler)
}
