package weixin

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

// Register request callback.
func (wx *Weixin) HandleFunc(pattern string, handler HandlerFunc) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		logger.Fatal(err)
		return
	}
	route := &route{regex, handler}
	wx.routes = append(wx.routes, route)
}

// Process weixin request and send response.
func (wx *Weixin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Host[0:4] == "api." {
		logger.Info(r.URL.Path)
		for _, route := range wx.routes {
			if !route.regex.MatchString(r.URL.Path) {
				continue
			}
			writer := responseWriter{}
			writer.wx = wx
			writer.writer = w
			req := Request{
				FormValues: r.Form,
			}
			route.handler(writer, &req)
			return
		}
		return
	}

	if !checkSignature(wx.token, w, r) {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	// Verify request
	if r.Method == "GET" {
		fmt.Fprintf(w, r.FormValue("echostr"))
		return
	}
	// Process message
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("Weixin receive message failed:", err)
		http.Error(w, "", http.StatusBadRequest)
	} else {
		var msg Request
		if err := xml.Unmarshal(data, &msg); err != nil {
			logger.Errorf("Weixin parse message failed:", err)
			http.Error(w, "", http.StatusBadRequest)
		} else {
			wx.routeRequest(w, &msg)
		}
	}
}

func (wx *Weixin) routeRequest(w http.ResponseWriter, r *Request) {
	requestPath := r.MsgType
	if requestPath == msgEvent {
		requestPath += "." + r.Event
	}
	for _, route := range wx.routes {
		if !route.regex.MatchString(requestPath) {
			continue
		}
		writer := responseWriter{}
		writer.wx = wx
		writer.writer = w
		writer.toUserName = r.FromUserName
		writer.fromUserName = r.ToUserName
		route.handler(writer, r)
		return
	}
	http.Error(w, "", http.StatusNotFound)
}
