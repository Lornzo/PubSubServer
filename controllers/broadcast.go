package controllers

import (
	"net/http"

	"github.com/Lornzo/PubSubServer/connection"
	"github.com/Lornzo/PubSubServer/connpool"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func BroadcastHandler(c *gin.Context) {
	var (
		upGrader     websocket.Upgrader = GetWsUpgrader()
		ws           *websocket.Conn
		err          error
		channel      string = c.Param("channel")
		wsConnection connection.IConnection
	)

	if ws, err = upGrader.Upgrade(c.Writer, c.Request, nil); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// 接通的那一刻就必需要開啟監聽了
	wsConnection = connection.NewWebsocket(connection.CONNECITON_TYPE_BROADCAST, ws)
	go wsConnection.Listen()

	// 如果加入頻道失敗的話，就把連線關掉
	if err = connpool.Use(channel).AddConnection(wsConnection); err != nil {
		wsConnection.Close()
	}

	c.Request.Context().Done()
}
