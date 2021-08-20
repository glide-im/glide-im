package comm

import "sync/atomic"

type AtomicBool struct {
	value int32
}

func NewAtomicBool(defValue bool) *AtomicBool {
	ret := new(AtomicBool)
	if defValue {
		ret.value = 1
	} else {
		ret.value = 0
	}
	return ret
}

func (a *AtomicBool) Set(v bool) {
	var v2 int32 = 0
	if v {
		v2 = 1
	} else {
		v2 = 0
	}
	atomic.StoreInt32(&a.value, v2)
}

func (a *AtomicBool) Get() bool {
	v := atomic.LoadInt32(&a.value)
	ret := true
	if v <= 0 {
		ret = false
	}
	return ret
}

type AtomicInt64 int64

func (a *AtomicInt64) Set(v int64) {
	atomic.StoreInt64((*int64)(a), v)
}

func (a *AtomicInt64) Get() int64 {
	return atomic.LoadInt64((*int64)(a))
}
