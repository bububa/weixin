package weixin

import (
	"fmt"
	"time"
)

// Format reply message header
func (w responseWriter) replyHeader() string {
	return fmt.Sprintf(replyHeader, w.toUserName, w.fromUserName, time.Now().Unix())
}

// Reply text message
func (w responseWriter) ReplyText(text string) {
	msg := fmt.Sprintf(replyText, w.replyHeader(), text)
	w.writer.Write([]byte(msg))
}

// Reply image message
func (w responseWriter) ReplyImage(mediaId string) {
	msg := fmt.Sprintf(replyImage, w.replyHeader(), mediaId)
	w.writer.Write([]byte(msg))
}

// Reply voice message
func (w responseWriter) ReplyVoice(mediaId string) {
	msg := fmt.Sprintf(replyVoice, w.replyHeader(), mediaId)
	w.writer.Write([]byte(msg))
}

// Reply video message
func (w responseWriter) ReplyVideo(mediaId string, title string, description string) {
	msg := fmt.Sprintf(replyVideo, w.replyHeader(), mediaId, title, description)
	w.writer.Write([]byte(msg))
}

// Reply music message
func (w responseWriter) ReplyMusic(m *Music) {
	msg := fmt.Sprintf(replyMusic, w.replyHeader(), m.Title, m.Description, m.MusicUrl, m.HQMusicUrl, m.ThumbMediaId)
	w.writer.Write([]byte(msg))
}

// Reply news message (max 10 news)
func (w responseWriter) ReplyNews(articles []Article) {
	var ctx string
	for _, article := range articles {
		ctx += fmt.Sprintf(replyArticle, article.Title, article.Description, article.PicUrl, article.Url)
	}
	msg := fmt.Sprintf(replyNews, w.replyHeader(), len(articles), ctx)
	w.writer.Write([]byte(msg))
}
