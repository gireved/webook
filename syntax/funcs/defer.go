package main

import (
	"os"
	"sync"
)

func Defer() {
	defer func() {
		println("第一个 defer")
	}()
	defer func() {
		println("第二个 defer")
	}()

	//println("第二个 defer")
	//println("第一个 defer")
}

func DeferLoop(max int) {
	for i := 0; i < max; i++ {
		defer func() {
			println("hello")
		}()
	}
}

func DeferClosure() {
	j := 0
	defer func() {
		println(j)
	}()
	j = 1
	//println(j)
}

func DeferClosureV1() {
	j := 0
	defer func(j int) {
		println(j)
	}(j)
	j = 1
}

func DeferReturn() int {
	a := 0
	defer func() {
		a = 1
	}()
	return a
}

func DeferReturnV0() (b int) {
	a := 0
	defer func() {
		a = 1
	}()
	b = a
	return
}

func DeferReturnV1() (a int) {
	a = 0
	defer func() {
		a = 1
	}()
	return a
}

func DeferReturnV2() *MyStruct {
	res := &MyStruct{
		name: "Tom",
	}
	defer func() {
		res.name = "Jerry"
	}()
	return res
}

func DeferReturnV3() MyStruct {
	res := MyStruct{
		name: "Tom",
	}
	defer func() {
		res.name = "Jerry"
	}()
	return res
}

type MyStruct struct {
	name string
}

type SafeResourceV1 struct {
	lock     *sync.Mutex
	resource any
}

func (s SafeResourceV1) UseResource() {
	s.lock.Lock()
	defer s.lock.Unlock()
}

type SafeResource struct {
	lock     sync.Mutex
	resource any
}

func (s *SafeResource) UseResource() {
	s.lock.Lock()
	defer s.lock.Unlock()
}

func ReadFile(file string) {
	f, err := os.Open(file)
	if err != nil {
		println(err)
		return
	}
	defer f.Close()
}
