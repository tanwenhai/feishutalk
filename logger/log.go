package logger

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func Init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	level, err := log.ParseLevel("debug")
	if err != nil {
		log.Fatalln(err)
	}
	log.SetLevel(level)
}
