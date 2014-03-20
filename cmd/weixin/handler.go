package main

import (
	"github.com/bububa/weixin"
	"github.com/bububa/pg"
    "io/ioutil"
    "net/http"
    "net/url"
    "strconv"
    "strings"
)

// 已关注扫描二维码事件
func MsgScan(w weixin.ResponseWriter, r *weixin.Request) {
    create := weixin.ParseCreateTime(r.CreateTime)
    logger.Debugf("USER:%s CREATE:%s", r.FromUserName, create)
    _, err := http.PostForm(callbackUrl["Scan"], url.Values{"openid": {r.FromUserName}, "evenkey": {r.EventKey}, "ticket": {r.Ticket}, "create": {create.Format("2006-01-02 15:04:05")}})
    if err != nil {
        logger.Error(err)
    }
}

// 文本消息的处理函数
func MsgTxt(w weixin.ResponseWriter, r *weixin.Request) {
    create := weixin.ParseCreateTime(r.CreateTime)
    logger.Debugf("USER:%s TXT:%s CREATE:%s", r.FromUserName, r.Content, create)
    _, err := http.PostForm(callbackUrl["MsgTxt"], url.Values{"openid": {r.FromUserName}, "msgid": {strconv.FormatInt(r.MsgId, 10)}, "content": {r.Content}, "create": {create.Format("2006-01-02 15:04:05")}}) 
    if err != nil {
        logger.Error(err)
    }
    _, err = w.PgDB().ExecOne(`SELECT * FROM plproxy.new_txt_msg(?::text, ?::text, ?::text, ?::timestamp, ?::text, ?::int8, ?::text)`, w.App(), r.FromUserName, r.ToUserName, create, r.MsgType, r.MsgId, r.Content)
	if err != nil && err != pg.ErrNoRows {
		logger.Error(err)
	}
}

// 图片消息的处理函数
func MsgImage(w weixin.ResponseWriter, r *weixin.Request) {
    create := weixin.ParseCreateTime(r.CreateTime)
    logger.Debugf("USER:%s IMG:%s CREATE:%s", r.FromUserName, r.MediaId, create)
    _, err := http.PostForm(callbackUrl["MsgImage"], url.Values{"openid": {r.FromUserName}, "msgid": {strconv.FormatInt(r.MsgId, 10)}, "mediaid": {r.MediaId}, "picurl": {r.PicUrl}, "create": {create.Format("2006-01-02 15:04:05")}})
    if err != nil {
        logger.Error(err)
    }
    //w.ReplyImage(r.MediaId) // 返回发送的图片
    _, err = w.PgDB().ExecOne(`SELECT * FROM plproxy.new_img_msg(?::text, ?::text, ?::text, ?::timestamp, ?::text, ?::int8, ?::text, ?::text)`, w.App(), r.FromUserName, r.ToUserName, create, r.MsgType, r.MsgId, r.PicUrl, r.MediaId)
	if err != nil && err != pg.ErrNoRows {
		logger.Error(err)
	}
}

// 语音消息的处理函数
func MsgVoice(w weixin.ResponseWriter, r *weixin.Request) {
    create := weixin.ParseCreateTime(r.CreateTime)
	logger.Debugf("USER:%s VOICE:%s CREATE:%s", r.FromUserName, r.MediaId, create)
    _, err := http.PostForm(callbackUrl["MsgVoice"], url.Values{"openid": {r.FromUserName}, "msgid": {strconv.FormatInt(r.MsgId, 10)}, "mediaid": {r.MediaId}, "format": {r.Format}, "create": {create.Format("2006-01-02 15:04:05")}})
    if err != nil {
        logger.Error(err)
    }
    //w.ReplyVoice(r.MediaId) // 返回发送的声音
    _, err = w.PgDB().ExecOne(`SELECT * FROM plproxy.new_voice_msg(?::text, ?::text, ?::text, ?::timestamp, ?::text, ?::int8, ?::text, ?::text)`, w.App(), r.FromUserName, r.ToUserName, create, r.MsgType, r.MsgId, r.MediaId, r.Format)
	if err != nil && err != pg.ErrNoRows {
		logger.Error(err)
	}
}

// 视频消息的处理函数
func MsgVideo(w weixin.ResponseWriter, r *weixin.Request) {
    create := weixin.ParseCreateTime(r.CreateTime)
	logger.Debugf("USER:%s VIDEO:%s CREATE:%s", r.FromUserName, r.MediaId, create)
    _, err := http.PostForm(callbackUrl["MsgVideo"], url.Values{"openid": {r.FromUserName}, "msgid": {strconv.FormatInt(r.MsgId, 10)}, "mediaid": {r.MediaId}, "thumbmediaid": {r.ThumbMediaId}, "create": {create.Format("2006-01-02 15:04:05")}})
    if err != nil {
        logger.Error(err)
    }
    _, err = w.PgDB().ExecOne(`SELECT * FROM plproxy.new_video_msg(?::text, ?::text, ?::text, ?::timestamp, ?::text, ?::int8, ?::text, ?::text)`, w.App(), r.FromUserName, r.ToUserName, create, r.MsgType, r.MsgId, r.MediaId, r.ThumbMediaId)
	if err != nil && err != pg.ErrNoRows {
		logger.Error(err)
	}
}

