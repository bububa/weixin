package weixin

// Post text message
func (wx *Weixin) PostText(touser string, text string) error {
	var msg struct {
		ToUser  string `json:"touser"`
		MsgType string `json:"msgtype"`
		Text    struct {
			Content string `json:"content"`
		} `json:"text"`
	}
	msg.ToUser = touser
	msg.MsgType = "text"
	msg.Text.Content = text
	return postMessage(wx.tokenChan, &msg)
}

// Post image message
func (wx *Weixin) PostImage(touser string, mediaId string) error {
	var msg struct {
		ToUser  string `json:"touser"`
		MsgType string `json:"msgtype"`
		Image   struct {
			MediaId string `json:"media_id"`
		} `json:"image"`
	}
	msg.ToUser = touser
	msg.MsgType = "image"
	msg.Image.MediaId = mediaId
	return postMessage(wx.tokenChan, &msg)
}

// Post voice message
func (wx *Weixin) PostVoice(touser string, mediaId string) error {
	var msg struct {
		ToUser  string `json:"touser"`
		MsgType string `json:"msgtype"`
		Voice   struct {
			MediaId string `json:"media_id"`
		} `json:"voice"`
	}
	msg.ToUser = touser
	msg.MsgType = "voice"
	msg.Voice.MediaId = mediaId
	return postMessage(wx.tokenChan, &msg)
}

// Post video message
func (wx *Weixin) PostVideo(touser string, m string, t string, d string) error {
	var msg struct {
		ToUser  string `json:"touser"`
		MsgType string `json:"msgtype"`
		Video   struct {
			MediaId     string `json:"media_id"`
			Title       string `json:"title"`
			Description string `json:"description"`
		} `json:"video"`
	}
	msg.ToUser = touser
	msg.MsgType = "video"
	msg.Video.MediaId = m
	msg.Video.Title = t
	msg.Video.Description = d
	return postMessage(wx.tokenChan, &msg)
}

// Post music message
func (wx *Weixin) PostMusic(touser string, music *Music) error {
	var msg struct {
		ToUser  string `json:"touser"`
		MsgType string `json:"msgtype"`
		Music   *Music `json:"music"`
	}
	msg.ToUser = touser
	msg.MsgType = "video"
	msg.Music = music
	return postMessage(wx.tokenChan, &msg)
}

// Post news message
func (wx *Weixin) PostNews(touser string, articles []Article) error {
	var msg struct {
		ToUser  string `json:"touser"`
		MsgType string `json:"msgtype"`
		News    struct {
			Articles []Article `json:"articles"`
		} `json:"news"`
	}
	msg.ToUser = touser
	msg.MsgType = "news"
	msg.News.Articles = articles
	return postMessage(wx.tokenChan, &msg)
}

func postMessage(c chan accessToken, msg interface{}) error {
	reqURL := weixinHost + "/message/custom/send?access_token="
	_, err := apiPOST(reqURL, c, msg)
	return err
}

// Post text message
func (w responseWriter) PostText(text string) error {
	return w.wx.PostText(w.toUserName, text)
}

// Post image message
func (w responseWriter) PostImage(mediaId string) error {
	return w.wx.PostImage(w.toUserName, mediaId)
}

// Post voice message
func (w responseWriter) PostVoice(mediaId string) error {
	return w.wx.PostVoice(w.toUserName, mediaId)
}

// Post video message
func (w responseWriter) PostVideo(mediaId string, title string, desc string) error {
	return w.wx.PostVideo(w.toUserName, mediaId, title, desc)
}

// Post music message
func (w responseWriter) PostMusic(music *Music) error {
	return w.wx.PostMusic(w.toUserName, music)
}

// Post news message
func (w responseWriter) PostNews(articles []Article) error {
	return w.wx.PostNews(w.toUserName, articles)
}
