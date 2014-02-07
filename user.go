package weixin

import (
	"encoding/json"
)

// Get User Info
func (wx *Weixin) GetUser(openId string, lang string) (user *User, err error) {
	if lang == "" {
		lang = "zh_CN"
	}
	gateway := weixinUserURL + "/info?openid=" + openId + "&lang=" + lang + "&access_token="
	reply, err := apiGET(gateway, wx.tokenChan)
	if err == nil && reply != nil {
		err = json.Unmarshal(reply, user)
	}
	return
}

// Get Subscribers
func (wx *Weixin) GetSubscribers(nextOpenId string) (subscribers *Subscribers, err error) {
	gateway := weixinUserURL + "/get?next_openid=" + nextOpenId + "&access_token="
	reply, err := apiGET(gateway, wx.tokenChan)
	if err == nil && reply != nil {
		err = json.Unmarshal(reply, subscribers)
	}
	return
}

// Get User Info
func (w responseWriter) GetUser(openId string, lang string) (user *User, err error) {
	return w.wx.GetUser(openId, lang)
}

// Get Subscribers
func (w responseWriter) GetSubscribers(nextOpenId string) (*Subscribers, error) {
	return w.wx.GetSubscribers(nextOpenId)
}
