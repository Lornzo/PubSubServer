package connpool

import (
	"sync"
)

var defaultMessage map[string]string = make(map[string]string)
var defaultMessageLock sync.RWMutex

func SetDefaultMessage(poolName string, msg string) {
	defaultMessageLock.Lock()
	defer defaultMessageLock.Unlock()
	defaultMessage[poolName] = msg
}

func GetDefaultMessage(poolName string) (msg string) {
	var (
		poolExist bool
	)
	defaultMessageLock.RLock()
	defer defaultMessageLock.RUnlock()
	if msg, poolExist = defaultMessage[poolName]; !poolExist {
		msg = ""
	}
	return
}
