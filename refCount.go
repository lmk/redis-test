package main

import "sync/atomic"

type RefCounter struct {
	count int32
}

func (r *RefCounter) Reset() {
	atomic.StoreInt32(&r.count, 0)
}

func (r *RefCounter) Inc() {
	atomic.AddInt32(&r.count, 1)
}

func (r *RefCounter) Dec() {
	atomic.AddInt32(&r.count, ^int32(0))
}

func (r *RefCounter) IsRef() bool {
	buf := r.Get()

	return buf > 0
}

func (r *RefCounter) Get() int32 {
	buf := atomic.LoadInt32(&r.count)

	return buf
}

func newRefCounter() RefCounter {
	var r RefCounter
	r.Reset()
	return r
}
