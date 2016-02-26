package wechat

import (
	"bytes"
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"net/http"
	"sort"

	"github.com/douglarek/bronx"
	"github.com/fatih/structs"
)

const orderURL = "https://api.mch.weixin.qq.com/pay/unifiedorder"

// Wechat ...
type Wechat struct {
	client *bronx.Client
}

// New makes a wechat ...
func New(httpClient *http.Client) *Wechat {
	c := bronx.NewClient(httpClient, bronx.MediaXML)
	return &Wechat{client: c}
}

// OrderReq ...
type OrderReq struct {
	XMLName        xml.Name `xml:"xml"`
	AppID          string   `xml:"appid" structs:"appid"`
	MchID          string   `xml:"mch_id" structs:"mch_id"`
	NonceStr       string   `xml:"nonce_str" structs:"nonce_str"`
	Sign           string   `xml:"sign" structs:"sign"`
	Body           string   `xml:"body" structs:"body"`
	OutTradeNo     string   `xml:"out_trade_no" structs:"out_trade_no"`
	TotalFee       int      `xml:"total_fee" structs:"total_fee"`
	SpbillCreateIP string   `xml:"spbill_create_ip" structs:"spbill_create_ip"`
	NotifyURL      string   `xml:"notify_url" structs:"notify_url"`
	TradeType      string   `xml:"trade_type" structs:"trade_type"`
}

// OrderResp ...
type OrderResp struct {
	XMLName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code"`
	ReturnMsg  string   `xml:"return_msg"`
	TradeType  string   `xml:"trade_type"`
	PrepayID   string   `xml:"prepay_id"`
	CodeURL    string   `xml:"code_url"`
}

// Sign ...
func (w *Wechat) Sign(s interface{}, secret string) string {
	m := bronx.Params(structs.Map(s))
	delete(m, "sign")
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf bytes.Buffer
	for _, k := range keys {
		buf.WriteString(fmt.Sprintf("%s=%s", k, m[k]))
		buf.WriteString("&")
	}
	buf.WriteString("key=" + secret)
	return fmt.Sprintf("%X", md5.Sum(buf.Bytes()))
}

// Order ...
func (w *Wechat) Order(r *OrderReq) (*OrderResp, error) {
	req, err := w.client.NewRequest("POST", orderURL, r)
	if err != nil {
		return nil, err
	}
	res := new(OrderResp)
	if _, err := w.client.Do(req, res); err != nil {
		return nil, err
	}
	return res, nil
}
