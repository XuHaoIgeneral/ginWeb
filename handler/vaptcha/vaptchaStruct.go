package vaptcha

type Vaotcha struct {
	Challenge string `form:"challenge" json:"challenge" binding:"-"`
	V string `form:"v" json:"v" binding:"-"`
	Action string `form:"action" json:"action" binding:"required"`
	Callback string `form:"callback" json:"callback" binding:"required"`
}

type VaKey struct {
	Key string `json:"key"`
	Active bool `json:"active"`
}

type Callback struct {
	Code string `json:"code"`
	Imgid string `json:"imgid"`
	Challenge string `json:"challenge"`
}

type Verify struct {
	Username string `json:"username"`
	Passwd string `json:"passwd"`
	Imgid string `json:"imgid"`
	Challenge string `json:"challenge"`
	V string `json:"v"`
}

type CNDstate struct {
	Code string `json:"code"`
}