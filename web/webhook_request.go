package web

//{
//"receiver": "feishu",
//"status": "firing",
//"alerts": [
//{
//"status": "firing",
//"labels": {
//"alertname": "test",
//"prometheus": "kubesphere-monitoring-system/k8s",
//"serverity": "warning"
//},
//"annotations": {
//"message": "test"
//},
//"startsAt": "2021-02-01T03:37:13.293116502Z",
//"endsAt": "0001-01-01T00:00:00Z",
//"generatorURL": "http://prometheus-k8s-0:9090/graph?g0.expr=vector%281%29&g0.tab=1",
//"fingerprint": "eacdcf4af175c687"
//}
//],
//"groupLabels": {
//"alertname": "test"
//},
//"commonLabels": {
//"alertname": "test",
//"prometheus": "kubesphere-monitoring-system/k8s",
//"serverity": "warning"
//},
//"commonAnnotations": {
//"message": "test"
//},
//"externalURL": "http://alertmanager-main-2:9093",
//"version": "4",
//"groupKey": "{}/{alerttype=\"\"}:{alertname=\"test\"}",
//"truncatedAlerts": 0
//}
const (
	Firing    = "firing"
	Resolved  = "resolved"
	Alertname = "alertname"
	Serverity = "Serverity"
	Message   = "message"
)

type PrometheusWebHookRequestDataOfAlert struct {
	Status       string                 `json:"status"`
	Labels       map[string]interface{} `json:"labels"`
	Annotations  map[string]interface{} `json:"annotations"`
	StartAt      string                 `json:"startsAt"`
	EndsAt       string                 `json:"endsAt"`
	GeneratorURL string                 `json:"generatorURL"`
	Fingerprint  string                 `json:"fingerprint"`
}

type PrometheusWebHookRequestData struct {
	Receiver          string                                `json:"receiver"`
	Status            string                                `json:"status"`
	Alerts            []PrometheusWebHookRequestDataOfAlert `json:"alerts"`
	GroupLabels       map[string]interface{}                `json:"groupLabels"`
	CommonLabels      map[string]interface{}                `json:"commonLabels"`
	CommonAnnotations map[string]interface{}                `json:"commonAnnotations"`
	ExternalURL       string                                `json:"externalURL"`
	Version           string                                `json:"version"`
	GroupKey          string                                `json:"groupKey"`
	TruncatedAlerts   int                                   `json:"truncatedAlerts"`
}
