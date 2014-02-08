package weixin

// Create Menu
func (wx *Weixin) CreateMenu(menu *Menu) error {
	gateway := weixinGroupURL + "/create?access_token="
	_, err := apiPOST(gateway, wx.tokenChan, menu)
	return err
}

// Create Menu
func (w responseWriter) CreateMenu(menu *Menu) error {
	return w.wx.CreateMenu(menu)
}
