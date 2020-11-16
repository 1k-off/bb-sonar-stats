package main

import "errors"

type Config struct {
	BaseUrl string `toml: "baseUrl"`
	Port string `toml: "port"`
	LogLevel string `toml: "logLevel""`
	BitbucketOauthKey string `toml: "bitbucketOauthKey"`
	BitbucketOauthSecret string `toml: "bitbucketOauthSecret"`
	SonarConfigPath string `toml: "sonarConfigPath"`
	Organization string `toml:"organization"`
}

func NewConfig() *Config{
	return &Config{
		Port: "8080",
		LogLevel: "debug",
		SonarConfigPath: "sonar.json",
	}
}

func (c *Config) CheckRequiredValues() error{
	if c.Organization == "" {
		return errors.New("You must fill in organization field in config.toml.")
	} else if c.BitbucketOauthSecret == "" {
		return errors.New("You must fill in bitbucketOauthSecret field in config.toml.")
	} else if c.BitbucketOauthKey == "" {
		return errors.New("You must fill in bitbucketOauthKey field in config.toml.")
	} else if c.BaseUrl == "" {
		return errors.New("You must fill in baseUrl field in config.toml.")
	} else {
		return nil
	}
}