package weixin

const (
	QR_SCENE       = "QR_SCENE"
	QR_LIMIT_SCENE = "QR_LIMIT_SCENE"
	// Event type
	msgEvent         = "event"
	EventSubscribe   = "subscribe"
	EventUnsubscribe = "unsubscribe"
	EventScan        = "SCAN"
	EventClick       = "CLICK"
	// Message type
	MsgTypeDefault          = ".*"
	MsgTypeText             = "text"
	MsgTypeImage            = "image"
	MsgTypeVoice            = "voice"
	MsgTypeVideo            = "video"
	MsgTypeLocation         = "location"
	MsgTypeLink             = "link"
	MsgTypeEvent            = msgEvent + ".*"
	MsgTypeEventSubscribe   = msgEvent + "\\." + EventSubscribe
	MsgTypeEventUnsubscribe = msgEvent + "\\." + EventUnsubscribe
	MsgTypeEventScan        = msgEvent + "\\." + EventScan
	MsgTypeEventClick       = msgEvent + "\\." + EventClick

	// Media type
	MediaTypeImage = "image"
	MediaTypeVoice = "voice"
	MediaTypeVideo = "video"
	MediaTypeThumb = "thumb"
	// Weixin host URL
	weixinHost          = "https://api.weixin.qq.com/cgi-bin"
	weixinFileURL       = "http://file.api.weixin.qq.com/cgi-bin/media"
	weixinUserURL       = "https://api.weixin.qq.com/cgi-bin/user"
	weixinGroupURL      = "https://api.weixin.qq.com/cgi-bin/groups"
	weixinMenuURL       = "https://api.weixin.qq.com/cgi-bin/menu"
	weixinQrcodeURL     = "https://api.weixin.qq.com/cgi-bin/qrcode"
	weixinShowQrcodeURL = "https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket="
	// Reply format
	replyText    = "<xml>%s<MsgType><![CDATA[text]]></MsgType><Content><![CDATA[%s]]></Content></xml>"
	replyImage   = "<xml>%s<MsgType><![CDATA[image]]></MsgType><Image><MediaId><![CDATA[%s]]></MediaId></Image></xml>"
	replyVoice   = "<xml>%s<MsgType><![CDATA[voice]]></MsgType><Voice><MediaId><![CDATA[%s]]></MediaId></Voice></xml>"
	replyVideo   = "<xml>%s<MsgType><![CDATA[video]]></MsgType><Video><MediaId><![CDATA[%s]]></MediaId><Title><![CDATA[%s]]></Title><Description><![CDATA[%s]]></Description></Video></xml>"
	replyMusic   = "<xml>%s<MsgType><![CDATA[music]]></MsgType><Music><Title><![CDATA[%s]]></Title><Description><![CDATA[%s]]></Description><MusicUrl><![CDATA[%s]]></MusicUrl><HQMusicUrl><![CDATA[%s]]></HQMusicUrl><ThumbMediaId><![CDATA[%s]]></ThumbMediaId></Music></xml>"
	replyNews    = "<xml>%s<MsgType><![CDATA[news]]></MsgType><ArticleCount>%d</ArticleCount><Articles>%s</Articles></xml>"
	replyHeader  = "<ToUserName><![CDATA[%s]]></ToUserName><FromUserName><![CDATA[%s]]></FromUserName><CreateTime>%d</CreateTime>"
	replyArticle = "<item><Title><![CDATA[%s]]></Title> <Description><![CDATA[%s]]></Description><PicUrl><![CDATA[%s]]></PicUrl><Url><![CDATA[%s]]></Url></item>"
)
