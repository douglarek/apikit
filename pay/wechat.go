package pay

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

const wechatUnifiedorderURL = "https://api.mch.weixin.qq.com/pay/unifiedorder"

// WechatService ...
type WechatService struct {
	client *bronx.Client
}

// NewWechatService ...
func NewWechatService(httpClient *http.Client) *WechatService {
	c := bronx.NewClient(httpClient, bronx.MediaXML)
	return &WechatService{client: c}
}

// WechatUnifiedorderReq ...
type WechatUnifiedorderReq struct {
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

// WechatUnifiedorderResp ...
type WechatUnifiedorderResp struct {
	XMLName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code"`
	ReturnMsg  string   `xml:"return_msg"`
	TradeType  string   `xml:"trade_type"`
	PrepayID   string   `xml:"prepay_id"`
	CodeURL    string   `xml:"code_url"`
}

// Sign ...
func (s *WechatService) Sign(s0 interface{}, secret string) string {
	m := bronx.Params(structs.Map(s0))
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
func (s *WechatService) Order(r *WechatUnifiedorderReq) (*WechatUnifiedorderResp, error) {
	req, err := s.client.NewRequest("POST", wechatUnifiedorderURL, r)
	if err != nil {
		return nil, err
	}
	res := new(WechatUnifiedorderResp)
	if _, err := s.client.Do(req, res); err != nil {
		return nil, err
	}
	return res, nil
}
