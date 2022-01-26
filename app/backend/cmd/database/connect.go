package database

import (
	"database/sql"
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
	defer db.Close()
	if err != nil {
		return err
	}

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
