package ali

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"github.com/douglarek/bronx"
	"github.com/fatih/structs"
)

const orderURL = "https://mapi.alipay.com/gateway.do"

// Ali ...
type Ali struct {
	client *bronx.Client
}

// New makes a wechat ...
func New(httpClient *http.Client) *Ali {
	c := bronx.NewClient(httpClient)
	return &Ali{client: c}
}

// Req ...
type Req struct {
	Service      string `structs:"service"`
	Partner      string `structs:"partner"`
	InputCharset string `structs:"_input_charset"`
	SignType     string `structs:"sign_type"`
	Sign         string `structs:"sign"`
	NotifyURL    string `structs:"notify_url"`
	ReturnURL    string `structs:"return_url"`
}

// OrderReq ...
type OrderReq struct {
	Req
	OutTradeNo        string `structs:"out_trade_no"`
	Subject           string `structs:"subject"`
	PaymentType       string `structs:"payment_type"`
	TotalFee          string `structs:"total_fee"`
	SellerID          string `structs:"seller_id"`
	SellerEmail       string `structs:"seller_email"`
	SellerAccountName string `structs:"seller_account_name"`
}

// Sign ...
func (a *Ali) Sign(s interface{}, secret string) string {
	m := bronx.Params(structs.Map(s))
	delete(m, "sign")
	delete(m, "sign_type")
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
	buf.Truncate(buf.Len() - 1)
	buf.WriteString(secret)
	return fmt.Sprintf("%x", md5.Sum(buf.Bytes()))
}

// PayStr ...
func (a *Ali) PayStr(s interface{}) string {
	m := bronx.Params(structs.Map(s))
	p := url.Values{}
	for k := range m {
		p.Add(k, m[k])
	}
	return p.Encode()
}

// PayURL ...
func (a *Ali) PayURL(s interface{}) string {
	u, err := url.Parse(orderURL)
	if err != nil {
		panic(err)
	}
	u.RawQuery = a.PayStr(s)
	return u.String()
}
