package tests

import (
	"database/sql"
	"fmt"
	"github.com/cenkalti/backoff"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	postgresDB *sql.DB
)

func TestMain(m *testing.M) {
	path := "../../docker-compose.yaml"
	compose, err := StartDockerCompose(path)
	if err != nil {
		log.Fatal(err)
	}

	host := os.Getenv("APP_HOST")
	if err = Retry(func() error {
		postgresDB, err = sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", "user", "123", host, "5432", "postgres"))
		if err != nil {
			return err
		}
		return postgresDB.Ping()
	}); err != nil {
		log.Fatalf("could not connect to postgres: %s", err)
	}

	code := m.Run()

	if err := compose.Down().Error; err != nil {
		log.Fatalf("could not stop containers from %v: %v", path, err)
	}

	os.Exit(code)
}

func StartDockerCompose(paths ...string) (*tc.LocalDockerCompose, error) {
	id := strings.ToLower(uuid.New().String())

	compose := tc.NewLocalDockerCompose(paths, id)
	execError := compose.
		WithCommand([]string{"up", "-d"}).
		Invoke()

	err := execError.Error
	if err != nil {
		return nil, fmt.Errorf("could not start containers from %v: %v", paths, err)
	}

	return compose, nil
}

func Retry(f func() error) error {
	b := backoff.NewExponentialBackOff()
	b.MaxInterval = time.Second * 5
	b.MaxElapsedTime = time.Minute

	return backoff.Retry(f, b)
}

func TestPostgres(t *testing.T) {
	rows, err := postgresDB.Query("SELECT * FROM BookStore")
	require.NoError(t, err)

	var table string
	for rows.Next() {
		err := rows.Scan(&table)
		require.NoError(t, err)
		t.Logf("table: %s", table)
	}

	require.NoError(t, rows.Err())
}
