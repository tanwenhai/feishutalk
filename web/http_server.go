package web

import (
	"github.com/buaazp/fasthttprouter"
	log "github.com/sirupsen/logrus"
	"github.com/tanwenhai/feishutalk/config"
	"github.com/valyala/fasthttp"
	"strconv"
)

func Run() {
	router := fasthttprouter.New()
	router.POST("/:name/webhook", Webhook)

	log.Info("http server listen on ", config.Server().Port)
	log.Fatal(fasthttp.ListenAndServe(":"+strconv.Itoa(config.Server().Port), router.Handler))
}
