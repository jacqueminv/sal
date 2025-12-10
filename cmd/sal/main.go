package main

import (
	"fmt"
	"github.com/jacqueminv/sal/internal/local"
	"github.com/jacqueminv/sal/internal/strava"
	"log"
	"os"
)

func main() {
	clientId := os.Getenv("SAL_CLIENT_ID")
	clientSecret := os.Getenv("SAL_CLIENT_SECRET")
	if clientId == "" || clientSecret == "" {
		log.Fatal("Missing SAL_CLIENT_ID, SAL_CLIENT_SECRET")
		os.Exit(1)
	}

	var tokens strava.StravaToken
	strava.Init(clientId, clientSecret)
	local.ReadState(&tokens)
	strava.AccessToken(&tokens)
	local.WriteState(&tokens)
	fmt.Print(strava.Activities(&tokens))
}
