package main

import (
	"math"
	"strconv"
	"unicode/utf8"
)

func main() {
	var a int = 456
	var b int = 123
	println(a + b)
	println(a - b)
	println(a * b)
	println(a / b)
	println(float64(a) / float64(b))
	// a = a + 1
	a++
	println(a)
	// b = b - 1
	b--
	println(b)
	//var c float64 = 1.23
	//println(a + c)
	//var d int64 = 987
	//println(a + d)
	println(math.Abs(-12.3))
	String()
	Byte()
}

func ExtremeNum() {
	println(math.MinInt64)
	println("float64 最小正数", math.SmallestNonzeroFloat64)
	println("float32 最小正数", math.SmallestNonzeroFloat32)
}

func String() {
	// he said "hello, go"
	println("he said \"hello, go\"")
	println(`hello, go
换行了。换行了
`)
	println("hello" + strconv.Itoa(123))

	println(len("hello"))
	println(len("hello你好"))
	println(utf8.RuneCountInString("hello你好"))
	//strings.CutPrefix()
}

func Byte() {
	var a byte = 'a'
	print(a)

	var str string = "hello"
	var bs []byte = []byte(str)
	var str1 string = string(bs)
	println(str1)
}

func Bool() {
	var a bool = true
	var b bool = false
	println(a && b)
	println(a || b)
	println(!a)
	// !(a&&b) => !a || !b
	// !(a||b) => !a && !b
}
