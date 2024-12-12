package main

func Func4() {
	myFunc3 := Func3
	s, err := myFunc3(1, 2)
	println(s, err)
	_, _ = myFunc3(2, 3)

	// _, _  = Func3(1, 2)
}

func Func5() {
	fn := func(name string) string {
		return "hello, " + name
	}
	str := fn("大明")
	println(str)
}

func Func6() func(name string) string {
	return func(name string) string {
		return "hello, " + name
	}
}
func Func6Invoke() {
	fn := Func6()
	str := fn("大明")
	println(str)
}

func Func7() {
	fn := func(name string) string {
		return "hello, " + name
	}("大明")
	println(fn)
}
