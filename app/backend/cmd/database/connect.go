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

type GenTable struct {
	Data_value        string  `json:"data_value"`
	Confirmed         int     `json:"confirmed"`
	Deaths            int     `json:"deaths"`
	Stringency_actual float32 `json:"stringency_actual"`
	Stringency        float32 `json:"stringency"`
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
		err := rows.Scan(&genTable.Data_value, &genTable.Confirmed, &genTable.Deaths, &genTable.Stringency_actual, &genTable.Stringency)
		if err != nil {
			return nil, err
		}
		list = append(list, genTable)
	}
	return list, nil
}

func ReturnBlockValue(query string, s Settings) (b bool, err error) {
	db, err := connect(s)
	if err != nil {
		return true, err
	}
	defer db.Close()
	//Retrieve data
	row := db.QueryRow(query)
	if err = row.Scan(&b); err != nil {
		return true, err
	}
	return b, err
}
