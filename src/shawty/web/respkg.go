package web

import (
	"net/http"
)

// ResPkg represents a response package
type ResPkg struct {
	HttpStatus int
	Errors     []error
	Data       map[string]interface{}
}

// NewResPkg creates and initialize and new response package with 200 status
func NewResPkg() *ResPkg {
	rp := new(ResPkg)
	rp.HttpStatus = http.StatusOK
	rp.Data = make(map[string]interface{})
	return rp
}
