// Package database for goCovid app
package database

import (
	"database/sql"
	"errors"
	"fmt"

	// Connect mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose"
)

// Settings of database
type Settings struct {
	Host  string
	Port  string
	Name  string
	User  string
	Pass  string
	Reset bool
}

// GenTable is created to parse json
type GenTable struct {
	DataValue        string  `json:"data_value"`
	Confirmed        int     `json:"confirmed"`
	Deaths           int     `json:"deaths"`
	StringencyActual float32 `json:"stringency_actual"`
	Stringency       float32 `json:"stringency"`
}

func connect(settings Settings) (db *sql.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", settings.User, settings.Pass, settings.Host, settings.Port, settings.Name)

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Migrate is exported func to use it in main module
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

// AddData is exported to use it in main module
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

// ReturnID is exported to use it in main module
func ReturnID(query string, s Settings) (id int, err error) {
	db, err := connect(s)
	if err != nil {
		return -1, err
	}
	defer db.Close()

	row := db.QueryRow(query)
	err = row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return -2, errors.New("zero rows found")
			//	} else {
			//		return -3, err
		}
	}
	return id, nil
}

// ReturnMulti is exported to use it in main module
func ReturnMulti(query string, s Settings) (list []GenTable, err error) {
	db, err := connect(s)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	genTable := GenTable{}
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err := rows.Scan(&genTable.DataValue, &genTable.Confirmed, &genTable.Deaths, &genTable.StringencyActual, &genTable.Stringency)
		if err != nil {
			return nil, err
		}
		list = append(list, genTable)
	}
	return list, nil
}

// ReturnBlockValue is exported to use it in main module
func ReturnBlockValue(query string, s Settings) (b bool, err error) {
	db, err := connect(s)
	if err != nil {
		return true, err
	}
	defer db.Close()
	// Retrieve data
	row := db.QueryRow(query)
	if err = row.Scan(&b); err != nil {
		return true, err
	}
	return b, err
}
