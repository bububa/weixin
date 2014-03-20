package weixin

import (
	//"encoding/json"
)

// Create Menu
func (wx *Weixin) CreateMenu(menu *Menu) error {
	gateway := weixinMenuURL + "/create?access_token="
	_, err := apiPOST(gateway, wx.tokenChan, menu)
	return err
}

// Create Menu
func (wx *Weixin) DeleteMenu() error {
	gateway := weixinMenuURL + "/delete?access_token="
	_, err := apiGET(gateway, wx.tokenChan)
	return err
}

// Create Menu
func (w responseWriter) CreateMenu(menu *Menu) error {
	err := w.wx.CreateMenu(menu)
    /*
	var js []byte
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
    */
	return err
}

// Delete Menu
func (w responseWriter) DeleteMenu() error {
    err := w.wx.DeleteMenu()
	return err
}
