package bronx

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
)

//
const (
	MediaJSON = "application/json;charset=utf-8"
	MediaForm = "application/x-www-form-urlencoded;charset=utf-8"

	HMAC = "hmac"
	MD5  = "md5"
)

// A Client manages communication with API.
type Client struct {
	client      *http.Client
	ContentType string
}

// NewClient returns a new API client.
func NewClient(httpClient *http.Client, contentType ...string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	c := &Client{client: httpClient}
	if len(contentType) > 0 {
		c.ContentType = contentType[0]
	}
	return c
}

// NewRequest creates an API request.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.Reader

	switch c.ContentType {
	case MediaJSON:
		if body != nil {
			var b []byte
			err := json.Unmarshal(b, body)
			if err != nil {
				return nil, err
			}
			buf = bytes.NewBuffer(b)
		}
	case MediaForm, "":
		f, _ := query.Values(body)
		buf = strings.NewReader(f.Encode())
		c.ContentType = MediaForm
	default:
		panic("unsupported content type!")
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", c.ContentType)
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
			err = json.NewDecoder(resp.Body).Decode(v)
			if err == io.EOF {
				err = nil // ignore EOF errors caused by empty response body
			}
		}
	}
	return resp, err
}
