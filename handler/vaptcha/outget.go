package vaptcha

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

var challenge string
var imgid string



// @Tags test
// @Accept  json
// @Produce  json
// @Success 200 {string} plain "OK"
// @Router /test1/downtime [get]
func Outage(c *gin.Context) {
	var vaotcha Vaotcha
	if err:=c.ShouldBindQuery(&vaotcha);err!=nil{
		log.Print("binding is fail")
	}

	switch vaotcha.Action {
	case "get":
		getOutage(c)
	case "verify":
		verifyOutage(c,&vaotcha)
	}
}

// string to json bind struct
func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	log.Print(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}
	return json.NewDecoder(r.Body).Decode(target)
}

// Generate random numbers
// return about string *32byte
func getRandomStr() string {
	str := `0123456789abcdef`
	randSource := rand.NewSource(time.Now().Unix())
	numPool := rand.New(randSource)
	var re []byte
	for i := 0; i < 4; i++ {
		randNum := numPool.Intn(16)
		randStr := str[randNum]
		re = append(re, randStr)
	}
	return string(re)
}

func md5str(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	md5Re := h.Sum(nil)
	return hex.EncodeToString(md5Re)
}

func getOutage(c *gin.Context) {
	log.Print("get1")
	url := `http://d.vaptcha.com/config`
	result := new(VaKey)
	getJson(url, result)
	imgid = md5str(result.Key + getRandomStr())
	challenge = imgid
	//index key with action process
	c.JSONP(http.StatusOK, Callback{
		Code:      "0103",
		Imgid:     imgid,
		Challenge: challenge,
	})
}

func verifyOutage(c *gin.Context,vaotcha *Vaotcha)  {
	log.Print("verify1")
	if challenge!=vaotcha.Challenge {
		log.Print("not equals")
		// verification failed
		c.JSONP(http.StatusOK,gin.H{
			"code":"0104",
			"msg":"验证不匹配",
		})
		return
	}

	log.Print("cnd")
	url:=`http://d.vaptcha.com/`
	validatekey:=md5str(imgid+vaotcha.V)

	result:=new(CNDstate)
	getJson(url+validatekey,result)
	log.Print(result)
	if result.Code=="200"{
		c.JSONP(http.StatusOK,gin.H{
			"token":md5str(challenge),
			"code":"0103",
		})
	}else {
		c.JSONP(http.StatusOK,gin.H{
			"code":"0104",
			"msg":"验证不匹配",
		})
	}
}