// 位置消息的处理函数
func MsgLocation(w weixin.ResponseWriter, r *weixin.Request) {
    create := weixin.ParseCreateTime(r.CreateTime)
    logger.Debugf("USER:%s LOCATION:%.2f, %.2f CREATE:%s", r.FromUserName, r.LocationX, r.LocationY, create)
    _, err := http.PostForm(callbackUrl["MsgLocation"], url.Values{"openid": {r.FromUserName}, "msgid": {strconv.FormatInt(r.MsgId, 10)}, "locationx": {strconv.FormatFloat(float64(r.LocationX), 'f', 1, 32)}, "locationy": {strconv.FormatFloat(float64(r.LocationY), 'f', 1, 32)}, "scale": {strconv.FormatFloat(float64(r.Scale), 'f', 1, 32)}, "label": {r.Label}, "create": {create.Format("2006-01-02 15:04:05")}})
    if err != nil {
        logger.Error(err)
    }
    _, err = w.PgDB().ExecOne(`SELECT * FROM plproxy.new_location_msg(?::text, ?::text, ?::text, ?::timestamp, ?::text, ?::int8, ?::float4, ?::float4, ?::int4, ?::text)`, w.App(), r.FromUserName, r.ToUserName, create, r.MsgType, r.MsgId, r.LocationX, r.LocationY, r.Scale, r.Label)
	if err != nil && err != pg.ErrNoRows {
		logger.Error(err)
	}
}

// 链接消息的处理函数
func MsgLink(w weixin.ResponseWriter, r *weixin.Request) {
    create := weixin.ParseCreateTime(r.CreateTime)
    logger.Debugf("USER:%s LINK:%s CREATE:%s", r.FromUserName, r.Url, create)
    _, err := http.PostForm(callbackUrl["MsgLink"], url.Values{"openid": {r.FromUserName}, "msgid": {strconv.FormatInt(r.MsgId, 10)}, "title": {r.Title}, "description": {r.Description}, "url": {r.Url}, "create": {create.Format("2006-01-02 15:04:05")}})
    if err != nil {
        logger.Error(err)
    }
	_, err = w.PgDB().ExecOne(`SELECT * FROM plproxy.new_link_msg(?, ?, ?, ?, ?, ?, ?, ?, ?)`, w.App(), r.FromUserName, r.ToUserName, create, r.MsgType, r.MsgId, r.Title, r.Description, r.Url)
	if err != nil && err != pg.ErrNoRows {
		logger.Error(err)
	}
}

// 关注事件的处理函数
func Subscribe(w weixin.ResponseWriter, r *weixin.Request) {
    create := weixin.ParseCreateTime(r.CreateTime)
	eventKey, err := strconv.ParseUint(strings.Replace(r.EventKey, "qrscene_", "", -1), 10, 64)
    if err != nil {
        logger.Error(err)
        return
    }
    logger.Debugf("USER:%s SUBSCRIBE CREATE:%s EVENT KEY:%d", r.FromUserName, create, eventKey)
    msg, err := subscribeCallback(eventKey, r.FromUserName)
    if err != nil {
        logger.Error(err)
        return
    }
    if msg != "" {
	    w.PostText(msg) // 有新人关注，返回欢迎消息
    }
    logger.Debugf("USER:%s SUBSCRIBE MSG: %s", r.FromUserName, msg)
    _, err = w.PgDB().ExecOne(`SELECT * FROM plproxy.subscribe(?::text, ?::text, ?::timestamp)`, w.App(), r.FromUserName, create)
	if err != nil && err != pg.ErrNoRows {
		logger.Error(err)
	}
    go func() {
        _, err := http.PostForm(apiHost + "/" + w.App() + "/user", url.Values{"openid": {r.FromUserName}})
        if err != nil {
            logger.Error(err)
        }
    }()
}

// 取消关注事件的处理函数
func Unsubscribe(w weixin.ResponseWriter, r *weixin.Request) {
    create := weixin.ParseCreateTime(r.CreateTime)
    logger.Debugf("USER:%s UNSUBSCRIBE CREATE:%s", r.FromUserName, create)
    _, err := http.PostForm(callbackUrl["Unsubscribe"], url.Values{"openid": {r.FromUserName}})
    if err != nil {
        logger.Error(err)
    }
    _, err = w.PgDB().ExecOne(`SELECT * FROM plproxy.unsubscribe(?::text, ?::text, ?::timestamp)`, w.App(), r.FromUserName, create)
	if err != nil && err != pg.ErrNoRows {
		logger.Error(err)
	}
}

func subscribeCallback(eventKey uint64, openId string) (msg string, err error) {
    if scene, found := scenesMap.Scenes[eventKey]; found {
        callback, err := url.QueryUnescape(scene.Callback)
        if err != nil {
            return msg, err
        }

        res, err := http.PostForm(callback, url.Values{"openid": {openId}})
        if err != nil {
            return msg, err
        }
        defer res.Body.Close()
        body, err := ioutil.ReadAll(res.Body)
        if err != nil {
            return msg, err
        }
        msg = string(body)
    }
    return msg, err
}
