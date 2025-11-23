package database

import (
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createLogger() *slog.Logger {
	options := &slog.HandlerOptions{Level: slog.LevelDebug}
	handler := slog.NewTextHandler(os.Stdout, options)
	return slog.New(handler)
}

type databaseFixture struct {
	db Database
}

func TestMain(m *testing.M) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("failed to load environment:", err.Error())
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func setUpDatabaseTest(t *testing.T) databaseFixture {
	logger := createLogger()
	dsn := os.Getenv("DATABASE_URL")
	db, err := NewDatabase(dsn, logger)
	require.Nil(t, err)
	err = db.Clear()
	require.Nil(t, err)
	t.Cleanup(db.Close)
	return databaseFixture{
		db: db,
	}
}

func TestDatabase_GetOrderReturnsErrNotFound(t *testing.T) {
	f := setUpDatabaseTest(t)

	_, err := f.db.GetOrder(1)

	assert.Equal(t, ErrNotFound, err)
}

func TestDatabase_CreateOrderIsSuccessful(t *testing.T) {
	f := setUpDatabaseTest(t)

	_, err := f.db.CreateOrder("something")

	assert.Nil(t, err)
}

func TestDatabase_GetOrderReturnsCreatedOrder(t *testing.T) {
	f := setUpDatabaseTest(t)
	id, err := f.db.CreateOrder("something")
	require.Nil(t, err)

	order, err := f.db.GetOrder(id)

	assert.Nil(t, err)
	assert.Equal(t, id, order.ID)
	assert.Equal(t, "something", order.Title)
}

func TestDatabase_StoresStateWhenClosed(t *testing.T) {
	logger := createLogger()
	dsn := os.Getenv("DATABASE_URL")
	db1, err := NewDatabase(dsn, logger)
	require.Nil(t, err)
	err = db1.Clear()
	if err != nil {
		db1.Close()
		t.Fatal(err.Error())
	}
	id, err := db1.CreateOrder("duck")
	db1.Close()
	require.Nil(t, err)
	db2, err := NewDatabase(dsn, logger)
	require.Nil(t, err)
	defer db2.Close()

	order, err := db2.GetOrder(id)

	assert.Nil(t, err)
	assert.Equal(t, id, order.ID)
	assert.Equal(t, "duck", order.Title)
}

func TestDatabase_GetOrderReturnsErrNotFoundIfWrongId(t *testing.T) {
	f := setUpDatabaseTest(t)
	id, err := f.db.CreateOrder("something")
	require.Nil(t, err)

	_, err = f.db.GetOrder(id + 1)

	assert.Equal(t, ErrNotFound, err)
}

func TestDatabase_GetOrderReturnsRespectiveOrder(t *testing.T) {
	f := setUpDatabaseTest(t)
	id1, err := f.db.CreateOrder("something")
	require.Nil(t, err)
	id2, err := f.db.CreateOrder("duck")
	require.Nil(t, err)
	id3, err := f.db.CreateOrder("pickle")
	require.Nil(t, err)

	order2, err := f.db.GetOrder(id2)
	require.Nil(t, err)
	order1, err := f.db.GetOrder(id1)
	require.Nil(t, err)
	order3, err := f.db.GetOrder(id3)
	require.Nil(t, err)

	assert.Equal(t, id1, order1.ID)
	assert.Equal(t, "something", order1.Title)
	assert.Equal(t, id2, order2.ID)
	assert.Equal(t, "duck", order2.Title)
	assert.Equal(t, id3, order3.ID)
	assert.Equal(t, "pickle", order3.Title)
}
