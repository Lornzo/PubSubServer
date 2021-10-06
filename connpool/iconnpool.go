package connpool

import (
	"github.com/Lornzo/PubSubServer/connection"
)

type IConnPool interface {
	AddConnection(conn connection.IConnection) (err error)
	RemoveConnection(connID string) (err error)
	close() (err error)
}
