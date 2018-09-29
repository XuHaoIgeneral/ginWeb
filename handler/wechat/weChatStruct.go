package wechat

type AccessToken struct {
	Access_token  string `form:"access_token" json:"access_token" binding:"-"`
	Expires_in    string `form:"expires_in" json:"expires_in" binding:"-"`
	Refresh_token string `form:"refresh_token" json:"refresh_token" binding:"-"`
	Openid        string `form:"openid" json:"openid" binding:"-"`
	Scope         string `form:"scope" json:"scope" binding:"-"`
	Unionid       string `form:"unionid" json:"unionid" binding:"-"`
	Errcode       int    `form:"errcode"json:"errcode" binding:"-"`
	Errmsg        string `form:"errmsg" json:"errmsg" binding:"-"`
}
