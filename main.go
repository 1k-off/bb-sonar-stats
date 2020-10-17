package main

import (
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"log"
)

type Context struct {
	Config               *Config
	Logger *logrus.Logger
	RepoOwner            string
	RepoName             string
	SonarServerUrl       string
	SonarProjectKey      string
	SonarToken           string
	//Per tenant meta-data and security info
	Tenants map[string]*TenantInfo
}

var (
	templatesPath = "./templates/"
	staticPath    = "./static/"
)

func main() {
	config := NewConfig()
	if _, err := toml.DecodeFile("./config.toml", config); err != nil {
		log.Fatalln(err)
	}
	c := &Context{
		Config:               config,
		Logger: logrus.New(),
		Tenants:              make(map[string]*TenantInfo),
	}
	c.ConfugureLogger()
	c.ListenAndServe()
}
