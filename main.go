package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/cloudflare/cloudflare-go/v4"

	"github.com/cloudflare/cloudflare-go/v4/option"
)

var cf *cloudflare.Client
var auth *Credentials

type Credentials struct {
	Username string
	Password string
}

func main() {
	fmt.Printf("connecting to cloudflare...\n")

	cf = cloudflare.NewClient(
		option.WithAPIKey(os.Getenv("CLOUDFLARE_API_TOKEN")),
	)

	auth = &Credentials{
		Username: os.Getenv("DYNDNS_USERNAME"),
		Password: os.Getenv("DYNDNS_PASSWORD"),
	}

	_, err := cf.User.Tokens.Verify(context.TODO())

	if err != nil {
		fmt.Printf("error verifying token: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("token verified")

	mux := http.NewServeMux()
	mux.HandleFunc("/", postDynDNS)
	mux.HandleFunc("/ping", getPing)

	listenAddr, found := os.LookupEnv("DYNDNS_LISTEN_ADDR")
	if !found {
		listenAddr = ":5000"
	}

	var listenAddrFriendly string

	// if listenAddr starts with a colon, prepend the hostname
	if listenAddr[0] == ':' {
		listenAddrFriendly = "localhost" + listenAddr
	} else {
		listenAddrFriendly = listenAddr
	}

	fmt.Printf("server listening on http://%s\n", listenAddrFriendly)
	err = http.ListenAndServe(listenAddr, mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
