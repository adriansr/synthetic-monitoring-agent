//go:build windows

package main

import (
	"net/http"
)

// disconnectHandler triggers a disconnection from the API.
func disconnectHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Disconnect feature is not supported on Windows", http.StatusInternalServerError)
}
