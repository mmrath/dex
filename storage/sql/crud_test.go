package sql

import (
	"database/sql"
	"os"
	"reflect"
	"testing"
)

func TestDecoder(t *testing.T) {
	db := getDB(t)
	defer db.Close()

	if _, err := db.Exec(`create table foo_decoder_test ( id integer primary key, bar bytea );`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into foo_decoder_test ( id, bar ) values (1, $1);`, []byte(`["a", "b"]`)); err != nil {
		t.Fatal(err)
	}
	var got []string
	if err := db.QueryRow(`select bar from foo_decoder_test where id = 1;`).Scan(decoder(&got)); err != nil {
		t.Fatal(err)
	}
	want := []string{"a", "b"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("wanted %q got %q", want, got)
	}
	_, _ = db.Exec(`drop table foo_decoder_test`)
}

func TestEncoder(t *testing.T) {
	db := getDB(t)
	defer db.Close()

	if _, err := db.Exec(`create table foo_encoder_test ( id integer primary key, bar bytea );`); err != nil {
		t.Fatal(err)
	}
	put := []string{"a", "b"}
	if _, err := db.Exec(`insert into foo_encoder_test ( id, bar ) values (1, $1)`, encoder(put)); err != nil {
		t.Fatal(err)
	}

	var got []byte
	if err := db.QueryRow(`select bar from foo_encoder_test where id = 1;`).Scan(&got); err != nil {
		t.Fatal(err)
	}
	want := []byte(`["a","b"]`)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("wanted %q got %q", want, got)
	}
	_, _ = db.Exec(`drop table foo_encoder_test`)
}

func getDB(t *testing.T) *sql.DB {
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
	return db
}
