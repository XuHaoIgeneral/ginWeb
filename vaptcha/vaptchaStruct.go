package vaptcha

type Vaotcha struct {
	Challenge string `json:"challenge"`
	V string `json:"v"`
	Action string `json:"action"`
	Callback string `json:"callback"`
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