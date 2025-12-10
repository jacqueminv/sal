package local

import (
	"encoding/json"
	"github.com/jacqueminv/sal/internal/strava"
	"log"
	"os"
	"path/filepath"
)

func ReadState(tokens *strava.StravaToken) {
	var (
		configPath string = configPath()
		origConfig []byte
		err        error
	)

	origConfig, err = os.ReadFile(configPath)
	if err != nil && !os.IsNotExist(err) {
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

func WriteState(tokens *strava.StravaToken) {
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

func configPath() string {
	var configPath string

	dir, dirErr := os.UserConfigDir()
	if dirErr == nil {
		configPath = filepath.Join(dir, "sal", "state.json")
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
