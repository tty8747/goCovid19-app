package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/spf13/viper"
	"github.com/tty8747/goCovid19/cmd/database"
)

func main() {
	app := &application{}

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	if err := initConfigs(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	app.settings.migrationDir = viper.GetString("migration_dir")
	app.settings.endPoint = viper.GetString("app_endpoint")
	app.dbSettings.Host = viper.GetString("db_host")
	app.dbSettings.Name = viper.GetString("db_name")
	app.dbSettings.Port = viper.GetString("db_port")
	app.dbSettings.User = viper.GetString("db_user")
	app.dbSettings.Pass = viper.GetString("db_pass")
	app.dbSettings.Reset = viper.GetBool("data_reset")

	addr := flag.String("addr", app.settings.endPoint, "API HTTP address")
	getData := flag.Bool("fill", false, "If it is true it fills db")
	flag.Parse()

	if *getData {
		// Migrations
		// if err := database.Migrate(app.settings.migrationDir, app.dbSettings); err != nil {
		// 	log.Println(err)
		// 	app.errLog.Fatal(err)
		// }

		// set block
		if err := app.setBlock(true); err != nil {
			app.errLog.Fatal(err)
		}

		// Get data
		app.parser()

		// Purge data from db tables
		app.purgeTables()

		// Insert data
		app.insertData()

		return
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Start api-server on %s", *addr)
	err := srv.ListenAndServe()
	errLog.Fatal(err)
}

type application struct {
	errLog *log.Logger
	// infoLog     *log.Logger
	listOfDates []string
	cList       []string // country list
	listObj     []Obj
	settings    appSettings
	dbSettings  database.Settings
	block       bool
	valueUpdate bool
}

// Obj is made for selected object
type Obj struct {
	DateValue            string  `json:"date_value"`
	CountryCode          string  `json:"country_code"`
	Confirmed            int     `json:"confirmed"`
	Deaths               int     `json:"deaths"`
	StringencyActual     float64 `json:"stringency_actual"`
	Stringency           float64 `json:"stringency"`
	StringencyLegacy     float64 `json:"stringency_legacy"`
	StringencyLegacyDisp float64 `json:"stringency_legacy_disp"`
}

type appSettings struct {
	migrationDir string
	endPoint     string
}

func initConfigs() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	// gets env variables
	viper.AutomaticEnv()
	// gets data from config
	return viper.ReadInConfig()
}
