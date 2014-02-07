package main

import (
	"github.com/bububa/weixin"
)

// 文本消息的处理函数
func Echo(w weixin.ResponseWriter, r *weixin.Request) {
	txt := r.Content          // 获取用户发送的消息
	w.ReplyText(txt)          // 回复一条文本消息
	w.PostText("Post:" + txt) // 发送一条文本消息
}

// 关注事件的处理函数
func Subscribe(w weixin.ResponseWriter, r *weixin.Request) {
	w.ReplyText("欢迎关注") // 有新人关注，返回欢迎消息
}
