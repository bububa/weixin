package weixin

import (
	"encoding/json"
)

// Get User Info
func (wx *Weixin) GetUser(openId string, lang string) (user *User, err error) {
	if lang == "" {
		lang = "zh_CN"
	}
	gateway := "https://api.weixin.qq.com/cgi-bin/user/info?openid=" + openId + "&lang=" + lang + "&access_token="
	reply, err := apiGET(gateway, wx.tokenChan)
	if err == nil && reply != nil {
		err = json.Unmarshal(reply, user)
	}
	return
}

// Get Subscribers
func (wx *Weixin) GetSubscribers(nextOpenId string) (subscribers *Subscribers, err error) {
	gateway := "https://api.weixin.qq.com/cgi-bin/user/get?next_openid=" + nextOpenId + "&access_token="
	reply, err := apiGET(gateway, wx.tokenChan)
	if err == nil && reply != nil {
		err = json.Unmarshal(reply, subscribers)
	}
	return
}
