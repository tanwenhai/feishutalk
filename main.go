package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/buaazp/fasthttprouter"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Webhook map[string]struct {
		Url string `yaml:"url"`
	} `yaml:"webhook"`
}

var serverConfig *ServerConfig

func Webhook(ctx *fasthttp.RequestCtx) {
	name := ctx.UserValue("name").(string)
	_, ok := serverConfig.Webhook[name]
	if !ok {
		_, _ = fmt.Fprint(ctx, name+" webhook not found")
		log.Error(name + " webhook not found")
		ctx.Response.SetStatusCode(500)
		return
	}

	url := serverConfig.Webhook[name].Url
	var body map[string]interface{}
	err := json.Unmarshal(ctx.Request.Body(), &body)
	if err != nil {
		log.Error(err)
		_, _ = fmt.Fprint(ctx, "read body error")
		ctx.Response.SetStatusCode(500)
		return
	}
	alerts, ok := body["alerts"]
	if !ok {
		ctx.Response.SetStatusCode(500)
		return
	}
	var title string
	var text string
	for _, v := range alerts.([]interface{}) {
		vMap := v.(map[string]interface{})
		labels := vMap["labels"].(map[string]interface{})
		annotations := vMap["annotations"].(map[string]interface{})
		title = labels["alertname"].(string)
		text = annotations["message"].(string)
		text += "\n严重级别:"
		_, ok := labels["severity"]
		if ok {
			text += labels["severity"].(string)
		} else {
			text += "未知"
		}
		app, ok := labels["app"]
		if ok {
			text += "\n应用:" + app.(string)
		}
		namespace, ok := labels["kubernetes_namespace"]
		if !ok {
			namespace, ok = labels["namespace"]
		}
		if ok {
			text += "\n命名空间:" + namespace.(string)
		}
		req := &fasthttp.Request{}
		req.SetRequestURI(url)
		resBody := make(map[string]string)
		resBody["title"] = title
		resBody["text"] = text
		requestBody, _ := json.Marshal(resBody)
		req.SetBody(requestBody)
		req.Header.SetContentType("application/json; charset=utf-8")
		req.Header.SetMethod("POST")
		resp := &fasthttp.Response{}
		client := &fasthttp.Client{}
		if err := client.Do(req, resp); err != nil {
			_, _ = fmt.Fprintln(ctx, "请求失败:", err.Error())
			log.Error("请求失败:", err.Error())
			return
		}
		_, _ = fmt.Fprint(ctx, string(resp.Body()))
	}

	ctx.Response.SetStatusCode(200)
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
	serverConfig = new(ServerConfig)
	err = yaml.Unmarshal(yamlFile, serverConfig)
	if err != nil {
		log.Fatalln(err)
	}

	router := fasthttprouter.New()
	router.POST("/:name/webhook", Webhook)

	log.Fatal(fasthttp.ListenAndServe(":"+strconv.Itoa(serverConfig.Server.Port), router.Handler))
}
