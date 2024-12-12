package main

import "fmt"

type Fish struct {
}

func (f Fish) Swim() {

}

type FakeFish Fish

func (f FakeFish) FakeSwim() {

}

type Yu = Fish

func UseFish() {
	f1 := Fish{}
	f1.Swim()

	f2 := FakeFish{}
	f2.FakeSwim()

	f3 := FakeFish(f1)

	f1 = Fish(f2)
	fmt.Println(f1, f3)

	yu := Yu{}
	yu.Swim()
}
