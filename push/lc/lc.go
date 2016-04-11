package lc

import (
	"net/http"
	"time"

	"github.com/douglarek/bronx"
)

const (
	installationURL = "https://leancloud.cn/1.1/installations"
	pushURL         = "https://leancloud.cn/1.1/push"
)

// LeanCloud ...
type LeanCloud struct {
	client *bronx.Client
}

// New makes a LeanCloud ...
func New(lcID, lcKey string) *LeanCloud {
	c := bronx.NewClient(http.DefaultClient, bronx.MediaJSON)
	c.Header = map[string]string{"X-LC-Id": lcID, "X-LC-Key": lcKey}
	return &LeanCloud{client: c}
}

// InstallationReq ...
type InstallationReq struct {
	DeviceType     string   `json:"deviceType"`
	DeviceToken    string   `json:"deviceToken"`
	InstallationID string   `json:"installationId"`
	Channels       []string `json:"channels"`
}

// Resp ...
type Resp struct {
	ObjectID  string    `json:"objectId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ChannelsReq ...
type ChannelsReq struct {
	Channels string `json:"channels"`
	Prod     string `json:"prod"`
	Data     struct {
		Alert string `json:"alert"`
	} `json:"data"`
}

// SaveInstallation ...
func (lc *LeanCloud) SaveInstallation(r *InstallationReq) (*Resp, error) {
	req, err := lc.client.NewRequest("POST", installationURL, r)
	if err != nil {
		return nil, err
	}
	res := new(Resp)
	if _, err := lc.client.Do(req, res); err != nil {
		return nil, err
	}
	return res, nil
}

// PushChannels ...
func (lc *LeanCloud) PushChannels(r *ChannelsReq) (*Resp, error) {
	req, err := lc.client.NewRequest("POST", pushURL, r)
	if err != nil {
		return nil, err
	}
	res := new(Resp)
	if _, err := lc.client.Do(req, res); err != nil {
		return nil, err
	}
	return res, nil

}
