package sql

import (
	"database/sql"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"os"
	"testing"
)

func TestMigrate(t *testing.T) {

	host := os.Getenv(testPostgresEnv)
	if host == "" {
		t.Skipf("test environment variable %q not set, skipping", testPostgresEnv)
	}
	baseCfg := &Postgres{
		NetworkDB: NetworkDB{
			Database: getenv("DEX_POSTGRES_DATABASE", "postgres"),
			User:     getenv("DEX_POSTGRES_USER", "postgres"),
			Password: getenv("DEX_POSTGRES_PASSWORD", "postgres"),
			Host:     host,
		},
		SSL: SSL{
			Mode: pgSSLDisable, // Postgres container doesn't support SSL.
		}}

	dataSourceName := baseCfg.createDataSourceName()

	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	logger := &logrus.Logger{
		Out:       os.Stderr,
		Formatter: &logrus.TextFormatter{DisableColors: true},
		Level:     logrus.DebugLevel,
	}

	errCheck := func(err error) bool {
		sqlErr, ok := err.(pq.Error)
		if !ok {
			return false
		}
		return sqlErr.Code != ""
	}

	c := &conn{db, flavorPostgres, logger, errCheck}
	for _, want := range []int{len(migrations), 0} {
		got, err := c.migrate()
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Errorf("expected %d migrations, got %d", want, got)
		}
	}
}
