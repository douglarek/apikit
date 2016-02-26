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

const (
	orderURL = "https://api.mch.weixin.qq.com/pay/unifiedorder"
	queryURL = "https://api.mch.weixin.qq.com/pay/orderquery"
)

// Wechat ...
type Wechat struct {
	client *bronx.Client
}

// New makes a wechat ...
func New(httpClient *http.Client) *Wechat {
	c := bronx.NewClient(httpClient, bronx.MediaXML)
	return &Wechat{client: c}
}

// Req ...
type Req struct {
	XMLName  xml.Name `xml:"xml"`
	AppID    string   `xml:"appid" structs:"appid"`
	MchID    string   `xml:"mch_id" structs:"mch_id"`
	NonceStr string   `xml:"nonce_str" structs:"nonce_str"`
	Sign     string   `xml:"sign" structs:"sign"`
}

// OrderReq ...
type OrderReq struct {
	Req
	DeviceInfo     string `xml:"device_info" structs:"device_info"`
	Body           string `xml:"body" structs:"body"`
	Detail         string `xml:"detail" structs:"detail"`
	Attach         string `xml:"attach" structs:"attach"`
	OutTradeNo     string `xml:"out_trade_no" structs:"out_trade_no"`
	FeeType        string `xml:"fee_type" structs:"fee_type"`
	TotalFee       int    `xml:"total_fee" structs:"total_fee"`
	SpbillCreateIP string `xml:"spbill_create_ip" structs:"spbill_create_ip"`
	TimeStart      string `xml:"time_start" structs:"time_start"`
	TimeExpire     string `xml:"time_expire" structs:"time_expire"`
	GoodsTag       string `xml:"goods_tag" structs:"goods_tag"`
	NotifyURL      string `xml:"notify_url" structs:"notify_url"`
	TradeType      string `xml:"trade_type" structs:"trade_type"`
	ProductID      string `xml:"product_id" structs:"product_id"`
	LimitPay       string `xml:"limit_pay" structs:"limit_pay"`
	OpenID         string `xml:"openid" structs:"openid"`
}

// Resp ...
type Resp struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`
}

// OrderResp ...
type OrderResp struct {
	Resp
	Req
	DeviceInfo string `xml:"device_info" structs:"device_info"`
	TradeType  string `xml:"trade_type"`
	PrepayID   string `xml:"prepay_id"`
	CodeURL    string `xml:"code_url"`
}

// QueryReq ...
type QueryReq struct {
	Req
	TransactionID string `xml:"transaction_id" structs:"transaction_id"`
	OutTradeNo    string `xml:"out_trade_no" structs:"out_trade_no"`
}

// QueryResp ...
type QueryResp struct {
	Resp
	Req
	DeviceInfo  string `xml:"device_info"`
	OpenID      string `xml:"openid"`
	IsSubscribe string `xml:"is_subscribe"`
	TradeType   string `xml:"trade_type"`
	TradeState  string `xml:"trade_state" structs:"trade_state"`
	BankType    string `xml:"bank_type"`
	TotalFee    string `xml:"total_fee"`
	FeeType     string `xml:"fee_type"`
	CashFee     int    `xml:"cash_fee"`
	CashFeeType string `xml:"cash_fee_type"`
	CouponFee   int    `xml:"coupon_fee"`
	CouponCount int    `xml:"coupon_count"`
	// coupon_batch_id_$n, coupon_id_$n, coupon_fee_$n
	TransactionID  string `xml:"transaction_id"`
	OutTradeNo     string `xml:"out_trade_no"`
	Attach         string `xml:"attach"`
	TimeEnd        string `xml:"time_end"`
	TradeStateDesc string `xml:"trade_state_desc"`
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

// Query ...
func (w *Wechat) Query(r *QueryReq) (*QueryResp, error) {
	req, err := w.client.NewRequest("POST", queryURL, r)
	if err != nil {
		return nil, err
	}
	res := new(QueryResp)
	if _, err := w.client.Do(req, res); err != nil {
		return nil, err
	}
	return res, nil
}
