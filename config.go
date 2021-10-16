package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type Config struct {
	Watchers []Watcher `json:"watchers"`
}

type Watcher struct {
	PathToLogFile string `json:"path_to_log_file"`
	ConnectionRegex string `json:"connection_regex"`
	SuccessRegex string `json:"success_regex"`
	BanCommand string `json:"ban_command"`
	UltimatumTimeInSeconds int `json:"ultimatum_time_in_seconds"`
	TrustKnown bool `json:"trust_known"`
	InstantBanAfter int `json:"instant_ban_after"`
}

func ParseConfig(configFileLocation string) (*Config, error) {
	configFile, err := ioutil.ReadFile(configFileLocation)

	if err != nil {
		return nil, errors.New("unable read config file")
	}

	config := Config{}

	err = json.Unmarshal(configFile, &config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}