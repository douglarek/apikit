package apikit

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/structs"
)

//
const (
	MediaJSON = "application/json;charset=utf-8"
	MediaForm = "application/x-www-form-urlencoded;charset=utf-8"
	MediaXML  = "application/xml;charset=utf-8"
)

// A Client manages communication with API.
type Client struct {
	client *http.Client
	header map[string]string
}

// H is a map shortcut.
type H map[string]string

// NewClient returns a new API client.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}
	c := &Client{client: httpClient, header: H{}}
	return c
}

// SetHeader sets a header.
func (c *Client) SetHeader(h H) {
	for k, v := range h {
		c.header[k] = v
	}
}

// AddHeader adds a header mapping.
func (c *Client) AddHeader(key, val string) {
	c.header[key] = val
}

// NewRequest creates an API request.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.Reader

	switch ct := c.header["Content-Type"]; ct {
	case MediaJSON:
		if body != nil {
			b, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			buf = bytes.NewBuffer(b)
		}
	case MediaForm, "":
		v := url.Values{}
		for k, val := range Params(structs.Map(body)) {
			v.Set(k, val)
		}
		buf = strings.NewReader(v.Encode())
		c.AddHeader("Content-Type", MediaForm)
	case MediaXML:
		if body != nil {
			var b []byte
			b, err := xml.Marshal(body)
			if err != nil {
				return nil, err
			}
			buf = bytes.NewBuffer(b)
		}
	default:
		panic("unsupported content type!")
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if len(c.header) != 0 {
		for k, v := range c.header {
			req.Header.Add(k, v)
		}
	}
	return req, nil
}

// Do sends an API request and returns the API response.
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			if r, _ := regexp.MatchString("(plain|xml|xhtml)", resp.Header.Get("Content-Type")); r {
				if err := xml.NewDecoder(resp.Body).Decode(v); err == io.EOF {
					err = nil
				}
			} else {
				if err := json.NewDecoder(resp.Body).Decode(v); err == io.EOF {
					err = nil
				}
			}
		}
	}
	return resp, err
}

// Params expands a nested map.
func Params(m0 map[string]interface{}) (m map[string]string) {
	if m == nil {
		m = make(map[string]string)
	}

	for k, v := range m0 {
		val := reflect.ValueOf(v)
		if v == nil {
			continue
		}
		switch val.Kind() {
		case reflect.Map:
			for k, v0 := range Params(v.(map[string]interface{})) {
				m[k] = v0
			}
		case reflect.String:
			if val.Len() != 0 {
				m[k] = v.(string)
			}
		case reflect.Int:
			if v != 0 {
				m[k] = strconv.FormatInt(int64(v.(int)), 10)
			}
		case reflect.Array, reflect.Slice:
			b, _ := json.Marshal(val.Interface())
			m[k] = string(b)
		default:
			panic(fmt.Sprintf("unsupported type: %T", v))
		}
	}
	return
}
