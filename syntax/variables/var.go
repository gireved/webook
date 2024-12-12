package main

import (
	"geektime-basic-go/syntax/variables/demo"
)

var (
	a1 int     = 123
	a2 float64 = 12.3
)

func main() {
	var a int = 123
	println(a)

	a = 234

	//var a = 124
	//println(a)

	var a1 = "123"
	println(a1)

	println(a1 + "hello")

	var b = 123
	println(b)

	var c = 12.4
	println(c)

	var str = "hello"
	println(str)

	var d uint = 123
	println(d)

	var e int
	println(e)

	println(demo.Global)
	println(demo.External)
	//demo.External = "ab c"
	//println(demo.internalV1)
	//println(demo.internal)

	var f int
	println(f)

	g := 123
	println(g)

	//var h int = "123"
	//println(h)
}

const (
	Status0 = iota
	Status1
	Status2
	Abc

	Sxxx

	Status6 = 6
	Status7 = 7
)

const (
	MyStatus0 = iota<<10 + 1
	MyStatus1
)
