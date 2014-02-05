// Weixin MP SDK (Golang)
package weixin

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/bububa/factorlog"
	"io/ioutil"
	"net/http"
	"sort"
	"time"
)

// Create a Weixin instance
func New(token string, appid string, secret string) *Weixin {
	wx := &Weixin{}
	wx.token = token
	if len(appid) > 0 && len(secret) > 0 {
		wx.tokenChan = make(chan accessToken)
		go createAccessToken(wx.tokenChan, appid, secret)
	}
	return wx
}

func checkSignature(t string, w http.ResponseWriter, r *http.Request) bool {
	r.ParseForm()
	var signature string = r.FormValue("signature")
	var timestamp string = r.FormValue("timestamp")
	var nonce string = r.FormValue("nonce")
	strs := sort.StringSlice{t, timestamp, nonce}
	sort.Strings(strs)
	var str string
	for _, s := range strs {
		str += s
	}
	h := sha1.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil)) == signature
}

func authAccessToken(appid string, secret string) (string, time.Duration) {
	resp, err := http.Get(weixinHost + "/token?grant_type=client_credential&appid=" + appid + "&secret=" + secret)
	if err != nil {
		log.Errorf("Get access token failed: ", err)
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("Read access token failed: ", err)
		} else {
			var res struct {
				AccessToken string `json:"access_token"`
				ExpiresIn   int64  `json:"expires_in"`
			}
			if err := json.Unmarshal(body, &res); err != nil {
				log.Errorf("Parse access token failed: ", err)
			} else {
				//log.Debugf("AuthAccessToken token=%s expires_in=%d", res.AccessToken, res.ExpiresIn)
				return res.AccessToken, time.Duration(res.ExpiresIn * 1000 * 1000 * 1000)
			}
		}
	}
	return "", 0
}

func createAccessToken(c chan accessToken, appid string, secret string) {
	token := accessToken{"", time.Now()}
	c <- token
	for {
		if time.Since(token.expires).Seconds() >= 0 {
			var expires time.Duration
			token.token, expires = authAccessToken(appid, secret)
			token.expires = time.Now().Add(expires)
		}
		c <- token
	}
}

func apiPOST(gateway string, c chan accessToken, msg interface{}) ([]byte, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	for i := 0; i < 3; i++ {
		token := <-c
		if time.Since(token.expires).Seconds() < 0 {
			r, err := http.Post(gateway+token.token, "application/json; charset=utf-8", bytes.NewReader(data))
			if err != nil {
				return nil, err
			}
			defer r.Body.Close()
			reply, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return nil, err
			}
			var result response
			if err := json.Unmarshal(reply, &result); err != nil {
				return nil, err
			} else {
				switch result.ErrorCode {
				case 0:
					return reply, nil
				case 42001: // access_token timeout and retry
					continue
				default:
					return reply, errors.New(fmt.Sprintf("WeiXin reply[%d]: %s", result.ErrorCode, result.ErrorMessage))
				}
			}
		}
	}
	return nil, errors.New("WeiXin post message too many times")
}

func apiGET(gateway string, c chan accessToken) ([]byte, error) {
	for i := 0; i < 3; i++ {
		token := <-c
		if time.Since(token.expires).Seconds() < 0 {
			r, err := http.Get(gateway + token.token)
			if err != nil {
				return nil, err
			}
			defer r.Body.Close()
			reply, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return nil, err
			}
			var result response
			if err := json.Unmarshal(reply, &result); err != nil {
				return nil, err
			} else {
				switch result.ErrorCode {
				case 0:
					return reply, nil
				case 42001: // access_token timeout and retry
					continue
				default:
					return nil, errors.New(fmt.Sprintf("WeiXin download[%d]: %s", result.ErrorCode, result.ErrorMessage))
				}
			}
		}
	}
	return nil, errors.New("WeiXin get too many times")
}
