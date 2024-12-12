package main

import "io"

type NameI interface {
	Name() string
}

type Outer struct {
	Inner
}

type Outer1 struct {
	*Inner
}

func (i Outer1) Name() string {
	return "outer"
}

type Inner struct {
}

func (i Inner) Name() string {
	return "inner"
}

func (i Inner) Hello() {
	println("hello, 我是", i.Name())
}

type Outer2 struct {
	io.Closer
}

func Components() {
	var o Outer
	o.Hello()

	//var o1 Outer1
	o1 := Outer1{
		Inner: &Inner{},
	}
	o1.Hello()
}
