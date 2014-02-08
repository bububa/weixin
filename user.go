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
	var js []byte
	if err == nil {
		js, _ = json.Marshal(user)
	} else {
		switch err.(type) {
		case response:
			js, _ = json.Marshal(err)
		default:
			res := response{
				ErrorCode:    0,
				ErrorMessage: err.Error(),
			}
			js, _ = json.Marshal(res)
		}

	}
	w.writer.Write(js)
	return
}

// Get Subscribers
func (w responseWriter) GetSubscribers(nextOpenId string) (subscribers *Subscribers, err error) {
	subscribers, err = w.wx.GetSubscribers(nextOpenId)
	var js []byte
	if err == nil {
		js, _ = json.Marshal(subscribers)
	} else {
		switch err.(type) {
		case response:
			js, _ = json.Marshal(err)
		default:
			res := response{
				ErrorCode:    0,
				ErrorMessage: err.Error(),
			}
			js, _ = json.Marshal(res)
		}

	}
	w.writer.Write(js)
	return
}

// Get SubscribersWithInfo
func (w responseWriter) GetSubscribersWithInfo(nextOpenId string) (subscribers *Subscribers, users []*User, err error) {
	subscribers, err = w.wx.GetSubscribers(nextOpenId)
	var js []byte
	if err != nil {
		switch err.(type) {
		case response:
			js, _ = json.Marshal(err)
		default:
			res := response{
				ErrorCode:    0,
				ErrorMessage: err.Error(),
			}
			js, _ = json.Marshal(res)
		}
		return
	}
	for _, openId := range subscribers.Data.OpenId {
		user, er := w.wx.GetUser(openId, "zh_CN")
		if er != nil {
			switch er.(type) {
			case response:
				js, _ = json.Marshal(er)
			default:
				res := response{
					ErrorCode:    0,
					ErrorMessage: er.Error(),
				}
				js, _ = json.Marshal(res)
			}
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
	js, err = json.Marshal(msg)
	if err != nil {
		res := response{
			ErrorCode:    0,
			ErrorMessage: err.Error(),
		}
		js, _ = json.Marshal(res)
		return
	}
	w.writer.Write(js)
	return
}
