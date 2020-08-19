package util

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	UseSSL    bool
	SentryDsn string
	Port      int
	Domain    string
	Email     string

	RequestLogging struct {
		Enable     bool
		ShowStatus bool
		ShowIP     bool
		ShowMethod bool
		ShowPath   bool
		ShowQuery  bool
	}

	CORS struct {
		Enable         bool
		AllowedOrigins []string
	}
}

func ReadConfig() Config {
	fileBytes, err := ioutil.ReadFile("config.json")
	CheckError(err)

	var config Config
	err = json.Unmarshal(fileBytes, &config)

	CheckError(err)

	return config

}
