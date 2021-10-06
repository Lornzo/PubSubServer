package main

import (
	"log"

	"github.com/Lornzo/PubSubServer/configs"
	"github.com/Lornzo/PubSubServer/controllers"
	"github.com/gin-gonic/gin"
)

func main() {

	var (
		router         *gin.Engine = gin.Default()
		configFilePath string      = "/Users/a6288678/github.com/Lornzo/PubSubServer/testfiles/config.json"
		err            error
	)

	if err = configs.SetConfigFile(configFilePath); err != nil {
		log.Fatal("Reading config file faild : ", err)
		return
	}

	router.GET("/ws/subscribe/:channel", controllers.WsMiddleWare, controllers.SubscribeHandler)

	router.GET("/ws/broadcast/:channel", controllers.WsMiddleWare, controllers.BroadcastHandler)

	router.GET("/ws/chatroom/:channel", controllers.WsMiddleWare, controllers.ChatroomHandler)

	if configs.UseSSL {
		err = router.RunTLS(configs.ServePort, configs.CertFile, configs.KeyFile)
	} else {
		err = router.Run(configs.ServePort)
	}

	if err != nil {
		log.Fatal("server start faild : ", err)
	}

}
