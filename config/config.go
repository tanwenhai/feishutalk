package config

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var serverConfig *ServerConfig

type ServerProperty struct {
	Port int `yaml:"port"`
}

type WebhookProperty map[string]struct {
	Url string `yaml:"url"`
}

type ServerConfig struct {
	Server  ServerProperty  `yaml:"server"`
	Webhook WebhookProperty `yaml:"webhook"`
}

func Init() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	var configFile = flag.String("config", pwd+"/config.yaml", "config file")
	flag.Parse()

	yamlFile, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatalln(err)
	}
	serverConfig = new(ServerConfig)
	err = yaml.Unmarshal(yamlFile, serverConfig)
	if err != nil {
		log.Fatalln(err)
	}
}

func Server() ServerProperty {
	return serverConfig.Server
}

func Webhook() WebhookProperty {
	return serverConfig.Webhook
}
