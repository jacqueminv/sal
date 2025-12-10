package strava

import (
	"encoding/json"
	"fmt"
	"github.com/jacqueminv/sal/internal/os"
	"github.com/jacqueminv/sal/internal/timeutils"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
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

var (
	authConfig StravaAuthConfig
)

func Init(clientId string, clientSecret string) {
	authConfig.clientId = clientId
	authConfig.clientSecret = clientSecret
}

func Activities(tokens *StravaToken) string {
	after := timeutils.FiftyOneWeeksBeforeMondayThisWeek(time.Now().UTC()).Unix()
	c := http.Client{Timeout: time.Duration(60) * time.Second}
	activitiesUrl := "https://www.strava.com/api/v3/athlete/activities?after=" + strconv.FormatInt(after, 10) + "&per_page=200"
	req, err := http.NewRequest("GET", activitiesUrl, nil)
	if err != nil {
		log.Fatal("Failed to create NewRequest", activitiesUrl)
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	resp, err := c.Do(req)
	if err != nil {
		log.Fatal("Failed to perform the activities request")
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body)
}

func AccessToken(tokens *StravaToken) {
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
		os.OpenBrowser(fmt.Sprintf("https://www.strava.com/oauth/authorize?client_id=%s&redirect_uri=http://localhost:8080&response_type=code&approval_prompt=auto&scope=activity:read_all", authConfig.clientId))
		t := <-quit
		tokens.ExpiresAt = t.ExpiresAt
		tokens.RefreshToken = t.RefreshToken
		tokens.AccessToken = t.AccessToken
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

func confirmCode(code string) StravaToken {
	form := url.Values{}
	form.Add("client_id", authConfig.clientId)
	form.Add("client_secret", authConfig.clientSecret)
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)

	return token(form)
}
