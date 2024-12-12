package main

import "fmt"

func main() {
	fmt.Printf("%v \n", Insert[int](0, 12, []int{1, 2}))
	fmt.Printf("%v \n", Insert[int](2, 12, []int{1, 2}))
	fmt.Printf("%v \n", Insert[int](1, 12, []int{1, 2}))
}
