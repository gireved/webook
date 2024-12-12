package main

import "sync/atomic"

func main() {
	var val int32 = 12
	//
	val = atomic.LoadInt32(&val)
	println(val)
	// 这个可以确保，缓存里面的 val 被改成 14 了
	atomic.StoreInt32(&val, 14)
	newVal := atomic.AddInt32(&val, 1)
	println(newVal)
	swapped := atomic.CompareAndSwapInt32(&val, 13, 15)
	println(swapped)
}
