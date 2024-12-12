package main

import "fmt"

func main() {
	//Invoke()
	//Func7()
	//Recursive()
	//ClosureInvoke()
	//Defer()
	//DeferClosure()
	//DeferClosureV1()
	fmt.Println("DeferReturn", DeferReturn())
	fmt.Println("DeferReturnV1", DeferReturnV1())
	fmt.Println("DeferReturnV2", DeferReturnV2().name)
	fmt.Println("DeferReturnV3", DeferReturnV3().name)
}

func Invoke() {
	str := Func0("大明")
	println(str)
	str1, err := Func2(12, 13)
	println(str1, err)
	_, err = Func3(1, 2)

	_, _ = Func1(1, 2, 3, "abc")
	Func1(1, 2, 3, "abc")
}

func Func0(name string) string {
	return "hello, " + name
}

func Func1(a, b, c int, d string) (string, error) {
	return "hello, world", nil
}

func Func2(a int, b int) (str string, err error) {
	str = "hello, world"
	return
}

func Func3(a int, b int) (str string, err error) {
	//str = "def"
	return "abc", nil
}

func Recursive() {
	Recursive()
}

func A() {
	B()
}

func B() {
	C()
}

func C() {
	A()
}
