package ali

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strings"

	"github.com/douglarek/bronx"
	"github.com/fatih/structs"
)

const orderURL = "https://mapi.alipay.com/gateway.do"

// sign methods
const (
	MD5 = "MD5"
	RSA = "RSA"
)

// Ali ...
type Ali struct {
	client *bronx.Client
}

// New makes an ali ...
func New(httpClient *http.Client) *Ali {
	c := bronx.NewClient(httpClient)
	return &Ali{client: c}
}

// Req ...
type Req struct {
	Service      string `structs:"service" json:"service"`
	Partner      string `structs:"partner" json:"partner"`
	InputCharset string `structs:"_input_charset" json:"inputCharset"`
	SignType     string `structs:"sign_type" json:"signType"`
	Sign         string `structs:"sign" json:"sign"`
	NotifyURL    string `structs:"notify_url" json:"notifyUrl"`
	ReturnURL    string `structs:"return_url" json:"returnUrl"`
}

// OrderReq ...
type OrderReq struct {
	Req
	OutTradeNo        string `structs:"out_trade_no" json:"outTradeNo"`
	Subject           string `structs:"subject" json:"subject"`
	PaymentType       string `structs:"payment_type" json:"paymentType"`
	TotalFee          string `structs:"total_fee" json:"totalFee"`
	SellerID          string `structs:"seller_id" json:"sellerId"`
	SellerEmail       string `structs:"seller_email" json:"sellerEmail"`
	SellerAccountName string `structs:"seller_account_name" json:"sellerAccountName"`
	ItBPay            string `structs:"it_b_pay" json:"itBPay"`
	Body              string `structs:"body" json:"body"`
}

func sortedParams(m map[string]string) bytes.Buffer {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf bytes.Buffer
	for i, k := range keys {
		if i > 0 {
			buf.WriteString("&")
		}
		buf.WriteString(fmt.Sprintf("%s=%q", k, m[k]))
	}
	return buf
}

func removeQuote(b []byte) []byte {
	return bytes.Replace(b, []byte(`"`), []byte(``), -1)
}

func removeKeys(m map[string]string, keys ...string) map[string]string {
	for _, k := range keys {
		if _, ok := m[k]; ok {
			delete(m, k)
		}
	}
	return m
}

// Sign ...
func (a *Ali) Sign(s interface{}, secretKey []byte) (b []byte) {
	m := bronx.Params(structs.Map(s))
	st := m["sign_type"]
	buf := sortedParams(removeKeys(m, "sign", "sign_type"))
	switch st {
	case RSA:
		p, _ := pem.Decode([]byte(secretKey))
		if p == nil {
			panic("Secret key broken!")
		}
		key, err := x509.ParsePKCS8PrivateKey(p.Bytes)
		if err != nil {
			panic(err)
		}
		h := crypto.Hash.New(crypto.SHA1)
		h.Write(buf.Bytes())
		sum := h.Sum(nil)
		sig, _ := rsa.SignPKCS1v15(rand.Reader, key.(*rsa.PrivateKey), crypto.SHA1, sum)
		return []byte(base64.StdEncoding.EncodeToString(sig))
	case MD5:
		buf.WriteString(string(secretKey))
		h := crypto.Hash.New(crypto.MD5)
		h.Write(removeQuote(buf.Bytes()))
		return h.Sum(nil)
	}
	return
}

// VerifyNotifyID ...
func (a *Ali) VerifyNotifyID(partner, notifyID string) bool {
	q := fmt.Sprintf("service=notify_verify&partner=%s&notify_id=%s", partner, notifyID)
	resp, err := http.Get(strings.Join([]string{orderURL, "?", q}, ""))
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	return err != nil && reflect.DeepEqual(b, []byte(`true`))
}

// Verify for RSA sign.
func (a *Ali) Verify(publicKey, sign []byte, req *NotifyReq) error {
	p, _ := pem.Decode(publicKey)
	if p == nil {
		panic("Public key broken!")
	}
	pub, err := x509.ParsePKIXPublicKey(p.Bytes)
	if err != nil {
		return err
	}
	h := crypto.Hash.New(crypto.SHA1)
	m := bronx.Params(structs.Map(req))
	b := sortedParams(removeKeys(m, "sign", "sign_type"))
	h.Write(removeQuote(b.Bytes()))
	sum := h.Sum(nil)
	if sign, err = base64.StdEncoding.DecodeString(string(sign)); err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA1, sum, sign)
}

// EncodedQuery ...
func (a *Ali) EncodedQuery(s interface{}) []byte {
	m := bronx.Params(structs.Map(s))
	m["sign"] = url.QueryEscape(m["sign"])
	buf := sortedParams(m)
	return buf.Bytes()
}

// PayURL ...
func (a *Ali) PayURL(s interface{}) string {
	u, err := url.Parse(orderURL)
	if err != nil {
		panic(err)
	}
	m := bronx.Params(structs.Map(s))
	p := url.Values{}
	for k := range m {
		p.Add(k, m[k])
	}
	u.RawQuery = p.Encode()
	return u.String()
}

// NotifyReq ...
type NotifyReq struct {
	NotifyTime       string `structs:"notify_time" form:"notify_time" json:"notify_time"`
	NotifyType       string `structs:"notify_type" form:"notify_type" json:"notify_type"`
	NotifyID         string `structs:"notify_id" form:"notify_id" json:"notify_id"`
	SignType         string `structs:"sign_type" form:"sign_type" json:"sign_type"`
	Sign             string `structs:"sign" form:"sign" json:"sign"`
	OutTradeNo       string `structs:"out_trade_no" form:"out_trade_no" json:"out_trade_no"`
	Subject          string `structs:"subject" form:"subject" json:"subject"`
	PaymentType      string `structs:"payment_type" form:"payment_type" json:"payment_type"`
	TradeNo          string `structs:"trade_no" form:"trade_no" json:"trade_no"`
	TradeStatus      string `structs:"trade_status" form:"trade_status" json:"trade_status"`
	GmtCreate        string `structs:"gmt_create" form:"gmt_create" json:"gmt_create"`
	GmtPayment       string `structs:"gmt_payment" form:"gmt_payment" json:"gmt_payment"`
	GmtClose         string `structs:"gmt_close" form:"gmt_close" json:"gmt_close"`
	RefundStatus     string `structs:"refund_status" form:"refund_status" json:"refund_status"`
	GmtRefund        string `structs:"gmt_refund" form:"gmt_refund" json:"gmt_refund"`
	SellerEmail      string `structs:"seller_email" form:"seller_email" json:"seller_email"`
	BuyerEmail       string `structs:"buyer_email" form:"buyer_email" json:"buyer_email"`
	SellerID         string `structs:"seller_id" form:"seller_id" json:"seller_id"`
	BuyerID          string `structs:"buyer_id" form:"buyer_id" json:"buyer_id"`
	Price            string `structs:"price" form:"price" json:"price"`
	TotalFee         string `structs:"total_fee" form:"total_fee" json:"total_fee"`
	Quantity         string `structs:"quantity" form:"quantity" json:"quantity"`
	Body             string `structs:"body" form:"body" json:"body"`
	Discount         string `structs:"discount" form:"discount" json:"discount"`
	IsTotalFeeAdjust string `structs:"is_total_fee_adjust" form:"is_total_fee_adjust" json:"is_total_fee_adjust"`
	UseCoupon        string `structs:"use_coupon" form:"use_coupon" json:"use_coupon"`
	ExtraCommonParam string `structs:"extra_common_param" form:"extra_common_param" json:"extra_common_param"`
	BusinessScene    string `structs:"business_scene" form:"business_scene" json:"business_scene"`
}
