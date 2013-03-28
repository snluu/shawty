package utils

import (
	"net/http"
	"strings"
)

func HttpRemoteIP(r *http.Request) string {
	idx := strings.LastIndex(r.RemoteAddr, ":")
	if idx == -1 {
		return r.RemoteAddr
	}
	return r.RemoteAddr[:idx]
}
