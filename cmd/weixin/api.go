package main

import (
	"github.com/bububa/weixin"
	"github.com/bububa/pg"
	"encoding/json"
    "net/http"
    "net/url"
    "io/ioutil"
	"strconv"
    "time"
)

// 取得订阅用户列表的处理函数
func GetSubscribers(w weixin.ResponseWriter, r *weixin.Request) {
	_, users, err := w.GetSubscribersWithInfo(r.FormValues.Get("nextopenid"))
	if err != nil {
		logger.Warn(err)
		return
	}
	for _, user := range users {
        create := weixin.ParseCreateTime(user.SubscribeTime).Format("2006-01-02 15:04:05")
		_, err := w.PgDB().ExecOne(`SELECT * FROM plproxy.update_user_info(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, w.App(), user.OpenId, user.Subscribe, user.Nick, user.Sex, user.City, user.Province, user.Country, user.Language, user.HeadImgUrl, create)
		if err != nil {
			logger.Error(err)
		}
	}
}

// 取得订阅用户信息的处理函数
func GetUser(w weixin.ResponseWriter, r *weixin.Request) {
	user, err := w.GetUser(r.FormValues.Get("openid"), r.FormValues.Get("lang"))
	if err != nil {
		logger.Warn(err)
		return
	}
    create := weixin.ParseCreateTime(user.SubscribeTime).Format("2006-01-02 15:04:05")
    _, err = w.PgDB().ExecOne(`SELECT * FROM plproxy.update_user_info(?::text, ?::text, ?::int2, ?::text, ?::int2, ?::text, ?::text, ?::text, ?::text, ?::text, ?::timestamp)`, w.App(), user.OpenId, user.Subscribe, user.Nick, user.Sex, user.City, user.Province, user.Country, user.Language, user.HeadImgUrl, create)
	if err != nil && err != pg.ErrNoRows {
		logger.Error(err)
	}
}

// 创建用户组的处理函数
func CreateGroup(w weixin.ResponseWriter, r *weixin.Request) {
    if ! checkAuth(r.FormValues.Get("openid")) {
        logger.Warn("Access denied")
        return
    }
	group, err := w.CreateGroup(r.FormValues.Get("name"))
	if err != nil {
		logger.Warn(err)
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
		logger.Warn(err)
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
    if ! checkAuth(r.FormValues.Get("openid")) {
        logger.Warn("Access denied")
        return
    }
	groupId, err := strconv.ParseUint(r.FormValues.Get("id"), 10, 64)
	if err != nil {
		logger.Error(err)
		return
	}
	group := &weixin.Group{
		Id:   groupId,
		Name: r.FormValues.Get("name"),
	}
	err = w.ChangeGroupName(group)
	if err != nil {
		logger.Warn(err)
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
		logger.Warn(err)
		return
	}
	_, err = w.PgDB().ExecOne(`SELECT * FROM plproxy.change_user_group(?, ?, ?)`, w.App(), r.FormValues.Get("openid"), group.Id)
	if err != nil {
		logger.Error(err)
	}
}

// 修改用户所在用户组的处理函数
func ChangeUserGroup(w weixin.ResponseWriter, r *weixin.Request) {
    if ! checkAuth(r.FormValues.Get("openid")) {
        logger.Warn("Access denied")
        return
    }
	groupId, err := strconv.ParseUint(r.FormValues.Get("groupid"), 10, 64)
	if err != nil {
		logger.Error(err)
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
    if ! checkAuth(r.FormValues.Get("openid")) {
        logger.Warn("Access denied")
        return
    }
	var menu *weixin.Menu
	err := json.Unmarshal([]byte(r.FormValues.Get("menu")), &menu)
	if err != nil {
		logger.Error(err)
		return
	}
	err = w.CreateMenu(menu)
	if err != nil {
		logger.Warn(err)
	}
}

// 删除菜单的处理函数
func DeleteMenu(w weixin.ResponseWriter, r *weixin.Request) {
    if ! checkAuth(r.FormValues.Get("openid")) {
        logger.Warn("Access denied")
        return
    }
    err := w.DeleteMenu()
    if err != nil {
        logger.Warn(err)
    }
}

// 生成QRCODE的处理函数
func CreateQrcode(w weixin.ResponseWriter, r *weixin.Request) {
    sceneId, _, err := createScene(w, r)
	if err != nil {
		logger.Error(err)
		return
	}
	_, err = w.CreateQrcode(sceneId)
	if err != nil {
		logger.Warn(err)
	}
}

// 显示QRCODE的处理函数
func ShowQrcode(w weixin.ResponseWriter, r *weixin.Request) {
    sceneId, _, err := createScene(w, r)
    if err != nil {
        logger.Error(err)
        return
    }
    err = w.ShowQrcode(sceneId)
    if err != nil {
        logger.Warn(err)
    }
}

// 生成临时QRCODE的处理函数
func CreateTempQrcode(w weixin.ResponseWriter, r *weixin.Request) {
    sceneId, expireSeconds, err := createScene(w, r)
	if err != nil {
		logger.Error(err)
		return
	}
	_, err = w.CreateTempQrcode(sceneId, uint(expireSeconds))
	if err != nil {
		logger.Warn(err)
	}
}

// 显示临时QRCODE的处理函数
func ShowTempQrcode(w weixin.ResponseWriter, r *weixin.Request) {
    sceneId, expireSeconds, err := createScene(w, r)
    if err != nil {
        logger.Error(err)
        return
    }
    err = w.ShowTempQrcode(sceneId, uint(expireSeconds))
    if err != nil {
        logger.Warn(err)
    }
}

// 发送客服消息
func SendCustomMessage(w weixin.ResponseWriter, r *weixin.Request) {
    var err error
    touser := r.FormValues.Get("touser")
    switch r.FormValues.Get("msgtype") {
    case "text":
        err = w.Wx().PostText(touser, r.FormValues.Get("text"))
    case "image":
        err = w.Wx().PostImage(touser, r.FormValues.Get("media_id"))
    case "voice":
        err = w.Wx().PostVoice(touser, r.FormValues.Get("media_id"))
    case "video":
        err = w.Wx().PostVideo(touser, r.FormValues.Get("media_id"), r.FormValues.Get("title"), r.FormValues.Get("description"))
    case "music":
        var music weixin.Music
        err = json.Unmarshal([]byte(r.FormValues.Get("music")), &music)
        if err == nil {
            err = w.Wx().PostMusic(touser, &music)
        }
    case "news":
        var articles []weixin.Article
        err = json.Unmarshal([]byte(r.FormValues.Get("articles")), &articles)
        if err == nil {
            err = w.Wx().PostNews(touser, articles)
        }
    }
    if err != nil {
        logger.Error(err)
    }
}

// 生成二维码参数
func createScene(w weixin.ResponseWriter, r *weixin.Request) (sceneId uint64, expiresSeconds uint64, err error) {
    sceneId, err = strconv.ParseUint(r.FormValues.Get("sceneid"), 10, 64)
    if err == nil && sceneId > 0 {
        return sceneId, expiresSeconds, nil
    }
    sceneParams := r.FormValues.Get("scene")
    var scene weixin.SceneParams
    err = json.Unmarshal([]byte(sceneParams), &scene)
    if err != nil {
        return
    }
    expiresSeconds = scene.Expires
    scene.Created = time.Now()
    if len(scenesMap.Scenes) > 0 {
        for k, s := range scenesMap.Scenes {
            if s.Created.Add(time.Duration(s.Expires) * time.Second).Before(time.Now()) {
                //logger.Infof("Now:%v, Expires:%v, Created:%v", time.Now(), s.Expires, s.Created)
                scenesMap.Mutex.Lock()
                delete(scenesMap.Scenes, k)
                scenesMap.Mutex.Unlock()
                continue
            }
            if k > sceneId {
                sceneId = k
            }
        }
    }
    sceneId += 1
    scenesMap.Mutex.Lock()
    scenesMap.Scenes[sceneId] = scene
    scenesMap.Mutex.Unlock()
    //logger.Trace(scenesMap.Scenes)
    return
}

// 验证身份
func checkAuth(openId string) bool {
    res, err := http.PostForm(callbackUrl["Auth"], url.Values{"openid": {openId}})
    if err != nil {
        logger.Error(err)
        return false
    }
    defer res.Body.Close()
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        logger.Error(err)
        return false
    }
    if string(body) == "1" {
        return true
    }
    return false
}
