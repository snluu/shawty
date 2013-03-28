package web

import "net/http"

// closeReqBody closes the request Body if it's not nil
func closeReqBody(r *http.Request) {
	if r != nil && r.Body != nil {
		r.Body.Close()
	}
}
