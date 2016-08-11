package lc

import (
	"time"

	"github.com/douglarek/apikit"
)

const (
	installationURL = "https://leancloud.cn/1.1/installations"
	pushURL         = "https://leancloud.cn/1.1/push"
)

// LeanCloud ...
type LeanCloud struct {
	client *apikit.Client
}

// New makes a LeanCloud ...
func New(lcID, lcKey string) *LeanCloud {
	c := apikit.NewClient(nil)
	c.SetHeader(apikit.H{"X-LC-Id": lcID, "X-LC-Key": lcKey, "Content-Type": apikit.MediaJSON})
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
	Channels []string               `json:"channels"`
	Prod     string                 `json:"prod"`
	Data     map[string]interface{} `json:"data"`
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

// ChannelsOpsReq ...
type ChannelsOpsReq struct {
	Channels struct {
		OP      string   `json:"__op"`
		Objects []string `json:"objects"`
	} `json:"channels"`
}

// UnsubscribeChannel ...
func (lc *LeanCloud) UnsubscribeChannel(objectID string, c *ChannelsOpsReq) (*Resp, error) {
	req, err := lc.client.NewRequest("PUT", installationURL+"/"+objectID, c)
	if err != nil {
		return nil, err
	}
	res := new(Resp)
	if _, err := lc.client.Do(req, res); err != nil {
		return nil, err
	}
	return res, nil
}
