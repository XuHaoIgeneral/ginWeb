package main

import (
	"errors"
	"ginweb/router"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"ginweb/config"
)

var (
	cfg = pflag.StringP("config", "c", "", "apiserver config file path")
)

func main() {
	
	pflag.Parse()
	//init config about viper
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}
	//初始化数据库
	//models.Init()
	
	//set gin mode
	gin.SetMode(viper.GetString("runmode"))
	
	//create the Gin engine
	g := gin.New()

	//gin middlewares
	middlewares := []gin.HandlerFunc{}

	router.Load(
		g,
		middlewares...,
	)
	//ping the server to make sure the router is working
	go func() {
		if err := pingServer(); err != nil {
			log.Fatal("The router has no response,or it might took too long to start up", err)
		}
		log.Print("The router has been deployed siccessfully")
	}()

	log.Printf("Start to listening the incoming requests on http address: %s", viper.GetString("addr"))
	log.Printf(http.ListenAndServe(viper.GetString("addr"), g).Error())
}

func pingServer() error {
	for i := 0; i < viper.GetInt("max_ping_count"); i++ {
		//ping the server by sending a get request to `/headlth `
		resp, err := http.Get(viper.GetString("url") + "/sd/health")
		log.Print(err)
		log.Print(resp.StatusCode)
		if err == nil && resp.StatusCode == 200 {
			return nil
		}
		//sleep for a second to continue the next ping
		log.Panicf("Waiting for the router,retry in 1 second")
		time.Sleep(time.Second)
	}
	return errors.New("Connot connet to the router")
}