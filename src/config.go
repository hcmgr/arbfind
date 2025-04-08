package main

import (
	"encoding/json"
	"os"
	"strconv"
)

var configFilePath string = "config.json"

type Config struct {
	// fields read directly from config file
	BaseURL                string            `json:"baseURL"`
	DefaultParams          DefaultOddsParams `json:"defaultOddsParams"`
	BackendPort            int               `json:"backendPort"`
	OutputDir              string            `json:"outputDir"`
	ApiKeyIndex            int               `json:"apiKeyIndex"`
	ApiKeys                []ApiKey          `json:"apiKeys"`
	BookmakerUrls          map[string]string `json:"bookmakerUrls"`
	MongodbUriPrefix       string            `json:"mongodbUriPrefix"`
	MongodbUriPrefixDocker string            `json:"mongodbUriPrefixDocker"`
	MongoDbPort            int               `json:"mongodbPort"`

	// fields populated in post-processing
	MongoDbUri string
}

type ApiKey struct {
	Email  string `json:"email"`
	ApiKey string `json:"apiKey"`
}

type DefaultOddsParams struct {
	Regions    string `json:"regions"`
	Markets    string `json:"markets"`
	OddsFormat string `json:"oddsFormat"`
}

func initConfig(cliArgs *CliArgs) {
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		panic(err)
	}

	config = &Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}

	config.postProcess(cliArgs)
}

func (config *Config) postProcess(cliArgs *CliArgs) {
	// resolve mongo db uri
	var mongodbUriPrefix string
	if cliArgs.runningInDocker {
		mongodbUriPrefix = config.MongodbUriPrefixDocker
	} else {
		mongodbUriPrefix = config.MongodbUriPrefix
	}
	config.MongoDbUri = mongodbUriPrefix + ":" + strconv.Itoa(config.MongoDbPort)
}

func (config *Config) findApiKey() string {
	return config.ApiKeys[config.ApiKeyIndex].ApiKey
}

func (config *Config) findBookmakerUrl(bookmakerKey string) string {
	if url, exists := config.BookmakerUrls[bookmakerKey]; exists {
		return url
	}
	return ""
}
