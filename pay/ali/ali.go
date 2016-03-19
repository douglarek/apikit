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
	"net/http"
	"net/url"
	"sort"

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
func (a *Ali) Sign(s interface{}, secretKey []byte) (b []byte) {
	m := bronx.Params(structs.Map(s))
	st := m["sign_type"]
	delete(m, "sign")
	delete(m, "sign_type")
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
		buf.WriteString(fmt.Sprintf("%s=%s", k, m[k]))
	}

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
		h.Write(buf.Bytes())
		return h.Sum(nil)
	}
	return
}

// Verify for RSA sign.
func Verify(publicKey, message, sign []byte) error {
	p, _ := pem.Decode(publicKey)
	if p == nil {
		panic("Public key broken!")
	}
	pub, err := x509.ParsePKIXPublicKey(p.Bytes)
	if err != nil {
		panic(err)
	}
	h := crypto.Hash.New(crypto.SHA1)
	h.Write(message)
	sum := h.Sum(nil)
	return rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA1, sum, sign)
}

// EncodedQuery ...
func (a *Ali) EncodedQuery(s interface{}) []byte {
	var buf bytes.Buffer
	for k, v := range bronx.Params(structs.Map(s)) {
		buf.WriteString(fmt.Sprintf("%s=%q&", k, url.QueryEscape(v)))
	}
	buf.Truncate(buf.Len() - 1)
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
