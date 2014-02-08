package weixin

import (
	"encoding/json"
)

// Create Group
func (wx *Weixin) CreateGroup(name string) (group *Group, err error) {
	gateway := weixinGroupURL + "/create?access_token="
	var msg struct {
		Group struct {
			Name string `json:"name"`
		} `json:"group"`
	}
	msg.Group.Name = name
	reply, err := apiPOST(gateway, wx.tokenChan, &msg)
	if err == nil && reply != nil {
		err = json.Unmarshal(reply, group)
	}
	return
}

// Get Groups
func (wx *Weixin) GetGroups() (groups []Group, err error) {
	gateway := weixinGroupURL + "/get?access_token="
	reply, err := apiGET(gateway, wx.tokenChan)
	if err == nil && reply != nil {
		err = json.Unmarshal(reply, &groups)
	}
	return
}

// Get User Group
func (wx *Weixin) GetUserGroup(openId string) (group *Group, err error) {
	gateway := weixinGroupURL + "/getid?access_token="
	var msg struct {
		OpenId string `json:"openid"`
	}
	msg.OpenId = openId
	reply, err := apiPOST(gateway, wx.tokenChan, &msg)
	if err == nil && reply != nil {
		err = json.Unmarshal(reply, group)
	}
	return
}

// Change Group Name
func (wx *Weixin) ChangeGroupName(group *Group) error {
	gateway := weixinGroupURL + "/update?access_token="
	var msg struct {
		Group struct {
			Id   uint64 `json:"id"`
			Name string `json:"name"`
		} `json:"group"`
	}
	msg.Group.Id = group.Id
	msg.Group.Name = group.Name
	_, err := apiPOST(gateway, wx.tokenChan, &msg)
	return err
}

// Change User Group
func (wx *Weixin) ChangeUserGroup(openId string, groupId uint64) error {
	gateway := weixinGroupURL + "/members/update?access_token="
	var msg struct {
		OpenId  string `json:"openidid"`
		GroupId uint64 `json:"to_groupid"`
	}
	msg.OpenId = openId
	msg.GroupId = groupId
	_, err := apiPOST(gateway, wx.tokenChan, &msg)
	return err
}

// Create Group
func (w responseWriter) CreateGroup(name string) (group *Group, err error) {
	group, err = w.wx.CreateGroup(name)
	var js []byte
	if err == nil {
		js, _ = json.Marshal(group)
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

// Get Groups
func (w responseWriter) GetGroups() (groups []Group, err error) {
	groups, err = w.wx.GetGroups()
	var js []byte
	if err == nil {
		js, _ = json.Marshal(groups)
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

// Get User Group
func (w responseWriter) GetUserGroup(openId string) (group *Group, err error) {
	group, err = w.wx.GetUserGroup(openId)
	var js []byte
	if err == nil {
		js, _ = json.Marshal(group)
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

// Change Group Name
func (w responseWriter) ChangeGroupName(group *Group) error {
	err := w.wx.ChangeGroupName(group)
	var js []byte
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
	w.writer.Write(js)
	return err
}

// Get User Group
func (w responseWriter) ChangeUserGroup(openId string, groupId uint64) error {
	err := w.wx.ChangeUserGroup(openId, groupId)
	var js []byte
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
	w.writer.Write(js)
	return err
}
