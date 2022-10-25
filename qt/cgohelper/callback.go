package cgohelper

import "sync"


var __GlobalCallbackRegistry = createCallbackRegistry()
func NewCallback(cb func()) (uint64, func()) {
	var id = __GlobalCallbackRegistry.register(cb)
	return id, func() {
		__GlobalCallbackRegistry.unregister(id)
	}
}
func GetCallback(id uint64) func() {
	return __GlobalCallbackRegistry.get(id)
}

type CallbackRegistry struct {
	mutex    sync.Mutex
	counter  uint64
	mapping  map[uint64] func()
}
func createCallbackRegistry() *CallbackRegistry {
	return &CallbackRegistry {
		counter: 0,
		mapping: make(map[uint64] func()),
	}
}
func (reg *CallbackRegistry) register(callback func()) uint64 {
	{ reg.mutex.Lock(); defer reg.mutex.Unlock() }
	var id = reg.counter
	reg.counter++
	reg.mapping[id] = callback
	return id
}
func (reg *CallbackRegistry) unregister(id uint64) {
	{ reg.mutex.Lock(); defer reg.mutex.Unlock() }
	delete(reg.mapping, id)
}
func (reg *CallbackRegistry) get(id uint64) func()  {
	{ reg.mutex.Lock(); defer reg.mutex.Unlock() }
	var callback, exists = reg.mapping[id]
	if exists {
		return callback
	} else {
		return func() {}
	}
}


