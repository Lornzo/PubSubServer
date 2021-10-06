package connpool

import "sync"

var pools map[string]IConnPool = make(map[string]IConnPool)
var poolsLock sync.RWMutex

// done
func Use(poolName string) (pool IConnPool) {

	var (
		poolExist bool
	)

	poolsLock.Lock()
	defer poolsLock.Unlock()

	if pool, poolExist = pools[poolName]; !poolExist {
		pool = newConnPool(poolName)
		pools[poolName] = pool
	}

	return
}

// not finish yet
func Close(poolName string) (err error) {
	return
}
