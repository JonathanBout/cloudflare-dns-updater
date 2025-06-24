package main

import (
	"fmt"
	"io"
	"net/http"
)

func getPing(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, "pong\n")
}

/*
postDynDNS is the handler for the / endpoint. It updates the DNS records for the configured with the supplied IP addresses.
*/
func postDynDNS(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	query := r.URL.Query()

	username := query.Get("username")
	password := query.Get("password")

	if username != auth.Username || password != auth.Password {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	ipv4 := query.Get("ipv4")
	ipv6 := query.Get("ipv6")

	fmt.Printf("Authenticated for DNS update. Supplied new addresses:\n")
	fmt.Printf("ipv4: %s\n", ipv4)
	fmt.Printf("ipv6: %s\n", ipv6)

	err := updateDNS(ipv4, ipv6)

	if err != nil {
		fmt.Printf("error updating DNS: %s\n", err)
		http.Error(w, "error updating DNS", http.StatusInternalServerError)
		return
	}

	_, _ = io.WriteString(w, "DNS updated\n")
}
