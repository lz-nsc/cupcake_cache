package cupcake_cache

import "sync"

type Call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type CallManager struct {
	calls map[string]*Call
	mu    sync.Mutex
}

func NewCallManger() *CallManager {
	return &CallManager{
		calls: map[string]*Call{},
		mu:    sync.Mutex{},
	}
}

func (cm *CallManager) Do(key string, fn func() (interface{}, error)) (val interface{}, err error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if call, ok := cm.calls[key]; ok {
		call.wg.Wait()
		return call.val, call.err
	}

	call := &Call{wg: sync.WaitGroup{}}
	cm.calls[key] = call
	call.wg.Add(1)
	val, err = fn()
	call.val = val
	call.err = err
	call.wg.Done()
	return val, err
}
