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
	user, err = w.wx.GetUser(openId, lang)
	if err != nil {
		logger.Warn(err)
		return
	}
	js, err := json.Marshal(user)
	if err != nil {
		logger.Warn(err)
		return
	}
	w.writer.Write(js)
	return
}

// Get Subscribers
func (w responseWriter) GetSubscribers(nextOpenId string) (subscribers *Subscribers, err error) {
	subscribers, err = w.wx.GetSubscribers(nextOpenId)
	if err != nil {
		logger.Warn(err)
		return
	}
	js, err := json.Marshal(subscribers)
	if err != nil {
		logger.Warn(err)
		return
	}
	w.writer.Write(js)
	return
}

// Get SubscribersWithInfo
func (w responseWriter) GetSubscribersWithInfo(nextOpenId string) (subscribers *Subscribers, users []*User, err error) {
	subscribers, err = w.wx.GetSubscribers(nextOpenId)
	if err != nil {
		logger.Warn(err)
		return
	}
	for _, openId := range subscribers.Data.OpenId {
		user, er := w.wx.GetUser(openId, "zh_CN")
		if er != nil {
			logger.Warn(er)
			err = er
			return
		}
		users = append(users, user)
	}
	var msg struct {
		Subscribers *Subscribers `json:"subscribers"`
		Users       []*User      `json:"users"`
	}
	msg.Subscribers = subscribers
	msg.Users = users
	js, err := json.Marshal(msg)
	if err != nil {
		logger.Warn(err)
		return
	}
	w.writer.Write(js)
	return
}
