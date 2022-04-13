package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestGettingCountryNames(t *testing.T) {
	var list []string = []string{"ABW", "AFG", "AGO", "ALB", "AND", "ARE", "ARG", "AUS", "AUT", "AZE", "BDI", "BEL", "BEN", "BFA", "BGD", "BGR", "BHR", "BHS", "BIH", "BLR", "BLZ", "BMU", "BOL", "BRA", "BRB", "BRN", "BTN", "BWA", "CAF", "CAN", "CHE", "CHL", "CHN", "CIV", "CMR", "COD", "COG", "COL", "CPV", "CRI", "CUB", "CYP", "CZE", "DEU", "DJI", "DMA", "DNK", "DOM", "DZA", "ECU", "EGY", "ERI", "ESP", "EST", "ETH", "FIN", "FJI", "FRA", "FRO", "GAB", "GBR", "GEO", "GHA", "GIN", "GMB", "GRC", "GRL", "GTM", "GUM", "GUY", "HKG", "HND", "HRV", "HTI", "HUN", "IDN", "IND", "IRL", "IRN", "IRQ", "ISL", "ISR", "ITA", "JAM", "JOR", "JPN", "KAZ", "KGZ", "KHM", "KIR", "KOR", "KWT", "LAO", "LBN", "LBR", "LBY", "LIE", "LKA", "LSO", "LTU", "LUX", "LVA", "MAC", "MAR", "MCO", "MDA", "MDG", "MEX", "MLI", "MLT", "MMR", "MNG", "MOZ", "MRT", "MUS", "MWI", "MYS", "NAM", "NER", "NGA", "NIC", "NLD", "NOR", "NPL", "NZL", "OMN", "PAK", "PAN", "PER", "PHL", "PNG", "POL", "PRI", "PRT", "PRY", "PSE", "QAT", "RKS", "ROU", "RUS", "RWA", "SAU", "SDN", "SEN", "SGP", "SLB", "SLE", "SLV", "SMR", "SOM", "SRB", "SSD", "SUR", "SVK", "SVN", "SWE", "SWZ", "SYC", "SYR", "TCD", "TGO", "THA", "TJK", "TKM", "TLS", "TON", "TTO", "TUN", "TUR", "TWN", "TZA", "UGA", "UKR", "URY", "USA", "UZB", "VEN", "VIR", "VNM", "VUT", "YEM", "ZAF", "ZMB", "ZWE"}
	app := &application{}
	if app.getCountryNames(list)["Australia"] != "AUS" {
		t.Fatalf(`Country name and alpha3 code are different!, want = AUS, got = %s`, app.getCountryNames(list)["Australia"])
	}
	if app.getCountryNames(list)["Kosovo"] != "RKS" {
		t.Fatalf(`Country name and alpha3 code are different!, want = RKS, got = %s`, app.getCountryNames(list)["Kosovo"])
	}
	if app.getCountryNames(list)["South Africa"] != "ZAF" {
		t.Fatalf(`Country name and alpha3 code are different!, want = ZAF, got = %s`, app.getCountryNames(list)["South Africa"])
	}
}
