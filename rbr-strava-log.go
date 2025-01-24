package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type StravaAuthConfig struct {
	clientId     string
	clientSecret string
}

type StravaToken struct {
	ExpiresAt    int64  `json:"expires_at"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

var (
	authConfig StravaAuthConfig
)

func configPath() string {
	var configPath string

	dir, dirErr := os.UserConfigDir()
	if dirErr == nil {
		configPath = filepath.Join(dir, "rbr-strava", "state.json")
		var err error
		_, err = os.ReadFile(configPath)
		if err != nil && !os.IsNotExist(err) {
			log.Fatal(err)
			panic(err)
		} else if err != nil && os.IsNotExist(err) {
			err = os.MkdirAll(filepath.Dir(configPath), 0700)
			if err != nil {
				log.Fatal(err)
				panic(err)
			}
		}
	} else {
		log.Fatal("Failed to get UserConfigDir")
		panic(dirErr)
	}

	return configPath
}

func readState(tokens *StravaToken) {
	var (
		configPath string = configPath()
		origConfig []byte
		err        error
	)

	origConfig, err = os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
		panic(err)
	} else if len(origConfig) == 0 {
		origConfig = []byte("{}")
		err = os.WriteFile(configPath, origConfig, 0600)
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
	}

	if err := json.Unmarshal(origConfig, &tokens); err != nil {
		panic(err)
	}
}

func main() {
	clientId := os.Getenv("SL_CLIENT_ID")
	clientSecret := os.Getenv("SL_CLIENT_SECRET")
	if clientId == "" || clientSecret == "" {
		log.Fatal("Missing SL_CLIENT_ID, SL_CLIENT_SECRET")
		os.Exit(1)
	}
	authConfig.clientId = clientId
	authConfig.clientSecret = clientSecret

	var tokens StravaToken
	readState(&tokens)

	// fetch an access token (if necessary)

	// write locally the updated Strava Tokens

	// outputs the activities over the past year
}
