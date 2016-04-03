package union

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/douglarek/bronx"
	"github.com/fatih/structs"
)

const appTransReqURL = "https://gateway.95516.com/gateway/api/appTransReq.do"

// OrderReq ...
type OrderReq struct {
	Version      string `structs:"version" json:"version"`
	Encoding     string `structs:"encoding" json:"encoding"`
	CertID       string `structs:"certId" json:"certId"`
	SignMethod   string `structs:"signMethod" json:"signMethod"`
	Signature    string `structs:"signature" json:"signature"`
	TxnType      string `structs:"txnType" json:"txnType"`
	TxnSubType   string `structs:"txnSubType" json:"txnSubType"`
	BizType      string `structs:"bizType" json:"bizType"`
	AccessType   string `structs:"accessType" json:"accessType"`
	MerID        string `structs:"merId" json:"merId"`
	OrderID      string `structs:"orderId" json:"orderId"`
	CurrencyCode string `structs:"currencyCode" json:"currencyCode"`
	TxnAmt       string `structs:"txnAmt" json:"txnAmt"`
	TxnTime      string `structs:"txnTime" json:"txnTime"`
	FrontURL     string `structs:"frontUrl" json:"frontUrl"`
	BackURL      string `structs:"backUrl" json:"backUrl"`
	ChannelType  string `structs:"channelType" json:"channelType"`
	AccType      string `structs:"accType" json:"accType"`
	OrderDesc    string `structs:"orderDesc" json:"orderDesc"`
	ReqReserved  string `structs:"reqReserved" json:"reqReserved"`
}

// DefaultOrderReq ...
func DefaultOrderReq() *OrderReq {
	return &OrderReq{
		Version:      "5.0.0",
		Encoding:     "UTF-8",
		TxnType:      "01",
		TxnSubType:   "01",
		BizType:      "000201",
		SignMethod:   "01",
		ChannelType:  "08",
		AccessType:   "0",
		TxnTime:      time.Now().UTC().Add(8 * time.Hour).Format("20060102150405"),
		AccType:      "01",
		TxnAmt:       "1",
		CurrencyCode: "156",
		ReqReserved:  "{}",
	}

}

func params(s interface{}) []byte {
	m := bronx.Params(structs.Map(s))
	delete(m, "signature")
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
	return buf.Bytes()
}

// Sign ...
func Sign(s interface{}, secretKey []byte) string {
	p, _ := pem.Decode(secretKey)
	if p == nil {
		panic("failed to parse PEM")
	}

	k, err := x509.ParsePKCS1PrivateKey(p.Bytes)
	if err != nil {
		panic(err)
	}

	hashed := sha1.Sum([]byte(fmt.Sprintf("%x", sha1.Sum(params(s)))))
	sign, err := rsa.SignPKCS1v15(rand.Reader, k, crypto.SHA1, hashed[:])
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(sign)
}

// OrderResp ...
type OrderResp struct {
	OrderReq
	RespCode string `structs:"respCode" json:"respCode"`
	RespMsg  string `structs:"respMsg" json:"respMsg"`
	Tn       string `structs:"tn" json:"tn"`
}

// AppConsume ...
func AppConsume(r *OrderReq) (oresp *OrderResp) {
	u := url.Values{}
	for k, v := range bronx.Params(structs.Map(r)) {
		u.Set(k, v)
	}
	resp, err := http.PostForm(appTransReqURL, u)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	m := map[string]string{}
	for _, v := range strings.Split(string(body), "&") {
		val := strings.SplitN(v, "=", 2)
		m[val[0]] = val[1]
	}
	if m["respCode"] != "00" {
		return
	}
	b, _ := json.Marshal(m)
	json.Unmarshal(b, &oresp)
	return oresp
}

// Verify ...
func Verify(s interface{}, publicKey, sign []byte) error {
	p, _ := pem.Decode(publicKey)
	if p == nil {
		panic("failed to parse pem")
	}
	c, err := x509.ParseCertificate(p.Bytes)
	if err != nil {
		panic(err)
	}
	sign, _ = base64.StdEncoding.DecodeString(string(sign))
	hashed := sha1.Sum([]byte(fmt.Sprintf("%x", sha1.Sum(params(s)))))
	return rsa.VerifyPKCS1v15(c.PublicKey.(*rsa.PublicKey), crypto.SHA1, hashed[:], sign)
}

// NotifyReq ...
type NotifyReq struct {
	OrderResp
	PayType            string `json:"payType"`
	AccNo              string `json:"accNo"`
	PayCardType        string `json:"payCardType"`
	Reserved           string `json:"reserved"`
	QueryID            string `json:"queryId"`
	TraceNo            string `json:"traceNo"`
	TraceTime          string `json:"traceTime"`
	SettleDate         string `json:"settleDate"`
	SettleCurrencyCode string `json:"settleCurrencyCode"`
	SettleAmt          string `json:"settleAmt"`
	PayCardNo          string `json:"payCardNo"`
	PayCardIssueName   string `json:"payCardIssueName"`
}
