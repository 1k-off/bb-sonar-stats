package main

type Config struct {
	BaseUrl string `toml: "baseUrl"`
	Port string `toml: "port"`
	LogLevel string `toml: "logLevel""`
	BitbucketOauthKey string `toml: "bitbucketOauthKey"`
	BitbucketOauthSecret string `toml: "bitbucketOauthSecret"`
	RepoBranch string `toml: "repoBranch"`
	SonarConfigPath string `toml: "sonarConfigPath"`
}

func NewConfig() *Config{
	return &Config{
		Port: "8080",
		LogLevel: "debug",
		RepoBranch: "main",
		SonarConfigPath: "sonar.json",
	}
}
