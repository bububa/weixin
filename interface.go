package weixin

import (
	"fmt"
	log "github.com/bububa/factorlog"
	"github.com/bububa/pg"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"
	"sync"
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
	FormValues   url.Values
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
	Subscribe       uint8  `json:"subscribe"`
	OpenId          string `json:"openid"`
	Nick            string `json:"nickname"`
	Sex             uint8  `json:"sex"`
	City            string `json:"city"`
	Country         string `json:"country"`
	Province        string `json:"province"`
	Language        string `json:"language"`
	HeadImgUrl      string `json:"headimgurl"`
	SubscribeTime   uint64 `json:"subscribe_time"`
    UnsubscribeTime uint64 `json:"unsubscribe_time"`
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

type Button struct {
	Type      string  `json:"type,omitempty"`
	Name      string  `json:"name,omitempty"`
	Key       string  `json:"key,omitempty"`
	Url       string  `json:"url,omitempty"`
	SubButton *Button `json:"sub_button,omitempty"`
}

type Menu struct {
	Button []*Button `json:"button,omitempty"`
}

type TicketReply struct {
	Ticket        string `json:"ticket"`
	ExpireSeconds uint   `json:"expire_seconds"`
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
	GetSubscribersWithInfo(nextOpenId string) (*Subscribers, []*User, error)
	// Menu operator
	CreateMenu(menu *Menu) error
    DeleteMenu() error
	// Orcode operator
	CreateQrcode(sceneId uint64) (*TicketReply, error)
	CreateTempQrcode(sceneId uint64, expireSeconds uint) (*TicketReply, error)
	ShowQrcode(sceneId uint64) error
	ShowTempQrcode(sceneId uint64, expireSeconds uint) error
	// Helper
	PgDB() *pg.DB
	App() string
    Wx() *Weixin
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

func (e response) Error() string {
	return fmt.Sprintf("CODE:%v, MSG:%v", e.ErrorCode, e.ErrorMessage)
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
	app         string
	token       string
	routes      []*route
	tokenChan   chan accessToken
	pg          *pg.DB
}

type SceneParams struct {
    Expires     uint64      `json:"expires,omitempty"`
    Callback    string      `json:"callback"`
    Created     time.Time   `json:"-"`
}

type ScenesMap struct {
	Mutex *sync.Mutex
	Scenes map[uint64]SceneParams
}
