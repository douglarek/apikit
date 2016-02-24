package sms

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"fmt"
	"net/http"
	"reflect"
	"sort"

	"github.com/douglarek/bronx"
	"github.com/fatih/structs"
)

// AlidayuURL ...
const AlidayuURL = "https://eco.taobao.com/router/rest"

// AlidayuService handles communication with related methods of the Alidayu API.
type AlidayuService struct {
	client *bronx.Client
}

// NewAlidayuService ...
func NewAlidayuService(httpClient *http.Client, contentType ...string) *AlidayuService {
	c := bronx.NewClient(httpClient, contentType...)
	return &AlidayuService{client: c}
}

// AlidayuRequest ...
type AlidayuRequest struct {
	Method     string `json:"method,omitempty" url:"method" structs:"method"`
	AppKey     string `json:"app_key,omitempty" url:"app_key" structs:"app_key"`
	Timestamp  string `json:"timestamp,omitempty" url:"timestamp" structs:"timestamp"`
	Format     string `json:"format,omitempty" url:"format" structs:"format"`
	Version    string `json:"v,omitempty" url:"v" structs:"v"`
	PartnerID  string `json:"partner_id,omitempty" url:"partner_id" structs:"partner_id"`
	SignMethod string `json:"sign_method,omitempty" url:"sign_method" structs:"sign_method"`
	Sign       string `json:"sign,omitempty" url:"sign" structs:"sign"`
}

// AlidayuSmsRequest ...
type AlidayuSmsRequest struct {
	AlidayuRequest
	Extend          string `json:"extend,omitempty" url:"extend" structs:"extend"`
	SmsType         string `json:"sms_type,omitempty" url:"sms_type" structs:"sms_type"`
	SmsFreeSignName string `json:"sms_free_sign_name,omitempty" url:"sms_free_sign_name" structs:"sms_free_sign_name"`
	SmsParam        string `json:"sms_param,omitempty" url:"sms_param" structs:"sms_param"`
	RecNum          string `json:"rec_num,omitempty" url:"rec_num" structs:"rec_num"`
	SmsTemplateCode string `json:"sms_template_code,omitempty" url:"sms_template_code" structs:"sms_template_code"`
}

func params(m0 map[string]interface{}) (m map[string]string) {
	if m == nil {
		m = make(map[string]string)
	}

	for k, v := range m0 {
		val := reflect.ValueOf(v)
		if k == "sign" || v == nil {
			continue
		}
		switch val.Kind() {
		case reflect.Map:
			for k, v0 := range params(v.(map[string]interface{})) {
				m[k] = v0
			}
		case reflect.String:
			m[k] = v.(string)
		default:
			panic(fmt.Sprintf("unsupported type: %T", v))
		}
	}
	return
}

// Sign signs an Alidayu request struct.
func (a *AlidayuService) Sign(s interface{}, secret, method string) string {
	m := params(structs.Map(s))
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
	case bronx.MD5:
		d := make([]byte, 0, len(s)+2*len(secret))
		d = append(d, secret...)
		d = append(d, s...)
		d = append(d, secret...)
		h = fmt.Sprintf("%X", md5.Sum(d))
	case bronx.HMAC:
		// TODO:
		mac := hmac.New(md5.New, secret)
		mac.Write(s)
		h = fmt.Sprintf("%X", mac.Sum(nil))
	default:
		panic("unsupported sign method!")
	}
	return h
}

// AlidayuResponse ...
type AlidayuResponse struct {
	Result struct {
		ErrorCode string `json:"err_code,omitempty"`
		Model     string `json:"model,omitempty"`
		Success   bool   `json:"success,omitempty"`
		Msg       string `json:"msg,omitempty"`
	} `json:"result"`
}

// AlidayuErrorResponse ...
type AlidayuErrorResponse struct {
	Code    int    `json:"code,omitempty"`
	Msg     string `json:"msg,omitempty"`
	SubCode string `json:"sub_code,omitempty"`
	SubMsg  string `json:"sub_msg,omitempty"`
}

// AlidayuSmsResponse ...
type AlidayuSmsResponse struct {
	AlidayuResponse      `json:"alibaba_aliqin_fc_sms_num_send_response,omitempty"`
	AlidayuErrorResponse `json:"error_response,omitempty"`
}

// SendSms sends a sms.
func (a *AlidayuService) SendSms(r *AlidayuSmsRequest) (*AlidayuSmsResponse, error) {
	req, err := a.client.NewRequest("POST", AlidayuURL, r)
	if err != nil {
		return nil, err
	}
	res := new(AlidayuSmsResponse)
	if _, err := a.client.Do(req, res); err != nil {
		return nil, err
	}
	return res, nil
}
