package appconfig

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	WebPort int `json:"web_port"`
}

func ParseConfig() (Config, error) {

	jsonFile, err := os.Open("config.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Println(err)
		return Config{}, err
	}
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Println(err)
		return Config{}, err
	}

	conf := Config{}

	err = json.Unmarshal(byteValue, &conf)
	if err != nil {
		log.Println(err)
		return Config{}, err
	}

	return conf, nil
}
