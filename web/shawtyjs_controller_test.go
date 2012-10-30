package web

import (
	"fmt"
	"go.3fps.com/shawty/data"
	"go.3fps.com/shawty/utils"
	"net/http"
	"testing"
	"time"
)

// getShawtyJSTestData creates the test data for ShortID tests
func getShawtyJSTestData() (config map[string]string, seed uint64, sh data.Shawties) {

	seed = uint64(time.Now().Unix()) % utils.BaseLen
	rand := utils.NewFakeRand()
	rand.Seed(seed)

	config = map[string]string{
		"SHAWTY_DOMAIN": fmt.Sprintf("shawty%d.local", seed),
		"SHAWTY_LPM":    "1000",
	}

	sh = data.NewMemSh(rand)
	sh.Create("", "http://test.com/url1", "127.0.0.1")
	sh.Create("", "http://test.com/url2", "127.0.0.1")
	sh.Create("", "http://test.com/url3", "127.0.0.1")

	return
}

func testShawtyJSInvalidUrl(t *testing.T, url string) {
	conf, _, sh := getShawtyJSTestData()
	controller := NewShawtyJSController(conf, sh)

	res := controller.GetJSResponse(url, "127.0.0.1", false)

	if res == nil {
		t.Error("No response")
		t.FailNow()
	}

	if res.Data["Success"] != 0 {
		t.Errorf("'Success' flag in the data needs to be '0', but got %v instead", res.Data["Success"])
	}

	// because we want the JS to always execute, even if fall, ensure it's a 200 response
	if res.HttpStatus != http.StatusOK {
		t.Error("JS response always needs to have 200 status")
	}
}

func testShawtyJSValidUrl(t *testing.T, url string, expectedID uint64) {
	conf, seed, sh := getShawtyJSTestData()
	controller := NewShawtyJSController(conf, sh)
	res := controller.GetJSResponse(url, "127.0.0.1", false)
	shortID := data.ShortID(expectedID, utils.ToSafeBase(seed))

	if res == nil {
		t.Error("No response")
		t.FailNow()
	}

	if res.Data["Success"] != 1 {
		t.Errorf("'Success' flag in the data needs to be '1', but got %v instead", res.Data["Success"])
	}

	if res.Data["Domain"] != conf["SHAWTY_DOMAIN"] {
		t.Errorf(
			"Wrong 'Domain' returned. '%s' expected, but '%s' returned instead",
			conf["SHAWTY_DOMAIN"], res.Data["Domain"])
	}

	if res.Data["Short"] != shortID {
		t.Errorf("Data[Short] expected to be %s, but %s returned", shortID, res.Data["Short"])
	}

	shawty, _ := sh.GetByUrl(url)
	if res.Data["Shawty"].(*data.Shawty).ID != shawty.ID {
		t.Errorf("Wrong 'Shawty' returned. Expecting %v, got %v",
			shawty, res.Data["Shawty"].(*data.Shawty))
	}
}

// TestShawtyJSInvalidUrl tests the response when the JS file is requested with an invalid Url
func TestShawtyJSInvalidUrl(t *testing.T) {
	testShawtyJSInvalidUrl(t, "some thing invalid")
}

// TestShawtyJSExistingUrl tests the response when JS is requested with a valid existing URL
func TestShawtyJSExistingUrl(t *testing.T) {
	testShawtyJSValidUrl(t, "http://test.com/url3", 3)
}

// TestShawtyJSNewUrl tests the response when JS is requested with a valid new URL
func TestShawtyJSNewUrl(t *testing.T) {
	testShawtyJSValidUrl(t, "http://test.com/something", 4)
}

// TestShawtyJSBookmarkletFlag makes sure the bookmarklet flag is correct
func TestShawtyJSBookmarkletFlag(t *testing.T) {
	conf, _, sh := getShawtyJSTestData()
	controller := NewShawtyJSController(conf, sh)

	res := controller.GetJSResponse("http://test.com/url3", "127.0.0.1", false)
	if res.Data["Bookmarklet"].(bool) != false {
		t.Error("Bookmarklet flag expecting 'false', but returned 'true'")
	}

	res = controller.GetJSResponse("http://test.com/url3", "127.0.0.1", true)
	if res.Data["Bookmarklet"].(bool) != true {
		t.Error("Bookmarklet flag expecting 'true', but returned 'false'")
	}
}

// TestShawtyJSDupDomain tests the response when JS is requested with its own domain
func TestShawtyJSDupDomain(t *testing.T) {
	conf, _, _ := getShawtyJSTestData()

	url1 := fmt.Sprintf("http://%s/something1", conf["SHAWTY_DOMAIN"])
	url2 := fmt.Sprintf("https://%s/something2", conf["SHAWTY_DOMAIN"])

	testShawtyJSInvalidUrl(t, url1)
	testShawtyJSInvalidUrl(t, url2)
}

// TestShawtyJSRateLimit tests the rate limit and make sure it doesn't allow the same IP
func TestShawtyJSRateLimit(t *testing.T) {
	conf, _, sh := getShawtyJSTestData()
	conf["SHAWTY_LPM"] = "3"

	controller := NewShawtyJSController(conf, sh)
	res := controller.GetJSResponse("http://test.com/url3", "127.0.0.1", false)
	if res.Data["Success"].(int) != 1 {
		t.Error("Existing URL should not be affected by rate limit")
	}

	newUrl := "http://example.com/something-new"

	res = controller.GetJSResponse(newUrl, "127.0.0.1", false)
	if res.Data["Success"].(int) != 0 {
		t.Error("Rate limit must be applied when creating a new URL")
	}

	_, err := sh.GetByUrl(newUrl)
	if err == nil {
		t.Error("When rate limit is hit, new URL requests must not be created")
	}
}
