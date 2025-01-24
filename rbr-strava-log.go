package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
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

const (
	fiftyWeeks = time.Hour * 24 * 7 * 50
)

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
func writeState(tokens *StravaToken) {
	var configPath string = configPath()

	content, err := json.Marshal(tokens)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(configPath, content, 0600)
	if err != nil {
		log.Printf("error saving config changes: %v", err)
		panic(err)
	}
}

func token(form url.Values) StravaToken {
	url := "https://www.strava.com/oauth/token"

	resp, err := http.PostForm(url, form)
	if err != nil || resp.StatusCode != 200 {
		log.Fatal("Failed to interact with the token api")
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var response StravaToken
	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatal("Failed to unmarshal Strava Response")
		panic(err)
	}

	return response
}
func confirmCode(code string) StravaToken {
	form := url.Values{}
	form.Add("client_id", authConfig.clientId)
	form.Add("client_secret", authConfig.clientSecret)
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)

	return token(form)
}
func refresh(tokens *StravaToken) {
	form := url.Values{}
	form.Add("client_id", authConfig.clientId)
	form.Add("client_secret", authConfig.clientSecret)
	form.Add("grant_type", "refresh_token")
	form.Add("refresh_token", tokens.RefreshToken)

	t := token(form)
	tokens.RefreshToken = t.RefreshToken
	tokens.AccessToken = t.AccessToken
	tokens.ExpiresAt = t.ExpiresAt
}
func accessToken(tokens *StravaToken) {
	var now = time.Now().UTC()
	if tokens.ExpiresAt < now.Unix() && len(tokens.RefreshToken) > 0 {
		refresh(tokens)
	}

	if tokens.ExpiresAt < now.Unix() {
		quit := make(chan StravaToken)
		acceptsTokenHandler := func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/" {
				q := req.URL.Query()
				code := q.Get("code")
				io.WriteString(w, "<!doctype html><html><body><p>You can close this page, the Strava Authorization proceeded successfully.</p></body>")
				resp := confirmCode(code)
				quit <- resp
			}
		}
		go func() {
			http.HandleFunc("/", acceptsTokenHandler)
			log.Fatal(http.ListenAndServe(":8080", nil))
		}()
		t := <-quit
		tokens.ExpiresAt = t.ExpiresAt
		tokens.RefreshToken = t.RefreshToken
		tokens.AccessToken = t.AccessToken
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
	accessToken(&tokens)
	writeState(&tokens)

	// outputs the activities over the past year
}
