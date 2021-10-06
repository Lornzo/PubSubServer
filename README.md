# PubSubServer
使用Golang寫成的PubSubServer

## Description
* websocket 的中轉 server
* 會記錄每個頻道最後推送的資料，如果推送的資料跟上一筆重復，則不會推送

## Hot To Use
* 啟動Serverw後
* 推送方連上路由：/ws/broadcast/{頻道名稱}
* 接收方連上路由：/ws/subscribe/{頻道名稱}
* 聊天室連上路由：/ws/chatroom/{頻道名稱}