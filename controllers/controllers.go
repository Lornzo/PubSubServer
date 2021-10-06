package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func GetWsUpgrader() (upgrader websocket.Upgrader) {
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) (checkOk bool) {
			checkOk = true
			return
		},
	}
	return
}

func WsMiddleWare(c *gin.Context) {
	if c.Param("channel") == "" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	c.Next()
}
