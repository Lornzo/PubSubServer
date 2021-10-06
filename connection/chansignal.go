package connection

const (
	CHANSIGNAL_TYPE_CLOSE = iota
	CHANSIGNAL_TYPE_DATA
	CHANSIGNAL_TYPE_BROADCAST
)

type ChanSignal struct {
	From string
	Type int
	Msg  []byte
}
