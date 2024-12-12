package main

var i = 0

func Closure(name string) func() string {
	// 闭包
	// name 变量
	// 方法本身
	return func() string {
		i++
		world := "world"
		return "hello, " + name + world
	}
}

func ClosureInvoke() {
	c := Closure("大明")
	println(c())
}
