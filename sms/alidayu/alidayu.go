package alidayu

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/douglarek/apikit"
	"github.com/fatih/structs"
	"github.com/imdario/mergo"
)

const url = "https://eco.taobao.com/router/rest"

// Alidayu handles communication with related methods of the Alidayu API.
type Alidayu struct {
	client *apikit.Client
}

// New ...
func New(httpClient *http.Client) *Alidayu {
	c := apikit.NewClient(httpClient)
	return &Alidayu{client: c}
}

// Config the alidayu configuration.
type Config struct {
	Method          string `json:"method,omitempty" structs:"method"`
	AppKey          string `json:"app_key,omitempty" structs:"app_key"`
	Timestamp       string `json:"timestamp,omitempty" structs:"timestamp"`
	Format          string `json:"format,omitempty" structs:"format"`
	Version         string `json:"v,omitempty" structs:"v"`
	PartnerID       string `json:"partner_id,omitempty" structs:"partner_id"`
	SignMethod      string `json:"sign_method,omitempty" structs:"sign_method"`
	Sign            string `json:"sign,omitempty" structs:"sign"`
	Extend          string `json:"extend,omitempty" structs:"extend"`
	SmsType         string `json:"sms_type,omitempty" structs:"sms_type"`
	SmsFreeSignName string `json:"sms_free_sign_name,omitempty" structs:"sms_free_sign_name"`
	SmsParam        string `json:"sms_param,omitempty" structs:"sms_param"`
	RecNum          string `json:"rec_num,omitempty" structs:"rec_num"`
	SmsTemplateCode string `json:"sms_template_code,omitempty" structs:"sms_template_code"`
}

// DefaultConfig returns the default alidayu configuration.
func DefaultConfig() Config {
	localLoc, _ := time.LoadLocation("Asia/Chongqing")
	return Config{
		Format:     "json",
		Method:     "alibaba.aliqin.fc.sms.num.send",
		SignMethod: "md5",
		Timestamp:  time.Now().In(localLoc).Format("2006-01-02 15:04:05"),
		Version:    "2.0",
		PartnerID:  "apidoc",
		Extend:     "123456",
		SmsType:    "normal"}
}

// Merge merges the default with the given config and returns the result.
func (c Config) Merge(cfg Config) (config Config) {
	config = cfg
	mergo.Merge(&config, c)
	return
}

// Sign signs an Alidayu request struct.
func (a *Alidayu) Sign(s interface{}, secret []byte) string {
	m := apikit.Params(structs.Map(s))
	delete(m, "sign")
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf bytes.Buffer
	for _, k := range keys {
		buf.WriteString(k)
		buf.WriteString(m[k])
	}

	return encrypt(buf.Bytes(), secret)
}

func encrypt(s, secret []byte) (h string) {
	d := make([]byte, 0, len(s)+2*len(secret))
	d = append(d, secret...)
	d = append(d, s...)
	d = append(d, secret...)
	return fmt.Sprintf("%X", md5.Sum(d))
}

// SendSms sends a sms.
func (a *Alidayu) SendSms(c Config) (map[string]interface{}, error) {
	req, err := a.client.NewRequest("POST", url, &c)
	if err != nil {
		return nil, err
	}
	res := map[string]interface{}{}
	if _, err := a.client.Do(req, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// SmsResult judges a sms sent ok or not.
func SmsResult(m map[string]interface{}) (bool, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return false, err
	}
	j, err := simplejson.NewJson(b)
	if err != nil {
		return false, err
	}
	return j.Get("alibaba_aliqin_fc_sms_num_send_response").Get("result").Get("success").Bool()
}
