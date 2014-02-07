package weixin

import (
	log "github.com/bububa/factorlog"
	"github.com/bububa/pg"
	"io"
	"net/http"
	"regexp"
	"time"
)

var logger *log.FactorLog

func SetLogger(aLogger *log.FactorLog) {
	logger = aLogger
}

// Common message header
type MessageHeader struct {
	ToUserName   string
	FromUserName string
	CreateTime   int
	MsgType      string
}

// Weixin request
type Request struct {
	MessageHeader
	MsgId        int64
	Content      string
	PicUrl       string
	MediaId      string
	Format       string
	ThumbMediaId string
	LocationX    float32 `xml:"Location_X"`
	LocationY    float32 `xml:"Location_Y"`
	Scale        float32
	Label        string
	Title        string
	Description  string
	Url          string
	Event        string
	EventKey     string
	Ticket       string
	Latitude     float32
	Longitude    float32
	Precision    float32
	Recognition  string
}

// Use to reply music message
type Music struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	MusicUrl     string `json:"musicurl"`
	HQMusicUrl   string `json:"hqmusicurl"`
	ThumbMediaId string `json:"thumb_media_id"`
}

// Use to reply news message
type Article struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	PicUrl      string `json:"picurl"`
	Url         string `json:"url"`
}

type User struct {
	Subscribe       uint8  `xml:"subscribe"`
	OpenId          string `xml:"openid"`
	Nick            string `xml:"nickname"`
	Sex             uint8  `xml:"sex"`
	City            string `xml:"city"`
	Country         string `xml:"country"`
	Province        string `xml:"province"`
	language        string `xml:"language"`
	HeadImgUrl      string `xml:"headimgurl"`
	SubscribeTime   uint64 `xml:"subscribe_time"`
	UnsubscribeTime uint64
	GroupId         uint64 `json:"groupid"`
}

type Group struct {
	Id   uint64 `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Num  uint64 `json:"count,omitempty"`
}

type Subscribers struct {
	Total uint64 `json:"total"`
	Count uint32 `json:"count"`
	Data  struct {
		OpenId []string `json:"openid"`
	} `json:"data"`
	NextOpenId string `json:"next_openid"`
}

// Use to output reply
type ResponseWriter interface {
	// Reply message
	ReplyText(text string)
	ReplyImage(mediaId string)
	ReplyVoice(mediaId string)
	ReplyVideo(mediaId string, title string, description string)
	ReplyMusic(music *Music)
	ReplyNews(articles []Article)
	// Post message
	PostText(text string) error
	PostImage(mediaId string) error
	PostVoice(mediaId string) error
	PostVideo(mediaId string, title string, description string) error
	PostMusic(music *Music) error
	PostNews(articles []Article) error
	// Media operator
	UploadMediaFromFile(mediaType string, filepath string) (string, error)
	DownloadMediaToFile(mediaId string, filepath string) error
	UploadMedia(mediaType string, filename string, reader io.Reader) (string, error)
	DownloadMedia(mediaId string, writer io.Writer) error
	// Group operator
	CreateGroup(name string) (*Group, error)
	GetGroups() ([]Group, error)
	GetUserGroup(openId string) (*Group, error)
	ChangeGroupName(group *Group) error
	ChangeUserGroup(openId string, groupId uint64) error
	// User operator
	GetUser(openId string, lang string) (*User, error)
	GetSubscribers(nextOpenId string) (*Subscribers, error)
	// Helper
	PgDB() *pg.DB
	App() string
}

type responseWriter struct {
	wx           *Weixin
	writer       http.ResponseWriter
	toUserName   string
	fromUserName string
}

type response struct {
	ErrorCode    int    `json:"errcode"`
	ErrorMessage string `json:"errmsg"`
}

// Callback function
type HandlerFunc func(ResponseWriter, *Request)

type route struct {
	regex   *regexp.Regexp
	handler HandlerFunc
}

type accessToken struct {
	token   string
	expires time.Time
}

type Weixin struct {
	app       string
	token     string
	routes    []*route
	tokenChan chan accessToken
	pg        *pg.DB
}
