package alidayu

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/douglarek/bronx"
	"github.com/fatih/structs"
)

const url = "https://eco.taobao.com/router/rest"

//
const (
	HMAC = "hmac"
	MD5  = "md5"
)

// Alidayu handles communication with related methods of the Alidayu API.
type Alidayu struct {
	client *bronx.Client
}

// New ...
func New(httpClient *http.Client, contentType ...string) *Alidayu {
	c := bronx.NewClient(httpClient, contentType...)
	return &Alidayu{client: c}
}

// Req ...
type Req struct {
	Method     string `json:"method,omitempty" url:"method" structs:"method"`
	AppKey     string `json:"app_key,omitempty" url:"app_key" structs:"app_key"`
	Timestamp  string `json:"timestamp,omitempty" url:"timestamp" structs:"timestamp"`
	Format     string `json:"format,omitempty" url:"format" structs:"format"`
	Version    string `json:"v,omitempty" url:"v" structs:"v"`
	PartnerID  string `json:"partner_id,omitempty" url:"partner_id" structs:"partner_id"`
	SignMethod string `json:"sign_method,omitempty" url:"sign_method" structs:"sign_method"`
	Sign       string `json:"sign,omitempty" url:"sign" structs:"sign"`
}

// SmsReq ...
type SmsReq struct {
	Req
	Extend          string `json:"extend,omitempty" url:"extend" structs:"extend"`
	SmsType         string `json:"sms_type,omitempty" url:"sms_type" structs:"sms_type"`
	SmsFreeSignName string `json:"sms_free_sign_name,omitempty" url:"sms_free_sign_name" structs:"sms_free_sign_name"`
	SmsParam        string `json:"sms_param,omitempty" url:"sms_param" structs:"sms_param"`
	RecNum          string `json:"rec_num,omitempty" url:"rec_num" structs:"rec_num"`
	SmsTemplateCode string `json:"sms_template_code,omitempty" url:"sms_template_code" structs:"sms_template_code"`
}

// DefaultSmsReq ...
func DefaultSmsReq() *SmsReq {
	req := Req{
		Format:     "json",
		Method:     "alibaba.aliqin.fc.sms.num.send",
		SignMethod: "md5",
		Timestamp:  time.Now().UTC().Add(time.Duration(8 * time.Hour)).Format("2006-01-02 15:04:05"),
		Version:    "2.0",
		PartnerID:  "apidoc",
	}
	return &SmsReq{Req: req, Extend: "123456"}
}

// Sign signs an Alidayu request struct.
func (a *Alidayu) Sign(s interface{}, secret, method string) string {
	m := bronx.Params(structs.Map(s))
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

	return encrypt(buf.Bytes(), []byte(secret), method)
}

func encrypt(s, secret []byte, method string) (h string) {
	switch method {
	case MD5:
		d := make([]byte, 0, len(s)+2*len(secret))
		d = append(d, secret...)
		d = append(d, s...)
		d = append(d, secret...)
		h = fmt.Sprintf("%X", md5.Sum(d))
	case HMAC:
		// TODO:
		mac := hmac.New(md5.New, secret)
		mac.Write(s)
		h = fmt.Sprintf("%X", mac.Sum(nil))
	default:
		panic("unsupported sign method!")
	}
	return h
}

// Resp ...
type Resp struct {
	Result struct {
		ErrorCode string `json:"err_code,omitempty"`
		Model     string `json:"model,omitempty"`
		Success   bool   `json:"success,omitempty"`
		Msg       string `json:"msg,omitempty"`
	} `json:"result"`
}

// ErrResp ...
type ErrResp struct {
	Code    int    `json:"code,omitempty"`
	Msg     string `json:"msg,omitempty"`
	SubCode string `json:"sub_code,omitempty"`
	SubMsg  string `json:"sub_msg,omitempty"`
}

// SmsResp ...
type SmsResp struct {
	Resp    `json:"alibaba_aliqin_fc_sms_num_send_response,omitempty"`
	ErrResp `json:"error_response,omitempty"`
}

// SendSms sends a sms.
func (a *Alidayu) SendSms(r *SmsReq) (*SmsResp, error) {
	req, err := a.client.NewRequest("POST", url, r)
	if err != nil {
		return nil, err
	}
	res := new(SmsResp)
	if _, err := a.client.Do(req, res); err != nil {
		return nil, err
	}
	return res, nil
}
