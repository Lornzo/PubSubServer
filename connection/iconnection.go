package connection

type IConnection interface {
	GetID() (id string)             // also connection memory address
	Listen()                        // Listen msg which from client
	SendMsg(msg []byte) (err error) // Send message to client
	Close() (err error)             // close websocket connection but not including send chan signal
	SetWindow(w chan ChanSignal)    // set chan to communicate with pool and connection object
	Ping() (err error)
}
