package web

import (
	"fmt"
	"github.com/3fps/shawty/data"
	"github.com/3fps/shawty/utils"
	"net/http"
	"testing"
	"time"
)

// getShortIDTestData creates the test data for ShortID tests
func getShortIDTestData() (config map[string]string, seed uint64, sh data.Shawties) {

	seed = uint64(time.Now().Unix()) % utils.BaseLen
	rand := utils.NewFakeRand()
	rand.Seed(seed)

	config = map[string]string{"SHAWTY_DOMAIN": fmt.Sprintf("shawty%d.local", seed), "SHAWTY_LPM": "1000"}

	sh = data.NewMemSh(rand)
	sh.Create("", "http://test.com/url1", "127.0.0.1")
	sh.Create("", "http://test.com/url2", "127.0.0.1")
	sh.Create("", "http://test.com/url3", "127.0.0.1")

	return
}

// TestShortIDFound tests the response when a short ID is requested which is in the data store
func TestShortIDFound(t *testing.T) {
	conf, seed, sh := getShortIDTestData()
	controller := NewShortIDController(conf, sh)

	shortID := data.ShortID(2, utils.ToSafeBase(seed))
	res := controller.Respond(shortID)

	if res == nil {
		t.Error("No response")
		t.FailNow()
	}

	if res.HttpStatus != http.StatusMovedPermanently {
		t.Errorf(
			"HTTP status needs to be %d when Shawty found, but %d returned instead",
			http.StatusMovedPermanently, res.HttpStatus)
	}

	if res.Data["Domain"] != conf["SHAWTY_DOMAIN"] {
		t.Errorf(
			"Wrong 'Domain' returned. '%s' expected, but '%s' returned instead",
			conf["SHAWTY_DOMAIN"], res.Data["Domain"])
	}

	shawty2, _ := sh.GetByUrl("http://test.com/url2")

	if res.Data["Shawty"].(*data.Shawty).ID != shawty2.ID {
		t.Errorf("Wrong 'Shawty' returned. Expected %v, but %v returned", shawty2, res.Data["Shawty"])
		t.FailNow()
	}

	// make sure it increased hits
	if shawty2.Hits != 1 {
		t.Error("Respond is expected to increase hits when the requested Shawty is found")
	}
}

// TestShortIDNotFound tests the response when a short ID is requested which is in the data store
func TestShortIDNotFound(t *testing.T) {
	conf, _, sh := getShortIDTestData()
	controller := NewShortIDController(conf, sh)

	shortID := data.ShortID(5, utils.ToSafeBase(1))
	res := controller.Respond(shortID)

	if res == nil {
		t.Error("No response")
		t.FailNow()
	}

	if res.HttpStatus != http.StatusNotFound {
		t.Errorf(
			"HTTP status needs to be %d when Shawty found, but %d returned instead",
			http.NotFound, res.HttpStatus)
	}
}
