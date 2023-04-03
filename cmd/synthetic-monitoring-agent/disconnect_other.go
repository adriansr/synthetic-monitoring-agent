//go:build !windows

package main

import (
	"fmt"
	"net/http"
	"os"
	"syscall"
)

// disconnectHandler triggers a disconnection from the API.
func disconnectHandler(w http.ResponseWriter, r *http.Request) {
	// TODO(mem): this is a hack to trigger a disconnection from the
	// API, it would be cleaner to do this through a channel that
	// communicates the request to the checks updater.

	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		msg := fmt.Sprintf("%s: %s", http.StatusText(http.StatusInternalServerError), err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	// SIGUSR1 will disconnect agent from API for 1 minute
	err = p.Signal(syscall.SIGUSR1)
	if err != nil {
		msg := fmt.Sprintf("%s: %s", http.StatusText(http.StatusInternalServerError), err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "disconnecting this agent from the API for 1 minute.")
}
