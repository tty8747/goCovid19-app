package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	app := &application{}

	addr := flag.String("addr", "localhost:4000", "HTTP address")
	app.api.hostname = flag.String("hostname", "localhost", "API hostname")
	app.api.port = flag.String("port", "8080", "API port")
	app.api.apiVers = flag.String("apiVers", "v1", "API version")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// gets list of alpha-3
	var list = []string{"ABW", "AFG", "AGO", "ALB", "AND", "ARE", "ARG", "AUS", "AUT", "AZE", "BDI", "BEL", "BEN", "BFA", "BGD", "BGR", "BHR", "BHS", "BIH", "BLR", "BLZ", "BMU", "BOL", "BRA", "BRB", "BRN", "BTN", "BWA", "CAF", "CAN", "CHE", "CHL", "CHN", "CIV", "CMR", "COD", "COG", "COL", "CPV", "CRI", "CUB", "CYP", "CZE", "DEU", "DJI", "DMA", "DNK", "DOM", "DZA", "ECU", "EGY", "ERI", "ESP", "EST", "ETH", "FIN", "FJI", "FRA", "FRO", "GAB", "GBR", "GEO", "GHA", "GIN", "GMB", "GRC", "GRL", "GTM", "GUM", "GUY", "HKG", "HND", "HRV", "HTI", "HUN", "IDN", "IND", "IRL", "IRN", "IRQ", "ISL", "ISR", "ITA", "JAM", "JOR", "JPN", "KAZ", "KGZ", "KHM", "KIR", "KOR", "KWT", "LAO", "LBN", "LBR", "LBY", "LIE", "LKA", "LSO", "LTU", "LUX", "LVA", "MAC", "MAR", "MCO", "MDA", "MDG", "MEX", "MLI", "MLT", "MMR", "MNG", "MOZ", "MRT", "MUS", "MWI", "MYS", "NAM", "NER", "NGA", "NIC", "NLD", "NOR", "NPL", "NZL", "OMN", "PAK", "PAN", "PER", "PHL", "PNG", "POL", "PRI", "PRT", "PRY", "PSE", "QAT", "RKS", "ROU", "RUS", "RWA", "SAU", "SDN", "SEN", "SGP", "SLB", "SLE", "SLV", "SMR", "SOM", "SRB", "SSD", "SUR", "SVK", "SVN", "SWE", "SWZ", "SYC", "SYR", "TCD", "TGO", "THA", "TJK", "TKM", "TLS", "TON", "TTO", "TUN", "TUR", "TWN", "TZA", "UGA", "UKR", "URY", "USA", "UZB", "VEN", "VIR", "VNM", "VUT", "YEM", "ZAF", "ZMB", "ZWE"}

	// gets data into a map
	if app.MapCountries == nil {
		app.MapCountries = app.getCountryNames(list)
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Start web-server on %s", *addr)
	err := srv.ListenAndServe()
	errLog.Fatal(err)
}

type application struct {
	errLog *log.Logger
	// infoLog          *log.Logger
	DateFrom, DateTo string
	RadioDD          string
	CountrySel       string
	Message          string
	Data             []genTable
	api              apiVariables
	MapCountries     map[string]string
	block            bool
}

type genTable struct {
	DataValue        string  `json:"data_value"`
	Confirmed        int     `json:"confirmed"`
	Deaths           int     `json:"deaths"`
	StringencyActual float32 `json:"stringency_actual"`
	Stringency       float32 `json:"stringency"`
}

type apiVariables struct {
	hostname, port, apiVers *string
}
