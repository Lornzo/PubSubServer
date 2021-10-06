package connection

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

func NewWebsocket(connType string, conn *websocket.Conn) (ws IConnection) {

	ws = &wsConnection{
		connType: connType,
		conn:     conn,
		memAddr:  fmt.Sprintf("%p", conn),
	}

	return
}

type wsConnection struct {
	memAddr  string
	connType string

	conn     *websocket.Conn
	connLock sync.RWMutex

	listenOnce sync.Once
	window     chan ChanSignal
	windowLock sync.RWMutex
}

func (thisObj *wsConnection) GetID() (id string) {
	id = thisObj.memAddr
	return
}

func (thisObj *wsConnection) Listen() {
	if thisObj.conn != nil {
		thisObj.listenOnce.Do(thisObj.listenHandler)
	}
}

func (thisObj *wsConnection) listenHandler() {
	for {
		var (
			msgType int
			msgByte []byte
			err     error
		)
		if msgType, msgByte, err = thisObj.conn.ReadMessage(); err != nil {
			thisObj.close()
			break
		}

		switch msgType {
		case websocket.TextMessage:
			thisObj.broadcast(msgByte)
		case websocket.PingMessage:
			fmt.Println("ping", string(msgByte))
			thisObj.SendMsg([]byte("pong"))
		case websocket.PongMessage:
			fmt.Println("pong", string(msgByte))
		case websocket.CloseMessage:
			fmt.Println("close message")
			thisObj.close()
		default:
			thisObj.close()
		}
	}
}

func (thisObj *wsConnection) Ping() (err error) {
	err = thisObj.conn.WriteMessage(websocket.PingMessage, []byte("I say ping and you say pong : ping !"))
	return
}

func (thisObj *wsConnection) SendMsg(msg []byte) (err error) {

	thisObj.connLock.RLock()
	defer thisObj.connLock.RUnlock()

	if thisObj.conn != nil {
		err = thisObj.conn.WriteMessage(websocket.TextMessage, msg)
	}

	return
}

// 外部的方法就只是單純的把conn 給關起來，不用傳資料到外層
func (thisObj *wsConnection) Close() (err error) {

	thisObj.connLock.Lock()
	defer thisObj.connLock.Unlock()

	if thisObj.conn != nil {
		err = thisObj.conn.Close()
		thisObj.conn = nil
	}

	return
}

// 內部方法不止要把conn關掉，也要通知外層的連線池（如果有的話）
func (thisObj *wsConnection) close() (err error) {

	thisObj.Close()

	thisObj.windowLock.RLock()
	defer thisObj.windowLock.RUnlock()

	if thisObj.window != nil {
		thisObj.window <- ChanSignal{From: thisObj.memAddr, Type: CHANSIGNAL_TYPE_CLOSE}
	}

	return
}

func (thisObj *wsConnection) SetWindow(w chan ChanSignal) {
	thisObj.windowLock.Lock()
	defer thisObj.windowLock.Unlock()
	thisObj.window = w
}

func (thisObj *wsConnection) broadcast(msg []byte) (err error) {

	if thisObj.connType != CONNECITON_TYPE_BROADCAST && thisObj.connType != CONNECTION_TYPE_CHATROOM {
		err = fmt.Errorf("type %s does not support broadcast", thisObj.connType)
		return
	}

	thisObj.windowLock.RLock()
	defer thisObj.windowLock.RUnlock()

	if thisObj.window != nil {
		thisObj.window <- ChanSignal{From: thisObj.memAddr, Msg: msg, Type: CHANSIGNAL_TYPE_BROADCAST}
	} else {
		err = fmt.Errorf("communicate unit is nil")
	}

	return
}
