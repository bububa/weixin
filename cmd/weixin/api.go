package main

import (
	"encoding/json"
	"github.com/bububa/weixin"
	"strconv"
)

// 取得订阅用户列表的处理函数
func GetSubscribers(w weixin.ResponseWriter, r *weixin.Request) {
	_, users, err := w.GetSubscribersWithInfo(r.FormValues.Get("nextopenid"))
	if err != nil {
		return
	}
	for _, user := range users {
		_, err := w.PgDB().ExecOne(`SELECT * FROM plproxy.update_user_info(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, w.App(), user.OpenId, user.Subscribe, user.Nick, user.Sex, user.City, user.Province, user.Country, user.Language, user.HeadImgUrl, user.SubscribeTime)
		if err != nil {
			logger.Error(err)
		}
	}
}

// 取得订阅用户信息的处理函数
func GetUser(w weixin.ResponseWriter, r *weixin.Request) {
	user, err := w.GetUser(r.FormValues.Get("openid"), r.FormValues.Get("lang"))
	if err != nil {
		return
	}
	_, err = w.PgDB().ExecOne(`SELECT * FROM plproxy.update_user_info(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, w.App(), user.OpenId, user.Subscribe, user.Nick, user.Sex, user.City, user.Province, user.Country, user.Language, user.HeadImgUrl, user.SubscribeTime)
	if err != nil {
		logger.Error(err)
	}
}

// 创建用户组的处理函数
func CreateGroup(w weixin.ResponseWriter, r *weixin.Request) {
	group, err := w.CreateGroup(r.FormValues.Get("name"))
	if err != nil {
		return
	}
	_, err = w.PgDB().ExecOne(`SELECT * FROM plproxy.new_group(?, ?, ?, ?)`, w.App(), group.Id, group.Name, group.Num)
	if err != nil {
		logger.Error(err)
	}
}

// 获取用户组列表的处理函数
func GetGroups(w weixin.ResponseWriter, r *weixin.Request) {
	groups, err := w.GetGroups()
	if err != nil {
		return
	}
	for _, group := range groups {
		_, err = w.PgDB().ExecOne(`SELECT * FROM plproxy.new_group(?, ?, ?, ?)`, w.App(), group.Id, group.Name, group.Num)
		if err != nil {
			logger.Error(err)
		}
	}
}

// 修改用户组名的处理函数
func ChangeGroupName(w weixin.ResponseWriter, r *weixin.Request) {
	groupId, err := strconv.ParseUint(r.FormValues.Get("id"), 10, 64)
	if err != nil {
		return
	}
	group := &weixin.Group{
		Id:   groupId,
		Name: r.FormValues.Get("name"),
	}
	err = w.ChangeGroupName(group)
	if err != nil {
		return
	}
	_, err = w.PgDB().ExecOne(`SELECT * FROM plproxy.new_group(?, ?, ?, ?)`, w.App(), group.Id, group.Name, group.Num)
	if err != nil {
		logger.Error(err)
	}
}

// 获得用户所在用户组的处理函数
func GetUserGroup(w weixin.ResponseWriter, r *weixin.Request) {
	group, err := w.GetUserGroup(r.FormValues.Get("openid"))
	if err != nil {
		return
	}
	_, err = w.PgDB().ExecOne(`SELECT * FROM plproxy.change_user_group(?, ?, ?)`, w.App(), r.FormValues.Get("openid"), group.Id)
	if err != nil {
		logger.Error(err)
	}
}

// 修改用户所在用户组的处理函数
func ChangeUserGroup(w weixin.ResponseWriter, r *weixin.Request) {
	groupId, err := strconv.ParseUint(r.FormValues.Get("groupid"), 10, 64)
	if err != nil {
		return
	}
	err = w.ChangeUserGroup(r.FormValues.Get("openid"), groupId)
	if err != nil {
		return
	}
	_, err = w.PgDB().ExecOne(`SELECT * FROM plproxy.change_user_group(?, ?, ?)`, w.App(), r.FormValues.Get("openid"), groupId)
	if err != nil {
		logger.Error(err)
	}
}

// 创建菜单的处理函数
func CreateMenu(w weixin.ResponseWriter, r *weixin.Request) {
	var menu *weixin.Menu
	err := json.Unmarshal([]byte(r.FormValues.Get("menu")), menu)
	if err != nil {
		return
	}
	err = w.CreateMenu(menu)
	if err != nil {
		return
	}
}
