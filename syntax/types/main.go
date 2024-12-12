package main

import "fmt"

func main() {
	u1 := &User{}
	println(u1)
	u1 = new(User)
	println(u1)

	u2 := User{}
	fmt.Printf("%+v \n", u2)
	u2.Name = "Jerry"
	println(u2.Name)

	var u3 User
	fmt.Printf("%+v \n", u3)
	var u4 *User
	// nil
	println(u4)

	u5 := User{Name: "Jerry"}
	fmt.Printf("%+v \n", u5)
	u5 = User{18, "nick-Jerry", "Jerry"}
	fmt.Printf("%+v \n", u5)
	ChangeUser()
	Components()
}

func UseList() {
	l1 := LinkedList{}
	l1Ptr := &l1
	var l2 LinkedList = *l1Ptr
	fmt.Printf("%+v \n", l2)

	var l3 *LinkedList
	// l3 是 nil
	println(l3)
}
