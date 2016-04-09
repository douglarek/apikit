package xg

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"sort"

	"github.com/douglarek/bronx"
	"github.com/fatih/structs"
)

const xgPushURL = "http://openapi.xg.qq.com/v2/push"

// Req ...
type Req struct {
	AccessID  string `structs:"access_id"`
	TimeStamp string `structs:"timestamp"`
	ValidTime string `structs:"valid_time"`
	Sign      string `structs:"sign"`
}

// SingleDeviceReq ...
type SingleDeviceReq struct {
	Req
	DeviceToken string `structs:"device_token"`
	MessageType string `structs:"message_type"`
	Message     string `structs:"message"`
	ExpireTime  string `structs:"expire_time"`
	SendTime    string `structs:"send_time"`
	MultiPkg    string `structs:"multi_pkg"`
	Environment string `structs:"environment"`
}

func params(prefix, secret string, m map[string]string) []byte {
	var keys []string
	for k := range m {
		if k == "sign" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var buf bytes.Buffer
	buf.WriteString(prefix)
	for _, k := range keys {
		buf.WriteString(k + "=" + m[k])
	}
	buf.WriteString(secret)
	return buf.Bytes()
}

func sign(s interface{}, prefix, secret string) string {
	m := bronx.Params(structs.Map(s))
	b := params(prefix, secret, m)
	sum := md5.Sum(b)
	return hex.EncodeToString(sum[:])
}

// Resp ...
type Resp struct {
	RetCode int    `json:"ret_code"`
	ErrMsg  string `json:"err_msg"`
	Result  struct {
		Status string `json:"status"`
		PushID string `json:"push_id"`
	} `json:"result"`
}

// Success ...
func (resp *Resp) Success() bool {
	return resp.RetCode == 0
}

func newURL(oldURL, suffix string) *url.URL {
	u, _ := url.Parse(oldURL)
	u.Path = path.Join(u.Path, suffix)
	return u
}

func values(s interface{}) url.Values {
	v := url.Values{}
	for k, val := range bronx.Params(structs.Map(s)) {
		v.Set(k, val)
	}
	return v
}

// SinglePush ...
func SinglePush(r *SingleDeviceReq, secret string) (*Resp, error) {
	u := newURL(xgPushURL, "single_device")
	prefix := "POST" + u.String()[len(u.Scheme)+3:]
	r.Sign = sign(r, prefix, secret)
	resp, err := http.PostForm(u.String(), values(r))
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	var resp0 Resp
	if err := json.NewDecoder(resp.Body).Decode(&resp0); err != nil && err != io.EOF {
		return nil, err
	}
	return &resp0, nil
}

// MultipleDeviceReq ...
type MultipleDeviceReq struct {
	SingleDeviceReq
	DeviceList []string `structs:"device_list"`
	PushID     string   `structs:"push_id"`
}

func unmarshalBody(body io.ReadCloser) (r *Resp, err error) {
	defer body.Close()
	if err != nil {
		return
	}
	err = json.NewDecoder(body).Decode(&r)
	if err != nil && err != io.EOF {
		return
	}
	return r, nil
}

func createMultiPush(r *MultipleDeviceReq, secret string) (s string, err error) {
	u := newURL(xgPushURL, "create_multipush")
	prefix := "POST" + u.String()[len(u.Scheme)+3:]
	r.Sign = sign(r, prefix, secret)
	resp, err := http.PostForm(u.String(), values(r))
	if err != nil {
		return
	}
	resp0, err := unmarshalBody(resp.Body)
	if err != nil {
		return
	}
	return resp0.Result.PushID, nil
}

func listMultiple(r *MultipleDeviceReq, secret string) (*Resp, error) {
	u := newURL(xgPushURL, "device_list_multiple")
	prefix := "POST" + u.String()[len(u.Scheme)+3:]
	r.Sign = sign(r, prefix, secret)
	resp, err := http.PostForm(u.String(), values(r))
	if err != nil {
		return nil, err
	}
	return unmarshalBody(resp.Body)
}

// MultiPush ...
func MultiPush(r *MultipleDeviceReq, secret string) (*Resp, error) {
	s, _ := createMultiPush(r, secret)
	r.PushID = s
	return listMultiple(r, secret)
}
