package main

import (
	"github.com/bububa/weixin"
)

// 文本消息的处理函数
func MsgTxt(w weixin.ResponseWriter, r *weixin.Request) {
	logger.Debugf("USER:%s TXT:%s", r.FromUserName, r.Content)
	w.ReplyText(r.Content)          // 回复一条文本消息
	w.PostText("Post:" + r.Content) // 发送一条文本消息
	_, err := w.PgDB().ExecOne(`SELECT * FROM plproxy.new_txt_msg(?, ?, ?, ?, ?)`, w.App(), r.FromUserName, r.ToUserName, r.CreateTime, r.MsgType, r.MsgId, r.Content)
	if err != nil {
		logger.Error(err)
	}
}

// 图片消息的处理函数
func MsgImage(w weixin.ResponseWriter, r *weixin.Request) {
	logger.Debugf("USER:%s IMG:%d", r.FromUserName, r.MediaId)
	_, err := w.PgDB().ExecOne(`SELECT * FROM plproxy.new_img_msg(?, ?, ?, ?, ?, ?)`, w.App(), r.FromUserName, r.ToUserName, r.CreateTime, r.MsgType, r.MsgId, r.PicUrl, r.MediaId)
	if err != nil {
		logger.Error(err)
	}
}

// 语音消息的处理函数
func MsgVoice(w weixin.ResponseWriter, r *weixin.Request) {
	logger.Debugf("USER:%s VOICE:%d", r.FromUserName, r.MediaId)
	_, err := w.PgDB().ExecOne(`SELECT * FROM plproxy.new_voice_msg(?, ?, ?, ?, ?, ?)`, w.App(), r.FromUserName, r.ToUserName, r.CreateTime, r.MsgType, r.MsgId, r.MediaId, r.Format)
	if err != nil {
		logger.Error(err)
	}
}

// 视频消息的处理函数
func MsgVideo(w weixin.ResponseWriter, r *weixin.Request) {
	logger.Debugf("USER:%s VIDEO:%d", r.FromUserName, r.MediaId)
	_, err := w.PgDB().ExecOne(`SELECT * FROM plproxy.new_video_msg(?, ?, ?, ?, ?, ?)`, w.App(), r.FromUserName, r.ToUserName, r.CreateTime, r.MsgType, r.MsgId, r.MediaId, r.ThumbMediaId)
	if err != nil {
		logger.Error(err)
	}
}

// 位置消息的处理函数
func MsgLocation(w weixin.ResponseWriter, r *weixin.Request) {
	logger.Debugf("USER:%s LOCATION:%.2f, %.2f", r.FromUserName, r.LocationX, r.LocationY)
	_, err := w.PgDB().ExecOne(`SELECT * FROM plproxy.new_location_msg(?, ?, ?, ?, ?, ?, ?, ?)`, w.App(), r.FromUserName, r.ToUserName, r.CreateTime, r.MsgType, r.MsgId, r.LocationX, r.LocationY, r.Scale, r.Label)
	if err != nil {
		logger.Error(err)
	}
}

// 链接消息的处理函数
func MsgLink(w weixin.ResponseWriter, r *weixin.Request) {
	logger.Debugf("USER:%s LINK:%s", r.FromUserName, r.Url)
	_, err := w.PgDB().ExecOne(`SELECT * FROM plproxy.new_video_msg(?, ?, ?, ?, ?, ?, ?)`, w.App(), r.FromUserName, r.ToUserName, r.CreateTime, r.MsgType, r.MsgId, r.Title, r.Description, r.Url)
	if err != nil {
		logger.Error(err)
	}
}

// 关注事件的处理函数
func Subscribe(w weixin.ResponseWriter, r *weixin.Request) {
	logger.Debugf("USER:%s SUBSCRIBE", r.FromUserName)
	w.ReplyText("欢迎关注Lens杂志") // 有新人关注，返回欢迎消息
	_, err := w.PgDB().ExecOne(`SELECT * FROM plproxy.subscribe(?, ?, ?)`, w.App(), r.FromUserName, r.CreateTime)
	if err != nil {
		logger.Error(err)
	}
}

// 取消关注事件的处理函数
func Unsubscribe(w weixin.ResponseWriter, r *weixin.Request) {
	logger.Debugf("USER:%s UNSUBSCRIBE", r.FromUserName)
	_, err := w.PgDB().ExecOne(`SELECT * FROM plproxy.unsubscribe(?, ?, ?)`, w.App(), r.FromUserName, r.CreateTime)
	if err != nil {
		logger.Error(err)
	}
}
