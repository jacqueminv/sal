package main

type StravaAuthConfig struct {
	clientId     string
	clientSecret string
}

var (
	authConfig StravaAuthConfig
)

func main() {
	clientId := os.Getenv("SL_CLIENT_ID")
	clientSecret := os.Getenv("SL_CLIENT_SECRET")
	if clientId == "" || clientSecret == "" {
		log.Fatal("Missing SL_CLIENT_ID, SL_CLIENT_SECRET")
		os.Exit(1)
	}
	authConfig.clientId = clientId
	authConfig.clientSecret = clientSecret

	// read locally the current Strava Tokens

	// fetch an access token (if necessary)

	// write locally the updated Strava Tokens

	// outputs the activities over the past year
}
