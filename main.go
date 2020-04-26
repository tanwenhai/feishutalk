package main

import (
	"bytes"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

type ProxyHostClient struct {
	Proxy []struct {
		Name string `yaml:"name"`
		Host string `yaml:"host"`
	} `yaml:"proxy"`
}

var proxyHostClient *ProxyHostClient

func prepareRequest(ctx *fasthttp.RequestCtx) {
	req := &ctx.Request
	req.Header.Del("Connection")
	rewrite := fasthttp.NewPathSlashesStripper(1)
	newRequestURI := string(rewrite(ctx))
	if ctx.QueryArgs().Len() == 0 {
		newRequestURI += "?" + ctx.QueryArgs().String()
	}
	req.SetRequestURI(newRequestURI)
}

var client = &fasthttp.Client{}

func Proxy(ctx *fasthttp.RequestCtx) {
	req := &ctx.Request
	resp := &ctx.Response
	path := ctx.Path()
	n := bytes.IndexByte(path[1:], '/')
	errorHandler := func(format string, a ...interface{}) {
		msg := fmt.Sprintf(format, a)
		log.Printf(msg)
		resp.SetStatusCode(500)
		resp.SetBody([]byte(msg))
	}

	if n < 0 {
		ctx.Response.SetStatusCode(200)
		ctx.Response.SetBody([]byte("This Is Wechat Proxy Server"))
		return
	}
	appname := string(path[1 : n+1])

	prepareRequest(ctx)
	var host = ""
	for _, v := range proxyHostClient.Proxy {
		if v.Name == appname {
			if strings.HasPrefix(v.Host, "https") {
				host = v.Host[8:]
			} else {
				host = v.Host[7:]
			}
		}
	}
	if host == "" {
		resp.SetStatusCode(500)
		errorHandler("host no match %s", appname)
		return
	}
	req.SetHost(host)
	req.Header.SetHost(host)

	if err := client.Do(req, resp); err != nil {
		errorHandler("error when proxying the request: %s", err)
		return
	}
	if log.IsLevelEnabled(log.DebugLevel) {
		log.WithField("req", req.String()).WithField("resp", resp.String()).Debugf("track request")
	}

	postprocessResponse(resp)
}

func postprocessResponse(resp *fasthttp.Response) {
	resp.Header.Del("Connection")
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	level, err := log.ParseLevel("debug")
	if err != nil {
		log.Fatalln(err)
	}
	log.SetLevel(level)

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
	proxyHostClient = new(ProxyHostClient)
	err = yaml.Unmarshal(yamlFile, proxyHostClient)
	if err != nil {
		log.Fatalln(err)
	}

	log.Fatal(fasthttp.ListenAndServe(":8888", Proxy))
}
