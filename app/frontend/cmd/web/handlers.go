package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.notAllowed(w)
		return
	}
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
		"./ui/html/header.partial.tmpl",
	}

	// gets list of alpha-3
	// var list []string = []string{"ABW", "AFG", "AGO", "ALB", "AND", "ARE", "ARG", "AUS", "AUT", "AZE", "BDI", "BEL", "BEN", "BFA", "BGD", "BGR", "BHR", "BHS", "BIH", "BLR", "BLZ", "BMU", "BOL", "BRA", "BRB", "BRN", "BTN", "BWA", "CAF", "CAN", "CHE", "CHL", "CHN", "CIV", "CMR", "COD", "COG", "COL", "CPV", "CRI", "CUB", "CYP", "CZE", "DEU", "DJI", "DMA", "DNK", "DOM", "DZA", "ECU", "EGY", "ERI", "ESP", "EST", "ETH", "FIN", "FJI", "FRA", "FRO", "GAB", "GBR", "GEO", "GHA", "GIN", "GMB", "GRC", "GRL", "GTM", "GUM", "GUY", "HKG", "HND", "HRV", "HTI", "HUN", "IDN", "IND", "IRL", "IRN", "IRQ", "ISL", "ISR", "ITA", "JAM", "JOR", "JPN", "KAZ", "KGZ", "KHM", "KIR", "KOR", "KWT", "LAO", "LBN", "LBR", "LBY", "LIE", "LKA", "LSO", "LTU", "LUX", "LVA", "MAC", "MAR", "MCO", "MDA", "MDG", "MEX", "MLI", "MLT", "MMR", "MNG", "MOZ", "MRT", "MUS", "MWI", "MYS", "NAM", "NER", "NGA", "NIC", "NLD", "NOR", "NPL", "NZL", "OMN", "PAK", "PAN", "PER", "PHL", "PNG", "POL", "PRI", "PRT", "PRY", "PSE", "QAT", "RKS", "ROU", "RUS", "RWA", "SAU", "SDN", "SEN", "SGP", "SLB", "SLE", "SLV", "SMR", "SOM", "SRB", "SSD", "SUR", "SVK", "SVN", "SWE", "SWZ", "SYC", "SYR", "TCD", "TGO", "THA", "TJK", "TKM", "TLS", "TON", "TTO", "TUN", "TUR", "TWN", "TZA", "UGA", "UKR", "URY", "USA", "UZB", "VEN", "VIR", "VNM", "VUT", "YEM", "ZAF", "ZMB", "ZWE"}

	//	if app.MapCountries == nil {
	//		app.MapCountries = app.getCountryNames(list)
	//	}

	//	for key, value := range app.MapCountries {
	//		log.Println(key, value)
	//	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		app.serverErr(w, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", &app)
	if err != nil {
		log.Println(err)
		app.errLog.Println(err.Error())
		app.serverErr(w, err)
	}

}

func (app *application) query(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost || r.Method == http.MethodGet {

		connString, ok := app.buildLink(r, *app.api.hostname, *app.api.port, *app.api.apiVers)
		if ok {

			// --- start of data preparation
			// gets raw data
			log.Println(connString)
			response, err := http.Get(connString)
			if err != nil {
				app.errLog.Fatalln(err.Error())
			}
			defer response.Body.Close()

			// gets array of raw bytes
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				app.errLog.Fatalln(err.Error())
			}

			// puts data in a struct
			if err := json.Unmarshal(body, &app.Data); err != nil {
				app.errLog.Fatalln(err.Error())
			}
		} else {
			app.Message = connString
		}
		// --- end of data preparation

		app.CountrySelFull = app.setCountryNameFull(app.CountrySel)

		files := []string{
			"./ui/html/home.page.tmpl",
			"./ui/html/base.layout.tmpl",
			"./ui/html/footer.partial.tmpl",
			"./ui/html/header.partial.tmpl",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.errLog.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
			app.serverErr(w, err)
			return
		}

		err = ts.Execute(w, &app)
		if err != nil {
			log.Println(err)
			app.errLog.Println(err.Error())
			app.serverErr(w, err)
		}

	} else {
		app.notAllowed(w)
		return
	}
}
