package service

import "sync/atomic"

var (
	taskIdCounter int64 = 0
)

func GetTaskId() int64 {
	atomic.AddInt64(&taskIdCounter, 1)
	return atomic.LoadInt64(&taskIdCounter)
}

func PeekTaskId() int64 {
	return atomic.LoadInt64(&taskIdCounter)
}
