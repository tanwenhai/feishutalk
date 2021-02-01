package web

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tanwenhai/feishutalk/config"
	"github.com/valyala/fasthttp"
)

type Content struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type FeiShuBotMsgTypeOfCardConfig struct {
	WideScreenMode bool `json:"wide_screen_mode"`
	EnableForWard  bool `json:"enable_for_ward"`
}

type FeiShuBotMsgTypeOfCardElements struct {
	Tag  string  `json:"tag"`
	Text Content `json:"text"`
}

type FeiShuBotMsgTypeOfCardHeader struct {
	Title Content `json:"title"`
}

type FeiShuBotMsgTypeOfCard struct {
	Config   FeiShuBotMsgTypeOfCardConfig     `json:"config"`
	Elements []FeiShuBotMsgTypeOfCardElements `json:"elements"`
	Header   FeiShuBotMsgTypeOfCardHeader
}

type FeiShuBotMsgBody struct {
	MsgType string                 `json:"msg_type"`
	Card    FeiShuBotMsgTypeOfCard `json:"card"`
}

func errorHandle(ctx *fasthttp.RequestCtx, err error) {
	_, _ = fmt.Fprint(ctx, err)
	log.Error("errorHandle ", err)
	ctx.Response.SetStatusCode(500)
}

func Webhook(ctx *fasthttp.RequestCtx) {
	name := ctx.UserValue("name").(string)
	webhookProperty := config.Webhook()
	_, ok := webhookProperty[name]
	if !ok {
		errorHandle(ctx, errors.New(name+" webhook not found"))
		return
	}

	url := webhookProperty[name].Url
	requestData := PrometheusWebHookRequestData{}
	err := json.Unmarshal(ctx.Request.Body(), &requestData)
	log.Info("webhook requestData ", string(ctx.Request.Body()))
	if err != nil {
		errorHandle(ctx, err)
		return
	}
	var title string
	var text string
	for _, v := range requestData.Alerts {
		if v.Status != Firing {
			continue
		}
		var cardElements []FeiShuBotMsgTypeOfCardElements
		labels := v.Labels
		annotations := v.Annotations
		title = labels[Alertname].(string)
		text = annotations[Message].(string)
		cardElements = append(cardElements, FeiShuBotMsgTypeOfCardElements{
			Tag: "div",
			Text: Content{
				Content: "**发生时间**:" + v.StartAt,
				Tag:     "lark_md",
			},
		})
		cardElements = append(cardElements, FeiShuBotMsgTypeOfCardElements{
			Tag: "div",
			Text: Content{
				Content: text,
				Tag:     "lark_md",
			},
		})
		serverityText := "**严重性**:"
		_, ok := labels[Serverity]
		if ok {
			serverityText += labels[Serverity].(string)
		} else {
			serverityText += "一般"
		}
		cardElements = append(cardElements, FeiShuBotMsgTypeOfCardElements{
			Tag: "div",
			Text: Content{
				Content: serverityText,
				Tag:     "lark_md",
			},
		})
		app, ok := labels["app"]
		if ok {
			cardElements = append(cardElements, FeiShuBotMsgTypeOfCardElements{
				Tag: "div",
				Text: Content{
					Content: "**应用**:" + app.(string),
					Tag:     "lark_md",
				},
			})
		}
		namespace, ok := labels["kubernetes_namespace"]
		if !ok {
			namespace, ok = labels["namespace"]
		}
		if ok {
			cardElements = append(cardElements, FeiShuBotMsgTypeOfCardElements{
				Tag: "div",
				Text: Content{
					Content: "**命名空间**:" + namespace.(string),
					Tag:     "lark_md",
				},
			})
		}
		req := &fasthttp.Request{}
		req.SetRequestURI(url)

		var feishuBodyData = FeiShuBotMsgBody{
			MsgType: "interactive",
			Card: FeiShuBotMsgTypeOfCard{
				Config: FeiShuBotMsgTypeOfCardConfig{
					WideScreenMode: true,
					EnableForWard:  true,
				},
				Elements: cardElements,
				Header: FeiShuBotMsgTypeOfCardHeader{
					Title: Content{
						Content: title,
						Tag:     "plain_text",
					},
				},
			},
		}
		outputBody, err := json.Marshal(feishuBodyData)
		if err != nil {
			errorHandle(ctx, err)
			return
		}
		req.SetBody(outputBody)
		req.Header.SetContentType("application/json; charset=utf-8")
		req.Header.SetMethod("POST")
		resp := &fasthttp.Response{}
		client := &fasthttp.Client{}
		if err := client.Do(req, resp); err != nil {
			errorHandle(ctx, err)
			return
		}
		respBody := string(resp.Body())
		log.WithField("respBody", respBody).Info("respBody")
	}

	ctx.Response.SetStatusCode(200)
}

//sum(ira
//te(istio_requests_total{reporter=\"$qrep\",destination_service=~\"$service\",response_code!~\"5.*\"}[5m])) / sum(irate(istio_requests_total{reporter=\"$qrep\",destination_service=~\"$service\"}[5m]))
