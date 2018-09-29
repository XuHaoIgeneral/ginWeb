package vaptcha

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// @Tags test
// @Accept  json
// @Produce  json
// @Success 200 {string} plain "OK"
// @Router /test/validate [post]
func Validate(c *gin.Context) {

	var verify Verify
	var err error
	contentType := c.Request.Header.Get("Content-Type")
	log.Print(contentType)
	switch contentType {
	case "application/json":
		err = c.BindJSON(&verify)
	case "application/x-www-form-urlencoded; charset=UTF-8":
		err = c.BindJSON(&verify)
	default:
		log.Print("error!!!!")
	}
	if err != nil {
		log.Print("error about request header application")
		log.Print(err)
	}
	log.Print(verify)
	message := verify.Imgid
	c.JSON(http.StatusOK, gin.H{
		"state":   true,
		"test":    "test",
		"message": message,
	})
}
