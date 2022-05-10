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
	var alpha3 = []string{"AUS", "RKS", "ZAF"}
	app := &application{}
	if app.getCountryNames(alpha3)["Australia"] != "AUS" {
		t.Fatalf(`Country name and alpha3 code are different!, want = AUS, got = %s`, app.getCountryNames(alpha3)["Australia"])
	}
	if app.getCountryNames(alpha3)["Kosovo"] != "RKS" {
		t.Fatalf(`Country name and alpha3 code are different!, want = RKS, got = %s`, app.getCountryNames(alpha3)["Kosovo"])
	}
	if app.getCountryNames(alpha3)["South Africa"] != "ZAF" {
		t.Fatalf(`Country name and alpha3 code are different!, want = ZAF, got = %s`, app.getCountryNames(alpha3)["South Africa"])
	}
}
