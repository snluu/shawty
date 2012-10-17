package web

import (
	"net/http"
	"testing"
)

// TestIndexPage makes sure the index page returns the data with the correct domain
func TestHomeIndexPage(t *testing.T) {
	domain := "shawty.local"
	conf := map[string]string{"SHAWTY_DOMAIN": domain}
	index := NewHomeController(conf)
	res := index.Index()
	if res == nil {
		t.Error("Response is nil")
		t.FailNow()
	}
	if res.HttpStatus != http.StatusOK {
		t.Errorf("HTTP status is %d, but %d expected", res.HttpStatus, http.StatusOK)
	}
	if res.Data == nil || len(res.Data) == 0 {
		t.Error("Response did not come back with any data")
		t.FailNow()
	}
	if res.Data["Domain"] != domain {
		t.Errorf("Response's 'Domain' key is expected to be '%s', but returns '%s'", domain, res.Data["Domain"])
	}
}
