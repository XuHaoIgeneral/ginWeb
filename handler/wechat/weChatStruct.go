package wechat

type AccessToken struct {
	Access_token  string `form:"access_token" json:"access_token" binding:"required"`
	Expires_in    string `form:"expires_in" json:"expires_in" binding:"required"`
	Refresh_token string `form:"refresh_token" json:"refresh_token" binding:"required"`
	Openid        string `form:"openid" json:"openid" binding:"required"`
	Scope         string `form:"scope" json:"scope" binding:"required"`
	Unionid       string `form:"unionid" json:"unionid" binding:"required"`
	Errcode       int    `form:"errcode"json:"errcode" binding:"required"`
	Errmsg        string `form:"errmsg" json:"errmsg" binding:"required"`
}

type Err struct {
	Errcode int    `form:"errcode"json:"errcode" binding:"required"`
	Errmsg  string `form:"errmsg" json:"errmsg" binding:"required"`
}
