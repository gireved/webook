package main

import "fmt"

func ChangeUser() {
	u1 := User{Name: "Tom", Age: 18}
	fmt.Printf("%+v \n", u1)
	fmt.Printf("u1 address %p \n", &u1)
	// 实际上，当 u1.ChangeName 的时候，u1 发生了复制
	u1.ChangeName("Jerry")
	u1.ChangeAge(35)
	// Jerry, 35
	fmt.Printf("%+v \n", u1)

	println("---------------u2---------")

	u2 := &User{Name: "Tom", Age: 18}
	fmt.Printf("%+v \n", u2)
	fmt.Printf("u2 address %p \n", &u2)
	u2.ChangeName("Jerry")
	u2.ChangeAge(35)
	fmt.Printf("%+v \n", u2)
}

type User struct {
	Age      int
	Name     string
	NickName string
}

func ChangeName(u User, name string) {

}

func (u User) ChangeName(name string) {
	fmt.Printf("ChangeName 的地址 %p \n", &u)
	u.NickName = name
}

func (u *User) ChangeAge(age int) {
	u.Age = age
}
