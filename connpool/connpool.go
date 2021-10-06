package connpool

import (
	"fmt"
	"sync"
	"time"

	"github.com/Lornzo/PubSubServer/connection"
)

type connPool struct {
	name      string
	conns     map[string]connection.IConnection
	connsLock sync.RWMutex

	chanSignal     chan connection.ChanSignal
	chanSignalOnce sync.Once

	pingHandlerOnce sync.Once
}

func newConnPool(poolName string) (pool IConnPool) {

	var thePool *connPool = &connPool{
		name:       poolName,
		conns:      make(map[string]connection.IConnection),
		chanSignal: make(chan connection.ChanSignal),
	}
	go thePool.signalHandler()
	go thePool.pingHandler()

	pool = thePool

	return
}

func (thisObj *connPool) signalHandler() {
	thisObj.chanSignalOnce.Do(func() {
		for {
			var signal connection.ChanSignal = <-thisObj.chanSignal
			switch signal.Type {
			case connection.CHANSIGNAL_TYPE_CLOSE:
				thisObj.RemoveConnection(signal.From)
			case connection.CHANSIGNAL_TYPE_BROADCAST:
				thisObj.Broadcast(signal.From, signal.Msg)
			case connection.CHANSIGNAL_TYPE_DATA:
				// 資料處理
			}

		}
	})
}

func (thisObj *connPool) pingHandler() {
	thisObj.pingHandlerOnce.Do(func() {
		var (
			timer       *time.Ticker = time.NewTicker(time.Second)
			removeConns []string
		)

		for range timer.C {
			thisObj.connsLock.RLock()
			for id, conn := range thisObj.conns {
				if err := conn.Ping(); err != nil {
					removeConns = append(removeConns, id)
				}
			}
			thisObj.connsLock.RUnlock()
			thisObj.RemoveConnections(removeConns)
		}
	})

}

func (thisObj *connPool) Broadcast(from string, msg []byte) {

	var (
		id       string
		conn     connection.IConnection
		err      error
		clearArr []string
	)

	if string(msg) == GetDefaultMessage(thisObj.name) {
		return
	}

	SetDefaultMessage(thisObj.name, string(msg))

	thisObj.connsLock.RLock()

	for id, conn = range thisObj.conns {

		if id == from {
			continue
		}

		if err = conn.SendMsg(msg); err != nil {
			clearArr = append(clearArr, id)
		}

	}

	thisObj.connsLock.RUnlock()

	if len(clearArr) > 0 {
		thisObj.RemoveConnections(clearArr)
	}

}

// done
func (thisObj *connPool) AddConnection(conn connection.IConnection) (err error) {

	var (
		defaultMsg string = GetDefaultMessage(thisObj.name)
		connExist  bool
	)

	if defaultMsg != "" {
		if err = conn.SendMsg([]byte(defaultMsg)); err != nil {
			return
		}
	}

	thisObj.connsLock.Lock()
	defer thisObj.connsLock.Unlock()

	if _, connExist = thisObj.conns[conn.GetID()]; connExist {
		err = fmt.Errorf("conn had already exist")
		return
	}

	conn.SetWindow(thisObj.chanSignal)
	thisObj.conns[conn.GetID()] = conn
	return
}

func (thisObj *connPool) RemoveConnections(connIDs []string) {
	for _, connID := range connIDs {
		thisObj.RemoveConnection(connID)
	}
}

func (thisObj *connPool) RemoveConnection(connID string) (err error) {

	thisObj.connsLock.Lock()
	defer thisObj.connsLock.Unlock()

	var (
		connExist bool
		conn      connection.IConnection
	)

	if conn, connExist = thisObj.conns[connID]; connExist {
		err = conn.Close()
		delete(thisObj.conns, connID)
	}

	return

}

// not finish yet
func (thisObj *connPool) close() (err error) {
	return
}
