package database

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

var PostgresDriver = &pqDriver{}

func CreateDnsPostgres() string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"
	return dsn
}

type pqDriver struct {
	db *sql.DB
}

func (d *pqDriver) Open(dsn string) (*sql.DB, error) {
	newDb, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	d.db = newDb
	return d.db, nil
}

func (d *pqDriver) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

func (d *pqDriver) GetDB() *sql.DB {
	return d.db
}

func (d *pqDriver) Migrate() error {
	db := d.GetDB()

	// read all migration script from ./scripts/create_table.sql
	script, err := os.ReadFile("./database/scripts/create_table.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(script))
	if err != nil {
		return err
	}

	return nil
}
