package database

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose"
)

type Settings struct {
	Host  string
	Port  string
	Name  string
	User  string
	Pass  string
	Reset bool
}

func connect(settings Settings) (db *sql.DB, err error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", settings.User, settings.Pass, settings.Host, settings.Port, settings.Name)

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Migrate(path string, s Settings) error {

	db, err := connect(s)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := goose.SetDialect("mysql"); err != nil {
		return err
	}

	if s.Reset {
		if err := goose.DownTo(db, path, 0); err != nil {
			return err
		}
	}

	if err := goose.Up(db, path); err != nil {
		return err
	}
	return nil
}

func AddData(query string, s Settings) error {

	db, err := connect(s)
	if err != nil {
		return err
	}
	defer db.Close()

	if _, err := db.Exec(query); err != nil {
		return err
	}

	return nil
}

func ReturnId(query string, s Settings) (id int, err error) {

	db, err := connect(s)
	if err != nil {
		return -1, err
	}
	defer db.Close()

	row := db.QueryRow(query)
	err = row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return -2, errors.New("Zero rows found")
		} else {
			return -3, err
		}
	}
	return id, nil
}
