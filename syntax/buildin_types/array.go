package main

import "fmt"

func Array() {
	a1 := [3]int{9, 8, 7}
	fmt.Printf("a1: %v, len: %d, cap: %d \n", a1, len(a1), cap(a1))

	a2 := [3]int{9, 8}
	fmt.Printf("a2: %v, len: %d, cap: %d \n", a2, len(a2), cap(a2))

	var a3 [3]int
	fmt.Printf("a3: %v, len: %d, cap: %d \n", a3, len(a3), cap(a3))

	//a2 = append(a2, 12)
	fmt.Printf("a1[1]: %d \n", a1[1])
	//arr1(4)
}

func arr1(idx int) {
	a1 := [3]int{9, 8, 7}
	fmt.Printf("a1[1]: %d \n", a1[idx])
}
