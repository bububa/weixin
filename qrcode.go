package weixin

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Create Permenent Qrcode
func (wx *Weixin) CreateQrcode(sceneId uint64) (ticketReply *TicketReply, err error) {
	var msg struct {
		ActionName string `json:"action_name"`
		ActionInfo struct {
			Scene struct {
				SceneId uint64 `json:"scene_id"`
			} `json:"scene"`
		} `json:"action_info"`
	}
	msg.ActionName = QR_LIMIT_SCENE
	msg.ActionInfo.Scene.SceneId = sceneId
	gateway := weixinQrcodeURL + "/create?access_token="
	reply, err := apiPOST(gateway, wx.tokenChan, &msg)
	if err == nil && reply != nil {
		err = json.Unmarshal(reply, &ticketReply)
	}
	return
}

// Create Temperary Qrcode
func (wx *Weixin) CreateTempQrcode(sceneId uint64, expireSeconds uint) (ticketReply *TicketReply, err error) {
	var msg struct {
        ExpireSeconds uint `json:"expire_seconds"`
		ActionName    string `json:"action_name"`
		ActionInfo    struct {
			Scene struct {
				SceneId uint64 `json:"scene_id"`
			} `json:"scene"`
		} `json:"action_info"`
	}
	msg.ActionName = QR_SCENE
	msg.ExpireSeconds = expireSeconds
	msg.ActionInfo.Scene.SceneId = sceneId
	gateway := weixinQrcodeURL + "/create?access_token="
	reply, err := apiPOST(gateway, wx.tokenChan, &msg)
	if err == nil && reply != nil {
		err = json.Unmarshal(reply, &ticketReply)
	}
	return
}

// Get Qrcode Image
func (wx *Weixin) GetQrcodeImage(sceneId uint64) ([]byte, http.Header, int, error) {
	ticketReply, err := wx.CreateQrcode(sceneId)
	if err != nil {
		return nil, nil, 0, err
	}
	r, err := http.Get(weixinShowQrcodeURL + ticketReply.Ticket)
	if err != nil {
		return nil, r.Header, r.StatusCode, err
	}
	defer r.Body.Close()
	img, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, r.Header, r.StatusCode, err
	}
	return img, r.Header, r.StatusCode, nil
}

// Get Temperary Qrcode Image
func (wx *Weixin) GetTempQrcodeImage(sceneId uint64, expireSeconds uint) ([]byte, http.Header, int, error) {
	ticketReply, err := wx.CreateTempQrcode(sceneId, expireSeconds)
	if err != nil {
		return nil, nil, 0, err
	}
	r, err := http.Get(weixinShowQrcodeURL + ticketReply.Ticket)
	if err != nil {
		return nil, r.Header, r.StatusCode, err
	}
	defer r.Body.Close()
	img, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, r.Header, r.StatusCode, err
	}
	return img, r.Header, r.StatusCode, nil
}

// Show Qrcode
func (w responseWriter) ShowQrcode(sceneId uint64) error {
	img, headers, statusCode, err := w.wx.GetQrcodeImage(sceneId)
	copyHeaders(w.writer.Header(), headers)
	w.writer.WriteHeader(statusCode)
	w.writer.Write(img)
	var js []byte
	if err != nil {
		switch err.(type) {
		case response:
			js, _ = json.Marshal(err)
		default:
			res := response{
				ErrorCode:    0,
				ErrorMessage: err.Error(),
			}
			js, _ = json.Marshal(res)
		}
		w.writer.Write(js)
		return err
	}
	return nil
}

// Show Temp Qrcode
func (w responseWriter) ShowTempQrcode(sceneId uint64, expireSeconds uint) error {
	img, headers, statusCode, err := w.wx.GetTempQrcodeImage(sceneId, expireSeconds)
	copyHeaders(w.writer.Header(), headers)
	w.writer.WriteHeader(statusCode)
	w.writer.Write(img)
	var js []byte
	if err != nil {
		switch err.(type) {
		case response:
			js, _ = json.Marshal(err)
		default:
			res := response{
				ErrorCode:    0,
				ErrorMessage: err.Error(),
			}
			js, _ = json.Marshal(res)
		}
		w.writer.Write(js)
		return err
	}
	return nil
}

// Create Qrcode
func (w responseWriter) CreateQrcode(sceneId uint64) (ticketReply *TicketReply, err error) {
	ticketReply, err = w.wx.CreateQrcode(sceneId)
	var js []byte
	if err == nil {
		js, _ = json.Marshal(ticketReply)
	} else {
		switch err.(type) {
		case response:
			js, _ = json.Marshal(err)
		default:
			res := response{
				ErrorCode:    0,
				ErrorMessage: err.Error(),
			}
			js, _ = json.Marshal(res)
		}
	}
	w.writer.Write(js)
	return
}

// Create Temperary Qrcode
func (w responseWriter) CreateTempQrcode(sceneId uint64, expireSeconds uint) (ticketReply *TicketReply, err error) {
	ticketReply, err = w.wx.CreateTempQrcode(sceneId, expireSeconds)
	var js []byte
	if err == nil {
		js, _ = json.Marshal(ticketReply)
	} else {
		switch err.(type) {
		case response:
			js, _ = json.Marshal(err)
		default:
			res := response{
				ErrorCode:    0,
				ErrorMessage: err.Error(),
			}
			js, _ = json.Marshal(res)
		}
	}
	w.writer.Write(js)
	return
}

func copyHeaders(dst, src http.Header) {
	for k, _ := range dst {
		dst.Del(k)
	}
	for k, vs := range src {
		for _, v := range vs {
			dst.Add(k, v)
		}
	}
}
